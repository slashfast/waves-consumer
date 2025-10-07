package store

import "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"

type Update *events.BlockchainUpdated
type UpdatesChannel = chan Update

type Connector interface {
	Height() (int32, error)
	StartLoader(UpdatesChannel, chan int32, int32)
}
