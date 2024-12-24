package clickhouse

import (
	"context"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/slashfast/wavesconsumer/v2/waves"
	"math"
	"math/bits"
)

func InsertToDB(db *ch.Client, input proto.Input) error {
	return db.Do(context.Background(), ch.Query{
		Body:  "INSERT INTO balances_updates VALUES",
		Input: input,
		Settings: []ch.Setting{
			{Key: "async_insert", Value: "1"},
			{Key: "wait_for_async_insert", Value: "0"},
		},
	})
}

func MulInt64(x int64, y uint64) (uint64, uint64) {
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
	deltaHigh, deltaLow := MulInt64(deltaSubunits, multiplier)
	delta := proto.Decimal128{Low: deltaLow, High: deltaHigh}

	return amount, delta
}

func GetInput(updates []waves.BalanceUpdate) proto.Input {
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
