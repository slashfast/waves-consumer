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
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	blockchain_updates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var txTypes = map[string]uint8{
	"Transaction_Genesis":         1,
	"Transaction_Payment":         2,
	"Transaction_Issue":           3,
	"Transaction_Transfer":        4,
	"Transaction_Reissue":         5,
	"Transaction_Exchange":        7,
	"Transaction_Lease":           8,
	"Transaction_LeaseCancel":     9,
	"Transaction_Burn":            6,
	"Transaction_CreateAlias":     10,
	"Transaction_MassTransfer":    11,
	"Transaction_DataTransaction": 12,
	"Transaction_SponsorFee":      14,
	"Transaction_SetScript":       13,
	"Transaction_SetAssetScript":  15,
	"Transaction_InvokeScript":    16,
	"Transaction_UpdateAssetInfo": 17,
	"Transaction_Ethereum":        18,
}

type BalanceUpdate struct {
	height          uint32
	asset           string
	address         string
	tx_version      uint8
	tx_type         uint8
	tx              string
	balance         uint64
	asset_decimals  uint8
	asset_name      string
	is_inbound      bool
	is_state_update bool
	timestamp       time.Time
	block_timestamp time.Time
}

type ReferencedAsset struct {
	id       string
	decimals uint8
	name     string
}

type BlocksHeightResponse struct {
	Height int32 `json:"height"`
}

var WAVES_ASSET = ReferencedAsset{
	id:       "WAVES",
	decimals: 8,
	name:     "Waves",
}

func toBase58String(addr []byte) string {
	return base58.Encode(addr)
}

const BLOCKS_STEP int32 = 1000

const (
	WAVES_ID           byte = 0x57 // blockchain ID
	RAW_ADDRESS_PREFIX byte = 0x01 // address first byte
)

var addressPrefix = []byte{RAW_ADDRESS_PREFIX, WAVES_ID}

func wavesHash(data []byte) []byte {
	blakeHash, err := blake2b.New256(nil)

	if err != nil {
		panic(err)
	}

	blakeHash.Write(data)

	blakeSum := blakeHash.Sum(nil)
	keccakHash := sha3.NewLegacyKeccak256()

	keccakHash.Write(blakeSum)

	return keccakHash.Sum(nil)
}

func addressFromPubKey(data []byte) string {
	// append generatorAddress hash (20 bytes) to the addressPrefix
	generatorAddress := append(addressPrefix, wavesHash(data)[:20]...)

	// append checksum (4 bytes) to the generatorAddress
	checksum := wavesHash(generatorAddress)[:4]
	generatorAddress = append(generatorAddress, checksum...)

	// encode to base58
	return base58.Encode(generatorAddress)
}

func exctractBalanceUpdate(ra map[string]ReferencedAsset, bu *events.StateUpdate_BalanceUpdate) (string, ReferencedAsset, uint64, bool) {
	var asset ReferencedAsset

	if assetId := bu.AmountAfter.AssetId; assetId != nil {
		asset = ra[toBase58String(assetId)]
	} else {
		asset = WAVES_ASSET
	}

	address := toBase58String(bu.Address)
	after := uint64(bu.AmountAfter.Amount)
	is_inbound := (bu.AmountAfter.Amount - bu.AmountBefore) > 0
	return address, asset, after, is_inbound
}

