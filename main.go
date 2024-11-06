package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"math/bits"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/joho/godotenv"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var restApiHost string
var referencedAssets = make(map[string]ReferencedAsset)
var httpClient = http.Client{}

func MulInt64(x int64, y uint64) (uint64, uint64) {
	if x < 0 {
		return math.MaxUint64, uint64(x * int64(y))
	}
	return bits.Mul64(uint64(x), y)
}

func decomposeBalanceUpdate(
	bu *events.StateUpdate_BalanceUpdate,
) (string, ReferencedAsset, uint64, int64, bool) {
	var asset ReferencedAsset

	if assetId := bu.AmountAfter.AssetId; assetId != nil {
		b58assetId := toBase58String(assetId)
		var ok bool
		asset, ok = referencedAssets[b58assetId]
		if !ok {
			asset = getAssetById(b58assetId)
			referencedAssets[b58assetId] = asset
		}
	} else {
		asset = WavesAsset
	}

	address := toBase58String(bu.Address)
	beforeSubunits := bu.AmountBefore
	afterSubunits := bu.AmountAfter.Amount
	deltaSubunits := afterSubunits - beforeSubunits
	isInbound := deltaSubunits > 0

	return address, asset, uint64(afterSubunits), deltaSubunits, isInbound
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

func getInsertion(updateResponse *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
	var insertions []BalanceUpdate
	update := updateResponse.GetUpdate()
	appendField := update.GetAppend()
	block := appendField.GetBlock()
	height := uint32(update.Height)
	blockTimestamp := time.Unix(0, block.Block.Header.Timestamp*int64(time.Millisecond)).UTC()
	stateUpdatesBalances := appendField.StateUpdate.GetBalances()
	transactionIds := appendField.TransactionIds
	transactionStateUpdates := appendField.TransactionStateUpdates
	orderedReferencedAssets := update.ReferencedAssets

	for _, ra := range orderedReferencedAssets {
		id := toBase58String(ra.Id)
		if _, ok := referencedAssets[id]; !ok {
			referencedAssets[id] = ReferencedAsset{
				id:       id,
				decimals: uint8(ra.Decimals),
				name:     ra.Name,
			}
		}
	}

	for _, bu := range stateUpdatesBalances {
		address, asset, subunits, deltaSubunits, isInbound := decomposeBalanceUpdate(bu)
		amount, delta := subunitsToDecimal(subunits, deltaSubunits, asset.decimals)
		insertions = append(insertions, BalanceUpdate{
			height:         height,
			assetId:        asset.id,
			address:        address,
			txVersion:      0,
			txType:         0,
			txId:           "",
			counterparty:   "",
			amount:         amount,
			delta:          delta,
			subunits:       subunits,
			deltaSubunits:  deltaSubunits,
			assetDecimals:  asset.decimals,
			assetName:      asset.name,
			isInbound:      isInbound,
			isStateUpdate:  true,
			txTimestamp:    blockTimestamp,
			blockTimestamp: blockTimestamp,
		})
	}

	for i, tsu := range transactionStateUpdates {
		txId := toBase58String(transactionIds[i])

		var timestamp time.Time
		var txType uint8
		var txVersion uint8
		var counterparty string

		if txData := block.Block.Transactions[i].GetWavesTransaction(); txData != nil {
			timestamp = time.Unix(0, txData.Timestamp*int64(time.Millisecond)).UTC()
			txVersion = uint8(txData.Version)
			txType = TxTypes[reflect.TypeOf(txData.Data).Elem().Name()]
			txMetadata := appendField.TransactionsMetadata[i]
			counterparty = toBase58String(txMetadata.SenderAddress)
		} else {
			txMetadata := appendField.TransactionsMetadata[i]
			txData := txMetadata.GetEthereum()
			timestamp = time.Unix(0, txData.Timestamp*int64(time.Millisecond)).UTC()
			txVersion = 0
			txType = TxTypes["Transaction_Ethereum"]
			counterparty = toBase58String(txMetadata.SenderAddress)
		}

		for _, bu := range tsu.Balances {
			address, asset, subunits, deltaSubunits, isInbound := decomposeBalanceUpdate(bu)
			amount, delta := subunitsToDecimal(subunits, deltaSubunits, asset.decimals)

			insertions = append(insertions, BalanceUpdate{
				height:         height,
				assetId:        asset.id,
				address:        address,
				txVersion:      txVersion,
				txType:         txType,
				txId:           txId,
				counterparty:   counterparty,
				amount:         amount,
				delta:          delta,
				subunits:       subunits,
				deltaSubunits:  deltaSubunits,
				assetDecimals:  asset.decimals,
				assetName:      asset.name,
				isInbound:      isInbound,
				isStateUpdate:  false,
				txTimestamp:    timestamp,
				blockTimestamp: blockTimestamp,
			})
		}
	}

	return insertions

}

func getChInput(ch chan *blockchainupdates.GetBlockUpdateResponse) proto.Input {
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

	for u := range ch {
		for _, ins := range getInsertion(u) {
			height.Append(ins.height)
			assetId.Append(ins.assetId)
			address.Append([]byte(ins.address))
			txVersion.Append(ins.txVersion)
			txType.Append(ins.txType)
			txId.Append(ins.txId)
			counterparty.Append(ins.counterparty)
			amount.Append(ins.amount)
			delta.Append(ins.delta)
			subunits.Append(ins.subunits)
			deltaSubunits.Append(ins.deltaSubunits)
			assetDecimals.Append(ins.assetDecimals)
			assetName.Append(ins.assetName)
			isInbound.Append(ins.isInbound)
			isStateUpdate.Append(ins.isStateUpdate)
			txTimestamp.Append(ins.txTimestamp)
			blockTimestamp.Append(ins.blockTimestamp)
		}
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

// []*blockchain_updates.GetBlockUpdateResponse
func getBlocks(client blockchainupdates.BlockchainUpdatesApiClient, ch chan *blockchainupdates.GetBlockUpdateResponse, from int32, len int32) {
	var wg sync.WaitGroup
	toHeight := from + len - 1
	for height := from; height <= toHeight; height++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				resp, err := client.GetBlockUpdate(
					context.Background(),
					&blockchainupdates.GetBlockUpdateRequest{Height: height},
				)
				if err == nil {
					// log.Printf("Height:\t%d", resp.Update.Height)
					ch <- resp
					return
				}
				log.Panicln(err)
				//log.Printf("Error: %v", err)
				//time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	close(ch)
}

func getCurrentHeight() int32 {
	r, err := httpClient.Get(restApiHost + "/blocks/height")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var target BlocksHeightResponse

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		panic(err)
	}

	return target.Height - 1
}

func getAssetById(assetId string) ReferencedAsset {
	r, err := httpClient.Get(restApiHost + "/assets/details/" + assetId + "?full=true")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var target AssetDetails

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		panic(err)
	}

	return ReferencedAsset{
		id:       target.AssetId,
		decimals: uint8(target.Decimals),
		name:     target.Name,
	}
}

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Println("The dotenv file is missing")
	}

	restApiHost = os.Getenv("REST_API_HOST")

	// Устанавливаем соединение с gRPC-сервером
	grpcConn, err := grpc.NewClient(
		os.Getenv("GRPC_HOSTNAME")+":"+os.Getenv("GRPC_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32*256)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		panic(err)
	}

	// Закрываем соединение после завершения работы
	defer func() {
		if err := grpcConn.Close(); err != nil {
			panic(err)
		}
	}()

	dbName := os.Getenv("CH_DATABASE")
	dbCtx := context.Background()
	dbConn, err := ch.Dial(
		dbCtx,
		ch.Options{
			Address:  os.Getenv("CH_HOSTNAME") + ":" + os.Getenv("CH_PORT"),
			Database: dbName,
			User:     os.Getenv("CH_USER"),
			Password: os.Getenv("CH_PASSWORD"),
		},
	)
	if err != nil {
		panic(err)
	}

	if err := createTables(dbConn, dbCtx, dbName); err != nil {
		panic(err)
	}

	client := blockchainupdates.NewBlockchainUpdatesApiClient(grpcConn)

	// Проверка соединения
	if _, err := client.GetBlockUpdate(
		context.Background(),
		&blockchainupdates.GetBlockUpdateRequest{Height: 1},
	); err != nil {
		panic(err)
	} else {
		log.Println("Probe OK")
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			panic(err)
		}
	}()

	// Получаем начальную максимальную высоту из базы данных
	var maxHeightResult proto.ColUInt32
	if err := dbConn.Do(dbCtx, ch.Query{
		Body:   "SELECT max(height) FROM balances_updates",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	}); err != nil {
		panic(err)
	}
	currentHeight := int32(maxHeightResult[0]) // Начальная высота из базы данных
	log.Printf("Starting from height: %d", currentHeight)

	newBlockchainHeight := getCurrentHeight() // Запрос текущей высоты блокчейна

	var heightDiff int32
	blocksStep := BlocksStep
	var prevBlockStep int32 = 0 // Переменная для хранения предыдущей разницы

	for {
		heightDiff = newBlockchainHeight - currentHeight
		realTimeMode := !(heightDiff >= BlocksStep)

		if realTimeMode {
			newBlockchainHeight = getCurrentHeight() // Запрос текущей высоты блокчейна
			log.Printf("DB Height: %d, Blockchain Height: %d", currentHeight, newBlockchainHeight)
			heightDiff = newBlockchainHeight - currentHeight
			if heightDiff == 0 {
				log.Println("No new blocks. Waiting for 1 minute...")
				time.Sleep(1 * time.Minute)
				continue
			} else if heightDiff < 0 {
				log.Panicln("Height difference is negative! Possible data inconsistency.")
			}
			blocksStep = heightDiff
		}

		// Печатаем информацию только если разница отличается от предыдущей
		if prevBlockStep != blocksStep {
			log.Printf("Processing blocks from height %d to %d with step %d", currentHeight, newBlockchainHeight, blocksStep)
		}

		results := make(chan *blockchainupdates.GetBlockUpdateResponse, blocksStep)
		getBlocks(client, results, currentHeight+1, blocksStep)

		insertions := getChInput(results)
		if err := dbConn.Do(context.Background(), ch.Query{
			Body:  "INSERT INTO balances_updates VALUES",
			Input: insertions,
			Settings: []ch.Setting{
				{Key: "async_insert", Value: "1"},
				{Key: "wait_for_async_insert", Value: "0"},
			},
		}); err != nil {
			panic(err)
		}

		// Обновляем текущую высоту
		currentHeight += blocksStep
		log.Printf("Updated current height to: %d", currentHeight)

		prevBlockStep = blocksStep
	}
}
