package updates

import "C"
import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/internal/store"
	bu "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Api struct {
	ctx    context.Context
	client bu.BlockchainUpdatesApiClient
	conn   *grpc.ClientConn
	log    zerolog.Logger
	cfg    Config
}

func NewApi(ctx context.Context, logger *zerolog.Logger, config Config) (*Api, error) {
	var creds credentials.TransportCredentials
	if config.Secure {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(
		config.Address,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(8*1024*1024)),
		grpc.WithDefaultServiceConfig(`{
			"loadBalancingConfig": [{"round_robin":{}}],
			"methodConfig": [{
				"name": [{"service": "waves.events.grpc.BlockchainUpdatesApi"}],
				"retryPolicy": {
					"maxAttempts": 4,
					"initialBackoff": "0.1s",
					"maxBackoff": "1s",
					"backoffMultiplier": 2,
					"retryableStatusCodes": ["UNAVAILABLE"]
				}
			}]
		}`),
	)
	if err != nil {
		return nil, err
	}

	client := bu.NewBlockchainUpdatesApiClient(conn)
	newLogger := logger.With().Str("instance", "updates_api").Logger()

	return &Api{
		ctx,
		client,
		conn,
		newLogger,
		config,
	}, nil
}

func Provider(ctx context.Context, lc fx.Lifecycle, logger *zerolog.Logger, configProvider ConfigProvider) (*Api, error) {
	cfg := NewConfigFromProvider(configProvider)
	api, err := NewApi(ctx, logger, cfg)
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
			api.log.Info().Msg("Updates GRPC probe is OK")
			return nil
		},
		OnStop: func(_ context.Context) error {
			if err := api.Close(); err != nil {
				return err
			}
			api.log.Info().Msg("Updates GRPC connection is closed")
			return nil
		},
	})

	return api, nil
}

func (a *Api) StartConsumer(outCh chan<- store.Update, heightCh <-chan int32) {
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

func (a *Api) worker(ctx context.Context, wg *sync.WaitGroup, heightCh <-chan int32, outCh chan<- store.Update) {
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

func (a *Api) NewUpdatesChannel() chan store.Update {
	return make(store.UpdatesChannel)
}

func (a *Api) Close() error {
	return a.conn.Close()
}
