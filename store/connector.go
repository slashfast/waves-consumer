package store

import "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"

type Update *events.BlockchainUpdated
type UpdatesChannel chan Update

type Connector interface {
	GetHeight() (int32, error)
	RunLoader(UpdatesChannel, chan int32, int32)
}
