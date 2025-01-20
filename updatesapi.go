package main

import (
	"context"
	"crypto/tls"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/store"
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

type UpdatesApiConfig interface {
	UpdatesApiAddress() string
	IsUpdatesApiSecure() bool
	UpdatesBufferSize() int
	UpdatesMaxWorkers() int
	UpdatesAfterSyncDelay() time.Duration
}

type UpdatesApi struct {
	ctx    context.Context
	client bu.BlockchainUpdatesApiClient
	log    zerolog.Logger
	config UpdatesApiConfig
}

func NewUpdatesApi(
	ctx context.Context,
	lc fx.Lifecycle,
	l *zerolog.Logger,
	config UpdatesApiConfig,
) (*UpdatesApi, error) {
	var creds credentials.TransportCredentials
	if config.IsUpdatesApiSecure() {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		creds = insecure.NewCredentials()
	}

	conn, err := grpc.NewClient(
		config.UpdatesApiAddress(),
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

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if _, err := client.GetBlockUpdate(
				ctx,
				&bu.GetBlockUpdateRequest{Height: 1},
			); err != nil {
				return err
			}
			l.Info().Msg("Updates GRPC probe is OK")
			return nil
		},
		OnStop: func(_ context.Context) error {
			err := conn.Close()
			if err != nil {
				l.Info().Msg("Updates GRPC connection is closed")
			}
			return err
		},
	})

	return &UpdatesApi{
		ctx,
		client,
		l.With().Str("instance", "updates_api").Logger(),
		config,
	}, nil
}

func (u *UpdatesApi) runHeightBufferFiller(ch chan int32, start int32) {
	defer close(ch)

	current := start

	for {
		select {
		case <-u.ctx.Done():
			return
		case ch <- current:
			current++
		}
	}
}

func (u *UpdatesApi) Consume(updatesCh store.UpdatesChannel, currentHeightCh chan int32, start int32) {
	wg := sync.WaitGroup{}
	defer func() {
		wg.Wait()
		close(updatesCh)
	}()

	sem := make(chan struct{}, u.config.UpdatesMaxWorkers())
	defer close(sem)

	bufferSize := int32(u.config.UpdatesBufferSize())
	afterSyncDelay := u.config.UpdatesAfterSyncDelay()

	heightBuffer := make(chan int32, bufferSize)
	go u.runHeightBufferFiller(heightBuffer, start)

	targetHeight := <-currentHeightCh
	logger := u.log.With().Str("worker", "consumer").Logger()
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

		resp, err := u.client.GetBlockUpdate(
			u.ctx,
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
		case <-u.ctx.Done():
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
				case <-u.ctx.Done():
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
