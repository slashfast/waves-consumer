package main

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

type NodeApiConfig interface {
	NodeApiAddress() string
	IsNodeApiSecure() bool
}
type NodeApi struct {
	ctx          context.Context
	blocksClient node.BlocksApiClient
	log          zerolog.Logger
}

func NewNodeApi(ctx context.Context, lc fx.Lifecycle, l *zerolog.Logger, config NodeApiConfig) (*NodeApi, error) {
	var creds credentials.TransportCredentials
	if config.IsNodeApiSecure() {
		creds = credentials.NewTLS(&tls.Config{})
	} else {
		creds = insecure.NewCredentials()
	}
	conn, err := grpc.NewClient(
		config.NodeApiAddress(),
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

	blocksClient := node.NewBlocksApiClient(conn)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if _, err := blocksClient.GetCurrentHeight(ctx, &emptypb.Empty{}); err != nil {
				return err
			}
			l.Info().Msg("Node GRPC probe is OK")
			return nil
		},
		OnStop: func(_ context.Context) error {
			err := conn.Close()
			if err != nil {
				l.Info().Msg("Closed Node GRPC connection")
			}
			return err
		},
	})
	return &NodeApi{
		ctx,
		blocksClient,
		l.With().Str("instance", "node_api").Logger(),
	}, nil
}

func (n *NodeApi) getCurrentHeight() (int32, error) {
	resp, err := n.blocksClient.GetCurrentHeight(n.ctx, &emptypb.Empty{})
	if err != nil {
		return 0, err
	}
	return int32(resp.Value), nil
}

func (n *NodeApi) GetCurrentHeight() int32 {
	height, err := n.getCurrentHeight()
	if err != nil {
		n.log.Error().Err(err).Msg("Failed to get current height")
		return height
	}
	return height
}

func (n *NodeApi) PollCurrentHeight(ch chan int32) {
	logger := n.log.With().Str("worker", "height_poller").Logger()
	for {
		select {
		case <-n.ctx.Done():
			return
		default:
			height, err := n.getCurrentHeight()
			if err == nil {
				logger.Debug().Int32("current_height", height).Send()
				ch <- height - 1
			}
		}
	}
}
