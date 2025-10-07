package node

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	node "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/node/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type API struct {
	ctx          context.Context
	blocksClient node.BlocksApiClient
	conn         *grpc.ClientConn
	log          zerolog.Logger
	cfg          Config
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
				"name": [{"service": "waves.events.node.grpc.BlocksAPI"}],
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
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	blocksClient := node.NewBlocksApiClient(conn)
	newLogger := logger.With().Str("instance", "node_api").Logger()

	return &API{
		ctx,
		blocksClient,
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
			if _, err := api.blocksClient.GetCurrentHeight(ctx, &emptypb.Empty{}); err != nil {
				return fmt.Errorf("failed to probe node GRPC: %w", err)
			}
			api.log.Info().Msg("Node GRPC probe is OK")
			return nil
		},
		OnStop: func(_ context.Context) error {
			err := api.Close()
			if err != nil {
				api.log.Info().Msg("Closed Node GRPC connection")
			}
			return err
		},
	})

	return api, nil
}

func (a *API) CurrentHeight() int32 {
	height, err := a.currentHeight()
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get current height")
		return height
	}
	return height
}

func (a *API) currentHeight() (int32, error) {
	resp, err := a.blocksClient.GetCurrentHeight(a.ctx, &emptypb.Empty{})
	if err != nil {
		return 0, fmt.Errorf("failed to get current height: %w", err)
	}
	// #nosec G115 -- the height of the blockchain will never exceed int32, overflow is impossible
	return int32(resp.Value), nil
}

func (a *API) StartHeightPoller(output chan int32) {
	logger := a.log.With().Str("worker", "height_poller").Logger()
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			height, err := a.currentHeight()
			if err == nil {
				logger.Debug().Int32("current_height", height).Send()
				output <- height - 1
			}
		}
	}
}

func (a *API) StartHeightProvider(outputCh chan int32, start int32) {
	logger := a.log.With().Str("worker", "height_consumer").Logger()

	currentHeight := a.CurrentHeight()
	accumulatorHeight := start

	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			delta := currentHeight - accumulatorHeight

			if delta != 0 {
				logger.Trace().
					Int32("current", currentHeight).
					Int32("accumulator", accumulatorHeight).
					Int32("delta", delta).
					Msg("sending block height")

				outputCh <- accumulatorHeight
				accumulatorHeight++
				continue
			}

			// checks if there are new blocks
			currentHeight = a.CurrentHeight()
			delta = currentHeight - accumulatorHeight

			if delta != 0 {
				continue
			}

			logger.Debug().
				Int32("current_height", currentHeight).
				Msg("No new blocks, waiting")

			// wait for a minute or context cancellation
			select {
			case <-time.After(a.cfg.NewHeightDelay):
			case <-a.ctx.Done():
				return
			}
		}
	}
}

func (a *API) Close() error {
	err := a.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}
	return nil
}
