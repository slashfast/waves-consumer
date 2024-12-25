package connectors

import (
	"context"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/waves"
	"go.uber.org/fx"
	"math"
	"math/bits"
	"net/url"
)

type ClickHouseConfig interface {
	GetClickHouseDSN() string
}

type ClickHouseConnector struct {
	conn *ch.Client
	log  *zerolog.Logger
}

func (chc *ClickHouseConnector) Insert(updates []waves.BalanceUpdate) error {
	return chc.insert(getInput(updates))
}

func (chc *ClickHouseConnector) insert(input proto.Input) error {
	return chc.conn.Do(context.Background(), ch.Query{
		Body:  "INSERT INTO balances_updates VALUES",
		Input: input,
		Settings: []ch.Setting{
			{Key: "async_insert", Value: "1"},
			{Key: "wait_for_async_insert", Value: "0"},
		},
	})
}

func (chc *ClickHouseConnector) GetHeight() (int32, error) {
	var maxHeightResult proto.ColUInt32
	if err := chc.conn.Do(context.Background(), ch.Query{
		Body:   "SELECT max(height) FROM balances_updates",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	}); err != nil {
		return 0, err
	}
	return int32(maxHeightResult[0]), nil
}

func getInput(updates []waves.BalanceUpdate) proto.Input {
	height := proto.ColUInt32{}
	assetId := proto.ColStr{}
	address := proto.ColFixedStr{Size: 35}
	txVersion := proto.ColUInt8{}
	txType := proto.ColUInt8{}
	txId := proto.ColStr{}
	counterparty := proto.ColStr{}
	amount := proto.ColDecimal128{}
	delta := proto.ColDecimal128{}
	subunits := proto.ColUInt64{}
	deltaSubunits := proto.ColInt64{}
	assetDecimals := proto.ColUInt8{}
	assetName := proto.ColStr{}
	isInbound := proto.ColBool{}
	isStateUpdate := proto.ColBool{}
	txTimestamp := proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}
	blockTimestamp := proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}

	for _, ins := range updates {
		_amount, _delta := subunitsToDecimal(ins.Subunits, ins.DeltaSubunits, ins.AssetDecimals)
		height.Append(ins.Height)
		assetId.Append(ins.AssetId)
		address.Append([]byte(ins.Address))
		txVersion.Append(ins.TxVersion)
		txType.Append(ins.TxType)
		txId.Append(ins.TxId)
		counterparty.Append(ins.Counterparty)
		amount.Append(_amount)
		delta.Append(_delta)
		subunits.Append(ins.Subunits)
		deltaSubunits.Append(ins.DeltaSubunits)
		assetDecimals.Append(ins.AssetDecimals)
		assetName.Append(ins.AssetName)
		isInbound.Append(ins.IsInbound)
		isStateUpdate.Append(ins.IsStateUpdate)
		txTimestamp.Append(ins.TxTimestamp)
		blockTimestamp.Append(ins.BlockTimestamp)
	}

	return proto.Input{
		{Name: "height", Data: &height},
		{Name: "asset_id", Data: &assetId},
		{Name: "address", Data: &address},
		{Name: "tx_version", Data: &txVersion},
		{Name: "tx_type", Data: &txType},
		{Name: "tx_id", Data: &txId},
		{Name: "counterparty", Data: &counterparty},
		{Name: "amount", Data: proto.Alias(&amount, "Decimal128(8)")},
		{Name: "delta", Data: proto.Alias(&delta, "Decimal128(8)")},
		{Name: "subunits", Data: &subunits},
		{Name: "delta_subunits", Data: &deltaSubunits},
		{Name: "asset_decimals", Data: &assetDecimals},
		{Name: "asset_name", Data: &assetName},
		{Name: "is_inbound", Data: &isInbound},
		{Name: "is_state_update", Data: &isStateUpdate},
		{Name: "tx_timestamp", Data: &txTimestamp},
		{Name: "block_timestamp", Data: &blockTimestamp},
	}
}

func (chc *ClickHouseConnector) createTables(dbName string) error {
	if err := chc.conn.Do(context.Background(), ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbName + `.balances_updates
		(
			asset_id        String,
			address         FixedString(35),
			height          UInt32,
			tx_version      UInt8,
			tx_type         UInt8,
			tx_id           String,
   		    counterparty    String,
			amount          Decimal128(8),
			delta           Decimal128(8),
			subunits        UInt64,
			delta_subunits  Int64,
			asset_decimals  UInt8,
			asset_name      String,
			is_inbound      Bool,
			is_state_update Bool,
			tx_timestamp    DateTime64(3),
			block_timestamp DateTime64(3)
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (asset_id, address, height, tx_id)
       SETTINGS non_replicated_deduplication_window = 1000
		`,
	}); err != nil {
		return err
	}

	return nil
}

func NewConnector(lc fx.Lifecycle, l *zerolog.Logger, config ClickHouseConfig) (*ClickHouseConnector, error) {
	chConfig, err := url.Parse(config.GetClickHouseDSN())
	if err != nil {
		return nil, err
	}

	chPassword, _ := chConfig.User.Password()
	dbName := chConfig.Path[1:]
	conn, err := ch.Dial(
		context.Background(),
		ch.Options{
			Address:  chConfig.Host,
			Database: dbName,
			User:     chConfig.User.Username(),
			Password: chPassword,
		},
	)
	if err != nil {
		return nil, err
	}

	connector := ClickHouseConnector{conn: conn, log: l}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := connector.createTables(dbName)
			if err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return conn.Close()
		},
	})

	return &connector, nil
}

func mulInt64(x int64, y uint64) (uint64, uint64) {
	if x < 0 {
		return math.MaxUint64, uint64(x * int64(y))
	}
	return bits.Mul64(uint64(x), y)
}

func subunitsToDecimal(subunits uint64, deltaSubunits int64, decimals uint8) (proto.Decimal128, proto.Decimal128) {
	power := 8 - decimals

	var i uint8
	var multiplier uint64 = 1

	for i = 0; i < power; i++ {
		multiplier *= 10
	}

	amountHigh, amountLow := bits.Mul64(subunits, multiplier)
	amount := proto.Decimal128{Low: amountLow, High: amountHigh}
	deltaHigh, deltaLow := mulInt64(deltaSubunits, multiplier)
	delta := proto.Decimal128{Low: deltaLow, High: deltaHigh}

	return amount, delta
}
