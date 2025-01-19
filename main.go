package main

import (
	"context"
	fxzerolog "github.com/efectn/fx-zerolog"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/config"
	"github.com/slashfast/wavesconsumer/v2/store"
	"github.com/slashfast/wavesconsumer/v2/store/clickhouse"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	"go.uber.org/fx"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

type App struct {
	ctx     context.Context
	log     zerolog.Logger
	updates *UpdatesApi
	node    *NodeApi
	store   store.Connector
}

func NewApp(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zerolog.Logger,
	updates *UpdatesApi,
	node *NodeApi,
	store store.Connector,
) *App {
	app := &App{ctx, l.With().Str("instance", "app").Logger(), updates, node, store}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.run()
			return nil
		},
	})
	return app
}

func (app App) run() {
	dbHeight, err := app.store.GetHeight()
	if err != nil {
		app.log.Fatal().Err(err).Msg("Failed to get current database height")
	}
	app.log.Info().Int32("from_height", dbHeight).Msg("Running")

	startHeight := dbHeight + 1
	updatesCh := make(chan *events.BlockchainUpdated, app.updates.config.GetUpdatesBufferSize())
	currentHeight := make(chan int32)

	go app.node.PollCurrentHeight(currentHeight)
	go app.updates.Consume(updatesCh, currentHeight, startHeight)
	go app.store.RunLoader(updatesCh, currentHeight, startHeight)

	<-app.ctx.Done()
}

func NewLogger(lc fx.Lifecycle) (*zerolog.Logger, error) {
	//wr := diode.NewWriter(os.Stderr, 1000, 0, func(missed int) {
	//	log.Warn().Msgf("Logger dropped %d messages", missed)
	//})
	wr := os.Stderr

	err := os.MkdirAll("logs", 0744)
	if err != nil {
		return nil, err
	}

	lf, err := os.OpenFile("logs/"+time.Now().Format(time.RFC3339)+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := lf.Close(); err != nil {
				return err
			}
			//return
			return nil
		},
	})

	logger := zerolog.
		New(io.MultiWriter(wr, lf)).
		With().
		Timestamp().
		Logger()

	return &logger, nil
}

func NewContext(lc fx.Lifecycle) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			cancel()
			return nil
		},
	})
	return ctx
}

func main() {
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	fx.New(
		fx.Invoke(func() {
			zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		}),
		fx.Provide(
			NewApp,
			fx.Annotate(
				clickhouse.NewConnector,
				fx.As(new(store.Connector)),
			),
			NewUpdatesApi,
			NewNodeApi,
			fx.Annotate(
				config.LoadConfigFromEnv,
				fx.As(new(NodeApiConfig)),
				fx.As(new(UpdatesApiConfig)),
				fx.As(new(clickhouse.ConnectorConfig)),
			),
			NewLogger,
			NewContext,
		),
		fx.Invoke(func(a *App) {}),
		fx.WithLogger(fxzerolog.InitPtr()),
	).Run()
}
