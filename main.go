package main

import (
	"context"
	"encoding/json"
	"log"
	"math"
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

func decomposeBalanceUpdate(
	ra map[string]ReferencedAsset,
	bu *events.StateUpdate_BalanceUpdate,
) (string, ReferencedAsset, uint64, bool) {
	var asset ReferencedAsset

	if assetId := bu.AmountAfter.AssetId; assetId != nil {
		asset = ra[toBase58String(assetId)]
	} else {
		asset = WavesAsset
	}

	address := toBase58String(bu.Address)
	after := uint64(bu.AmountAfter.Amount)
	isInbound := (bu.AmountAfter.Amount - bu.AmountBefore) > 0
	return address, asset, after, isInbound
}

func getInsertion(updateResponse *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
	var insertions []BalanceUpdate
	update := updateResponse.GetUpdate()
	appendField := update.GetAppend()
	block := appendField.GetBlock()
	height := uint32(update.Height)
	blockTimestamp := time.Unix(0, block.Block.Header.Timestamp*int64(time.Millisecond))
	stateUpdatesBalances := appendField.StateUpdate.GetBalances()
	transactionIds := appendField.TransactionIds
	transactionStateUpdates := appendField.TransactionStateUpdates
	orderedReferencedAssets := update.ReferencedAssets
	referencedAssets := make(map[string]ReferencedAsset)

	for _, ra := range orderedReferencedAssets {
		id := toBase58String(ra.Id)
		referencedAssets[id] = ReferencedAsset{
			id:       id,
			decimals: uint8(ra.Decimals),
			name:     ra.Name,
		}
	}

	for _, bu := range stateUpdatesBalances {
		address, asset, after, isInbound := decomposeBalanceUpdate(referencedAssets, bu)

		insertions = append(insertions, BalanceUpdate{
			height:         height,
			asset:          asset.id,
			address:        address,
			txVersion:      0,
			txType:         0,
			tx:             "",
			balance:        after,
			assetDecimals:  asset.decimals,
			assetName:      asset.name,
			isInbound:      isInbound,
			isStateUpdate:  true,
			timestamp:      blockTimestamp,
			blockTimestamp: blockTimestamp,
		})
	}

	for i, tsu := range transactionStateUpdates {
		txId := toBase58String(transactionIds[i])

		var timestamp time.Time
		var txType uint8
		var txVersion uint8

		if txData := block.Block.Transactions[i].GetWavesTransaction(); txData != nil {
			timestamp = time.Unix(0, txData.Timestamp*int64(time.Millisecond))
			txVersion = uint8(txData.Version)
			txType = TxTypes[reflect.TypeOf(txData.Data).Elem().Name()]
		} else {
			timestamp = time.Unix(0, appendField.TransactionsMetadata[i].GetEthereum().Timestamp*int64(time.Millisecond))
			txVersion = 0
			txType = TxTypes["Transaction_Ethereum"]
		}

		for _, bu := range tsu.Balances {
			address, asset, after, isInbound := decomposeBalanceUpdate(referencedAssets, bu)

			insertions = append(insertions, BalanceUpdate{
				height:         height,
				asset:          asset.id,
				address:        address,
				txVersion:      txVersion,
				txType:         txType,
				tx:             txId,
				balance:        after,
				assetDecimals:  asset.decimals,
				assetName:      asset.name,
				isInbound:      isInbound,
				isStateUpdate:  false,
				timestamp:      timestamp,
				blockTimestamp: blockTimestamp,
			})
		}
	}

	return insertions

}

func getChInput(ch chan *blockchainupdates.GetBlockUpdateResponse) proto.Input {
	height := proto.ColUInt32{}
	asset := proto.ColStr{}
	address := proto.ColStr{}
	txVersion := proto.ColUInt8{}
	txType := proto.ColUInt8{}
	tx := proto.ColStr{}
	balance := proto.ColUInt64{}
	assetDecimals := proto.ColUInt8{}
	assetName := proto.ColStr{}
	isInbound := proto.ColBool{}
	isStateUpdate := proto.ColBool{}
	timestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMilli)
	blockTimestamp := new(proto.ColDateTime64).WithPrecision(proto.PrecisionMilli)

	for u := range ch {
		for _, ins := range getInsertion(u) {
			height.Append(ins.height)
			asset.Append(ins.asset)
			address.Append(ins.address)
			txVersion.Append(ins.txVersion)
			txType.Append(ins.txType)
			tx.Append(ins.tx)
			balance.Append(ins.balance)
			assetDecimals.Append(ins.assetDecimals)
			assetName.Append(ins.assetName)
			isInbound.Append(ins.isInbound)
			isStateUpdate.Append(ins.isStateUpdate)
			timestamp.Append(ins.timestamp)
			blockTimestamp.Append(ins.blockTimestamp)
		}
	}

	return proto.Input{
		{Name: "height", Data: &height},
		{Name: "asset", Data: &asset},
		{Name: "address", Data: &address},
		{Name: "tx_version", Data: &txVersion},
		{Name: "tx_type", Data: &txType},
		{Name: "tx", Data: &tx},
		{Name: "balance", Data: &balance},
		{Name: "asset_decimals", Data: &assetDecimals},
		{Name: "asset_name", Data: &assetName},
		{Name: "is_inbound", Data: &isInbound},
		{Name: "is_state_update", Data: &isStateUpdate},
		{Name: "timestamp", Data: timestamp},
		{Name: "block_timestamp", Data: blockTimestamp},
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

func getCurrentHeight(client *http.Client, path string) int32 {
	r, err := (client).Get(restApiHost + path)
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

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	restApiHost = os.Getenv("REST_API_HOST")
	httpClient := http.Client{}

	// Устанавливаем соединение с gRPC-сервером
	grpcConn, err := grpc.NewClient(
		os.Getenv("GRPC_HOSTNAME")+":"+os.Getenv("GRPC_PORT"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32*256)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,  // send pings every 10 seconds if there is no activity
			Timeout:             30 * time.Second, // wait 1 second for ping back
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
			Address: os.Getenv("CH_HOSTNAME") +
				":" + os.Getenv("CH_PORT"),
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

	// probe
	// consume_ctx := context.Background()
	// probe_ctx, cancel_probe_ctx := context.WithTimeout(consume_ctx, 1*time.Second)
	if _, err := client.GetBlockUpdate(
		context.Background(),
		&blockchainupdates.GetBlockUpdateRequest{Height: 1},
	); err != nil {
		// cancel_probe_ctx()
		panic(err)
	} else {
		log.Println("Probe OK")
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			panic(err)
		}
	}()

	var maxHeightResult proto.ColUInt32

	if err := dbConn.Do(dbCtx, ch.Query{
		Body:   "SELECT max(height) FROM blocks",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	}); err != nil {
		panic(err)
	}

	counter := int32(maxHeightResult[0])

	log.Printf("Consume from: %d", counter)

	toHeight := getCurrentHeight(&httpClient, "/blocks/height")

	estimatedBlocks := toHeight - counter
	batchCount := int(estimatedBlocks / BlocksStep)
	blocksStep := BlocksStep

	if estimatedBlocks < blocksStep {
		blocksStep = estimatedBlocks
	}

	var chWg sync.WaitGroup

	start := time.Now()
	workTime := time.Duration(0)
	for i := 1; i <= batchCount; i++ {
		from := counter + 1
		start := time.Now()
		results := make(chan *blockchainupdates.GetBlockUpdateResponse, blocksStep)
		getBlocks(client, results, from, blocksStep)

		insertions := getChInput(results)

		chWg.Add(1)
		go func() {
			defer chWg.Done()

			if err := dbConn.Do(dbCtx, ch.Query{
				Body:  "INSERT INTO blocks VALUES",
				Input: proto.Input{insertions[0]},
			}); err != nil {
				panic(err)
			}

			if err := dbConn.Do(dbCtx, ch.Query{
				Body:  "INSERT INTO balances_updates VALUES",
				Input: insertions,
			}); err != nil {
				panic(err)
			}
		}()
		counter = counter + BlocksStep
		workTime += time.Since(start)
		avgBatchTime := workTime.Seconds() / float64(i)
		log.Printf(
			"Batch:\t%d/%d Range:\t(%d-%d) Time:\t%.2fs Avg:\t%.2fs ETA:\t%.2fh",
			i,
			batchCount,
			from,
			counter,
			time.Since(start).Seconds(),
			avgBatchTime,
			float64(batchCount-i)*avgBatchTime/60/60,
		)
	}

	chWg.Wait()

	log.Printf("Total: %.2fs", time.Since(start).Seconds())
}
