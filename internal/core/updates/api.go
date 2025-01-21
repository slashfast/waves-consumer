package updates

import "C"
import (
	"context"
	"crypto/tls"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/internal/store"
	bu "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"sync"
	"sync/atomic"
	"time"
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
			err := api.Close()
			if err != nil {
				api.log.Info().Msg("Updates GRPC connection is closed")
			}
			return err
		},
	})

	return api, nil
}

func (a *Api) runHeightBufferFiller(ch chan int32, start int32) {
	defer close(ch)

	current := start

	for {
		select {
		case <-a.ctx.Done():
			return
		case ch <- current:
			current++
		}
	}
}

func (a *Api) StartConsumer(updatesCh store.UpdatesChannel, currentHeightCh chan int32, start int32) {
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		close(updatesCh)
	}()

	sem := make(chan struct{}, a.cfg.MaxWorkers)
	defer close(sem)

	bufferSize := int32(a.cfg.BufferSize)
	afterSyncDelay := a.cfg.AfterSyncDelay

	heightBuffer := make(chan int32, bufferSize)
	go a.runHeightBufferFiller(heightBuffer, start)

	targetHeight := <-currentHeightCh
	logger := a.log.With().Str("worker", "consumer").Logger()
	logger.Info().Int32("target_height", targetHeight).Send()

	var (
		previousConsumerHeight,
		currentConsumerHeight int32
		consumerHeight atomic.Int32
	)
	consumerHeight.Store(start)

	worker := func(height int32) {
		wg.Add(1)
		defer func() { <-sem }()
		defer wg.Done()

		resp, err := a.client.GetBlockUpdate(
			a.ctx,
			&bu.GetBlockUpdateRequest{Height: height},
		)

		if err == nil {
			updatesCh <- resp.Update
			if consumerHeight.Load() < height {
				consumerHeight.Store(height)
			}
		} else {
			st, ok := status.FromError(err)
			if ok {
				switch st.Code() {
				case codes.Canceled:
				default:
					logger.Fatal().
						Err(err).
						Int32("height", height).
						Str("details", st.String()).
						Msg("Failed to get block update")
				}
			}
		}
	}

	for {
		heightDelta := targetHeight - currentConsumerHeight
		isQuietMode := !(heightDelta >= bufferSize)

		select {
		case <-a.ctx.Done():
			return
		default:
			if isQuietMode {
				logger.Info().
					Int32("target_height", targetHeight).
					Int32("previous_height", previousConsumerHeight).
					Int32("current_height", currentConsumerHeight).
					Int32("height_delta", heightDelta).
					Send()
				h := currentConsumerHeight
				for range heightDelta {
					h++
					sem <- struct{}{}
					worker(h)
					logger.Info().Int32("consumed_height", h).Send()
				}
			} else {
				sem <- struct{}{}
				go worker(<-heightBuffer)
			}
		}

		previousConsumerHeight = currentConsumerHeight
		currentConsumerHeight = consumerHeight.Load()

		if !isQuietMode &&
			previousConsumerHeight < currentConsumerHeight &&
			currentConsumerHeight%bufferSize == 0 {
			logger.Info().
				Int32("consumed_height", currentConsumerHeight).
				Send()
		}

		if heightDelta == 0 {
			logger.Info().Msg("Synced")
		DelayLoop:
			for {
				select {
				case <-a.ctx.Done():
					return
				case <-time.After(afterSyncDelay):
					break DelayLoop
				}
			}
		}

		if isQuietMode {
			targetHeight = <-currentHeightCh
		}
	}
}

func (a *Api) NewUpdatesChannel() store.UpdatesChannel {
	return make(store.UpdatesChannel, a.cfg.BufferSize)
}

func (a *Api) Close() error {
	return a.conn.Close()
}