func getInsertion(updateResponse *blockchain_updates.GetBlockUpdateResponse) []BalanceUpdate {
	var insertions []BalanceUpdate
	var update *events.BlockchainUpdated = updateResponse.GetUpdate()
	appendField := update.GetAppend()
	block := appendField.GetBlock()
	height := uint32(update.Height)
	block_timestamp := time.Unix(0, block.Block.Header.Timestamp*int64(time.Millisecond))
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
		address, asset, after, is_inbound := exctractBalanceUpdate(referencedAssets, bu)

		insertions = append(insertions, BalanceUpdate{
			height:          height,
			asset:           asset.id,
			address:         address,
			tx_version:      0,
			tx_type:         0,
			tx:              "",
			balance:         after,
			asset_decimals:  asset.decimals,
			asset_name:      asset.name,
			is_inbound:      is_inbound,
			is_state_update: true,
			timestamp:       block_timestamp,
			block_timestamp: block_timestamp,
		})
	}

	for i, tsu := range transactionStateUpdates {
		tx_id := toBase58String(transactionIds[i])

		var timestamp time.Time
		var tx_type uint8
		var tx_version uint8

		if tx_data := block.Block.Transactions[i].GetWavesTransaction(); tx_data != nil {
			timestamp = time.Unix(0, tx_data.Timestamp*int64(time.Millisecond))
			tx_version = uint8(tx_data.Version)
			tx_type = txTypes[reflect.TypeOf(tx_data.Data).Elem().Name()]
		} else {
			timestamp = time.Unix(0, appendField.TransactionsMetadata[i].GetEthereum().Timestamp*int64(time.Millisecond))
			tx_version = 0
			tx_type = txTypes["Transaction_Ethereum"]
		}

		for _, bu := range tsu.Balances {
			address, asset, after, is_inbound := exctractBalanceUpdate(referencedAssets, bu)

			insertions = append(insertions, BalanceUpdate{
				height:          height,
				asset:           asset.id,
				address:         address,
				tx_version:      tx_version,
				tx_type:         tx_type,
				tx:              tx_id,
				balance:         after,
				asset_decimals:  asset.decimals,
				asset_name:      asset.name,
				is_inbound:      is_inbound,
				is_state_update: false,
				timestamp:       timestamp,
				block_timestamp: block_timestamp,
			})
		}
	}

	return insertions

}

func getChInput(ch chan *blockchain_updates.GetBlockUpdateResponse) proto.Input {
	var (
		height          proto.ColUInt32
		asset           proto.ColStr
		address         proto.ColStr
		tx_version      proto.ColUInt8
		tx_type         proto.ColUInt8
		tx              proto.ColStr
		balance         proto.ColUInt64
		asset_decimals  proto.ColUInt8
		asset_name      proto.ColStr
		is_inbound      proto.ColBool
		is_state_update proto.ColBool
		timestamp       = new(proto.ColDateTime64).WithPrecision(proto.PrecisionMilli)
		block_timestamp = new(proto.ColDateTime64).WithPrecision(proto.PrecisionMilli)
	)

	for u := range ch {
		for _, ins := range getInsertion(u) {
			height.Append(ins.height)
			asset.Append(ins.asset)
			address.Append(ins.address)
			tx_version.Append(ins.tx_version)
			tx_type.Append(ins.tx_type)
			tx.Append(ins.tx)
			balance.Append(ins.balance)
			asset_decimals.Append(ins.asset_decimals)
			asset_name.Append(ins.asset_name)
			is_inbound.Append(ins.is_inbound)
			is_state_update.Append(ins.is_state_update)
			timestamp.Append(ins.timestamp)
			block_timestamp.Append(ins.block_timestamp)
		}
	}

	return proto.Input{
		{Name: "height", Data: &height},
		{Name: "asset", Data: &asset},
		{Name: "address", Data: &address},
		{Name: "tx_version", Data: &tx_version},
		{Name: "tx_type", Data: &tx_type},
		{Name: "tx", Data: &tx},
		{Name: "balance", Data: &balance},
		{Name: "asset_decimals", Data: &asset_decimals},
		{Name: "asset_name", Data: &asset_name},
		{Name: "is_inbound", Data: &is_inbound},
		{Name: "is_state_update", Data: &is_state_update},
		{Name: "timestamp", Data: timestamp},
		{Name: "block_timestamp", Data: block_timestamp},
	}
}

// []*blockchain_updates.GetBlockUpdateResponse
func getBlocks(client blockchain_updates.BlockchainUpdatesApiClient, ch chan *blockchain_updates.GetBlockUpdateResponse, from int32, len int32) {
	var wg sync.WaitGroup
	toHeight := from + len - 1
	for height := from; height <= toHeight; height++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				resp, err := client.GetBlockUpdate(
					context.Background(),
					&blockchain_updates.GetBlockUpdateRequest{Height: height},
				)
				if err == nil {
					// log.Printf("Height:\t%d", resp.Update.Height)
					ch <- resp
					return
				}
				log.Panicln(err)
				log.Printf("Error: %v", err)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	close(ch)
}

