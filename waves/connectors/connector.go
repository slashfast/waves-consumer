package connectors

import "github.com/slashfast/wavesconsumer/v2/waves"

type Connector interface {
	Insert([]waves.BalanceUpdate) error
	GetHeight() (int32, error)
}
