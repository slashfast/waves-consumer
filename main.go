package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/waves"
	"github.com/slashfast/wavesconsumer/v2/waves/connectors"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"os"
	"sync"
	"time"
)

const NoNewBlocksTimeout = 1 * time.Minute
const BlocksStep int32 = 1000

type App struct {
	wg  *sync.WaitGroup
	log *zerolog.Logger
}

func (app App) isCtxDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		app.log.Info().Msg("Shutdown signal received")
		return true
	default:
		return false
	}
}

func (app App) serve(ctx context.Context, updatesClient *waves.UpdatesClient, conn connectors.Connector) {
	app.wg.Add(1)
	defer app.wg.Done()

	currentHeight, err := conn.GetHeight()
	if err != nil {
		app.log.Fatal().Err(err).Msg("Failed to get current height")
	}
	app.log.Info().Int32("from_height", currentHeight).Msg("Serving")

	newBlockchainHeight := updatesClient.GetCurrentHeight()

	var heightDiff int32
	blocksStep := BlocksStep
	var prevBlockStep int32 = 0 // Переменная для хранения предыдущей разницы
MainLoop:
	for {
		if app.isCtxDone(ctx) {
			return
		}
		heightDiff = newBlockchainHeight - currentHeight
		isRealTimeMode := !(heightDiff >= BlocksStep)

		if isRealTimeMode {
			newBlockchainHeight = updatesClient.GetCurrentHeight() // Запрос текущей высоты блокчейна
			heightDiff = newBlockchainHeight - currentHeight
			if heightDiff == 0 {
				app.log.Info().
					Int32("db", currentHeight).
					Int32("blockchain", newBlockchainHeight).
					Dur("next_check_in", NoNewBlocksTimeout).
					Msg("Synced")
				tick := time.Tick(1 * time.Second)
				var timeout <-chan time.Time
			NoNewBlockWaiting:
				for range tick {
					if timeout == nil {
						timeout = time.After(NoNewBlocksTimeout)
					}
					select {
					case <-timeout:
						break NoNewBlockWaiting
					default:
						if app.isCtxDone(ctx) {
							return
						}
					}
				}
				continue MainLoop
			} else if heightDiff < 0 {
				app.log.Fatal().Msg("Height difference is negative")
			}
			blocksStep = heightDiff
		}

		// Печатаем информацию только если разница отличается от предыдущей
		if prevBlockStep != blocksStep {
			app.log.Info().
				Int32("from", currentHeight).
				Int32("to", newBlockchainHeight).
				Int32("step", blocksStep).
				Msg("Sync")
		}

		results := make(chan *blockchainupdates.GetBlockUpdateResponse, blocksStep)
		updatesClient.GetBlocks(results, currentHeight+1, blocksStep)
		updates := updatesClient.GetUpdates(results)
		if err := conn.Insert(updates, !isRealTimeMode); err != nil {
			app.log.Fatal().Err(err).Msg("Failed to insert block updates")
		}

		// Обновляем текущую высоту
		currentHeight += blocksStep
		app.log.Info().Int32("to", currentHeight).Msg("Updated")

		prevBlockStep = blocksStep
	}
}

func NewApp(lc fx.Lifecycle, l *zerolog.Logger, uc *waves.UpdatesClient, conn connectors.Connector) *App {
	bg, cancel := context.WithCancel(context.Background())
	app := &App{wg: &sync.WaitGroup{}, log: l}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.serve(bg, uc, conn)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			app.wg.Wait()
			return nil
		},
	})
	return app
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
			NewApp,
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
		fx.Invoke(func(app *App) {}),
	).Run()
}