func getCurrentHeight(url string) int32 {
	r, err := (&http.Client{}).Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var target BlocksHeightResponse

	json.NewDecoder(r.Body).Decode(&target)

	return target.Height - 1
}

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// Получаем параметры для gRPC-сервера
	grpcServerHostname := os.Getenv("GRPC_HOSTNAME")
	grpcServerPort := os.Getenv("GRPC_PORT")
	grpcServerHost := grpcServerHostname + ":" + grpcServerPort
	restApiHost := os.Getenv("REST_API_HOST")

	// Устанавливаем соединение с gRPC-сервером
	conn, err := grpc.NewClient(
		grpcServerHost,
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
	defer conn.Close()

	dbServerHostname := os.Getenv("CH_HOSTNAME")
	dbServerPort := os.Getenv("CH_PORT")
	dbServerHost := dbServerHostname + ":" + dbServerPort
	dbUser := os.Getenv("CH_USER")
	dbPassword := os.Getenv("CH_PASSWORD")
	dbDatabase := os.Getenv("CH_DATABASE")

	db_ctx := context.Background()
	db_conn, err := ch.Dial(
		db_ctx,
		ch.Options{
			Address:  dbServerHost,
			Database: dbDatabase,
			User:     dbUser,
			Password: dbPassword,
		},
	)
	if err != nil {
		panic(err)
	}

	client := blockchain_updates.NewBlockchainUpdatesApiClient(conn)

	// probe
	// consume_ctx := context.Background()
	// probe_ctx, cancel_probe_ctx := context.WithTimeout(consume_ctx, 1*time.Second)
	if _, err := client.GetBlockUpdate(
		context.Background(),
		&blockchain_updates.GetBlockUpdateRequest{Height: 1},
	); err != nil {
		// cancel_probe_ctx()
		panic(err)
	} else {
		log.Println("Probe OK")
	}

	if err := db_conn.Do(db_ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbDatabase + `.balances_updates
		(
			height          UInt32,
			asset           String,
			address         String,
			tx_version      UInt8,
			tx_type         UInt8,
			tx              String,
			balance         UInt64,
			asset_decimals  UInt8,
			asset_name      String,
			is_inbound      Bool,
			is_state_update Bool,
			timestamp       DateTime64(3),
			block_timestamp DateTime64(3)
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (height, asset, address, tx)
		`,
	}); err != nil {
		panic(err)
	}

	if err := db_conn.Do(db_ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbDatabase + `.blocks
		(
			height          UInt32,
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (height)
		`,
	}); err != nil {
		panic(err)
	}

	defer db_conn.Close()

	var maxHeightResult proto.ColUInt32

	if err := db_conn.Do(db_ctx, ch.Query{
		Body:   "SELECT max(height) FROM blocks",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	}); err != nil {
		panic(err)
	}

	var counter int32 = int32(maxHeightResult[0])

	log.Printf("Consume from: %d", counter)

	toHeight := getCurrentHeight(restApiHost + "/blocks/height")

	estimatedBlocks := toHeight - counter
	batchCount := int(estimatedBlocks / BLOCKS_STEP)
	var blocksStep int32 = BLOCKS_STEP

	if estimatedBlocks < blocksStep {
		blocksStep = estimatedBlocks
	}

	var ch_wg sync.WaitGroup

	start := time.Now()
	workTime := time.Duration(0)
	for i := 1; i <= batchCount; i++ {
		from := counter + 1
		start := time.Now()
		results := make(chan *blockchain_updates.GetBlockUpdateResponse, blocksStep)
		getBlocks(client, results, from, blocksStep)

		insertions := getChInput(results)

		ch_wg.Add(1)
		go func() {
			defer ch_wg.Done()

			if err := db_conn.Do(db_ctx, ch.Query{
				Body:  "INSERT INTO blocks VALUES",
				Input: proto.Input{insertions[0]},
			}); err != nil {
				panic(err)
			}

			if err := db_conn.Do(db_ctx, ch.Query{
				Body:  "INSERT INTO balances_updates VALUES",
				Input: insertions,
			}); err != nil {
				panic(err)
			}
		}()
		counter = counter + BLOCKS_STEP
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

	ch_wg.Wait()

	log.Printf("Total: %.2fs", time.Since(start).Seconds())
}
