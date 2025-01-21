package node

import (
	"context"
	"crypto/tls"
	"github.com/rs/zerolog"
	node "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/node/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Api struct {
	ctx          context.Context
	blocksClient node.BlocksApiClient
	conn         *grpc.ClientConn
	log          zerolog.Logger
	cfg          Config
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
				"name": [{"service": "waves.events.node.grpc.BlocksApi"}],
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

	blocksClient := node.NewBlocksApiClient(conn)
	newLogger := logger.With().Str("instance", "node_api").Logger()

	return &Api{
		ctx,
		blocksClient,
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
			if _, err := api.blocksClient.GetCurrentHeight(ctx, &emptypb.Empty{}); err != nil {
				return err
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

func (a *Api) CurrentHeight() int32 {
	height, err := a.currentHeight()
	if err != nil {
		a.log.Error().Err(err).Msg("Failed to get current height")
		return height
	}
	return height
}

func (a *Api) currentHeight() (int32, error) {
	resp, err := a.blocksClient.GetCurrentHeight(a.ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}
	return int32(resp.Value), nil
}

func (a *Api) StartHeightPoller(ch chan int32) {
	logger := a.log.With().Str("worker", "height_poller").Logger()
	for {
		select {
		case <-a.ctx.Done():
			return
		default:
			height, err := a.currentHeight()
			if err == nil {
				logger.Debug().Int32("current_height", height).Send()
				ch <- height - 1
			}
		}
	}
}

func (a *Api) Close() error {
	return a.conn.Close()
}
