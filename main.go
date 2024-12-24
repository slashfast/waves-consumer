package main

import (
	"context"
	"github.com/slashfast/wavesconsumer/v2/connectors/clickhouse"
	"github.com/slashfast/wavesconsumer/v2/waves"
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/joho/godotenv"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

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

func main() {
	// Загружаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Println("The dotenv file is missing")
	}

	waves.RestApiHost = os.Getenv("REST_API_HOST")

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
		log.Fatal(err)
	}

	// Закрываем соединение после завершения работы
	defer func() {
		if err := grpcConn.Close(); err != nil {
			log.Fatal(err)
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
		log.Fatal(err)
	}

	if err := createTables(dbConn, dbCtx, dbName); err != nil {
		log.Fatal(err)
	}

	client := blockchainupdates.NewBlockchainUpdatesApiClient(grpcConn)

	// Проверка соединения
	if _, err := client.GetBlockUpdate(
		context.Background(),
		&blockchainupdates.GetBlockUpdateRequest{Height: 1},
	); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Probe OK")
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Получаем начальную максимальную высоту из базы данных
	var maxHeightResult proto.ColUInt32
	if err := dbConn.Do(dbCtx, ch.Query{
		Body:   "SELECT max(height) FROM balances_updates",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	}); err != nil {
		log.Fatal(err)
	}
	currentHeight := int32(maxHeightResult[0]) // Начальная высота из базы данных
	log.Printf("Starting from height: %d", currentHeight)

	newBlockchainHeight := waves.GetCurrentHeight() // Запрос текущей высоты блокчейна

	var heightDiff int32
	blocksStep := BlocksStep
	var prevBlockStep int32 = 0 // Переменная для хранения предыдущей разницы

	for {
		heightDiff = newBlockchainHeight - currentHeight
		realTimeMode := !(heightDiff >= BlocksStep)

		if realTimeMode {
			newBlockchainHeight = waves.GetCurrentHeight() // Запрос текущей высоты блокчейна
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

		updates := waves.GetUpdates(results)
		insertions := clickhouse.GetInput(updates)
		if err := clickhouse.InsertToDB(dbConn, insertions); err != nil {
			log.Fatal(err)
		}

		// Обновляем текущую высоту
		currentHeight += blocksStep
		log.Printf("Updated current height to: %d", currentHeight)

		prevBlockStep = blocksStep
	}
}
