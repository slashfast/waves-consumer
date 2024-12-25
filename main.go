package main

import (
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/waves"
	"github.com/slashfast/wavesconsumer/v2/waves/connectors"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"os"
	"strconv"
	"time"
)

const BlocksStep int32 = 1000

func Serve(l *zerolog.Logger, updatesClient *waves.UpdatesClient, conn connectors.Connector) {
	currentHeight, err := conn.GetHeight()
	if err != nil {
		l.Fatal().Err(err).Msg("failed to get current height")
	}
	l.Info().Str("from_height", strconv.Itoa(int(currentHeight))).Msg("Start serving...")

	newBlockchainHeight := updatesClient.GetCurrentHeight()

	var heightDiff int32
	blocksStep := BlocksStep
	var prevBlockStep int32 = 0 // Переменная для хранения предыдущей разницы

	for {
		heightDiff = newBlockchainHeight - currentHeight
		realTimeMode := !(heightDiff >= BlocksStep)

		if realTimeMode {
			newBlockchainHeight = updatesClient.GetCurrentHeight() // Запрос текущей высоты блокчейна
			l.Info().
				Str("db", strconv.Itoa(int(currentHeight))).
				Str("blockchain", strconv.Itoa(int(newBlockchainHeight))).
				Msg("Height")
			heightDiff = newBlockchainHeight - currentHeight
			if heightDiff == 0 {
				l.Info().Msg("No new blocks. Waiting for 1 minute...")
				time.Sleep(1 * time.Minute)
				continue
			} else if heightDiff < 0 {
				l.Fatal().Msg("Height difference is negative! Possible data inconsistency.")
			}
			blocksStep = heightDiff
		}

		// Печатаем информацию только если разница отличается от предыдущей
		if prevBlockStep != blocksStep {
			l.Info().
				Str("from", strconv.Itoa(int(currentHeight))).
				Str("to", strconv.Itoa(int(newBlockchainHeight))).
				Str("step", strconv.Itoa(int(blocksStep))).
				Msg("Sync blocks")
		}

		results := make(chan *blockchainupdates.GetBlockUpdateResponse, blocksStep)
		updatesClient.GetBlocks(results, currentHeight+1, blocksStep)
		updates := updatesClient.GetUpdates(results)
		if err := conn.Insert(updates); err != nil {
			l.Fatal().Err(err).Msg("")
		}

		// Обновляем текущую высоту
		currentHeight += blocksStep
		l.Info().Str("to", strconv.Itoa(int(currentHeight))).Msg("Updated")

		prevBlockStep = blocksStep
	}
}

func NewLogger(_ fx.Lifecycle) (*zerolog.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.
		New(os.Stderr).
		With().
		Timestamp().
		Logger()
	return &logger, nil
}

func main() {
	fx.New(
		fx.Provide(
			fx.Annotate(
				connectors.NewConnector,
				fx.As(new(connectors.Connector)),
			),
			waves.NewUpdatesClient,
			fx.Annotate(
				LoadConfigFromEnv,
				fx.As(new(waves.NodeClientConfig)),
				fx.As(new(connectors.ClickHouseConfig)),
			),
			NewLogger,
		),
		fx.Invoke(Serve),
	).Run()
}
