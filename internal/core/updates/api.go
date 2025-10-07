package updates

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/rs/zerolog"
	bu "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/slashfast/wavesconsumer/v2/internal/store"
)

type API struct {
	ctx    context.Context
	client bu.BlockchainUpdatesApiClient
	conn   *grpc.ClientConn
	log    zerolog.Logger
	cfg    Config
}

func NewAPI(ctx context.Context, logger *zerolog.Logger, config Config) (*API, error) {
	var creds credentials.TransportCredentials
	if config.Secure {
		creds = credentials.NewTLS(&tls.Config{
			MinVersion: tls.VersionTLS12,
		})
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(8*1024*1024)),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{"service": "waves.events.grpc.BlockchainUpdatesAPI"}],
				"retryPolicy": {
					"maxAttempts": 4,
					"initialBackoff": "0.1s",
					"maxBackoff": "5s",
					"backoffMultiplier": 2,
					"retryableStatusCodes": [
					    "UNAVAILABLE",
					    "DEADLINE_EXCEEDED",
                        "RESOURCE_EXHAUSTED"
				    ]
				}
			}]
		}`),
	)
	if err != nil {
		return nil, err
	}

	client := bu.NewBlockchainUpdatesApiClient(conn)
	newLogger := logger.With().Str("instance", "updates_api").Logger()

	return &API{
		ctx,
		client,
		conn,
		newLogger,
		config,
	}, nil
}

func Provider(ctx context.Context, lc fx.Lifecycle, logger *zerolog.Logger, configProvider ConfigProvider) (*API, error) {
	cfg := NewConfigFromProvider(configProvider)
	api, err := NewAPI(ctx, logger, cfg)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			request := &bu.GetBlockUpdateRequest{Height: 1}
			_, err := api.client.GetBlockUpdate(ctx, request)
			if err != nil {
				return err
			}
			api.log.Info().Msg("blockchain node gRPC test is OK")
			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := api.Close(); err != nil {
				return err
			}
			api.log.Info().Msg("gRPC to the blockchain node is closed")
			return nil
		},
	})

	return api, nil
}

func (a *API) StartConsumer(outCh chan<- store.Update, heightCh <-chan int32) {
	wg := &sync.WaitGroup{}
	wg.Add(int(a.cfg.MaxWorkers))

	for range a.cfg.MaxWorkers {
		go a.worker(a.ctx, wg, heightCh, outCh)
	}

	go func() {
		defer close(outCh)
		wg.Wait()
	}()
}

func (a *API) worker(ctx context.Context, wg *sync.WaitGroup, heightCh <-chan int32, outCh chan<- store.Update) {
	defer wg.Done()

	for {
		select {
		case height := <-heightCh:
			req := &bu.GetBlockUpdateRequest{Height: height}
			resp, err := a.client.GetBlockUpdate(ctx, req)
			if err != nil {
				if s, ok := status.FromError(err); ok {
					switch s.Code() {
					case codes.Canceled:
						return
					default:
						a.log.Fatal().
							Err(err).
							Int32("height", height).
							Str("details", s.String()).
							Msg("Failed to get the block update")
					}
				}
			}

			outCh <- resp.Update
		case <-ctx.Done():
			return
		}
	}
}

func (a *API) NewUpdatesChannel() chan store.Update {
	return make(store.UpdatesChannel)
}

func (a *API) Close() error {
	return a.conn.Close()
}
