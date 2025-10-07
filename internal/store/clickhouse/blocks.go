package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type BlocksBatch struct {
	ID        proto.ColStr
	Timestamp proto.ColDateTime64
	Height    proto.ColInt32
}

func NewBlocksBatch() *BlocksBatch {
	return &BlocksBatch{
		Height:    proto.ColInt32{},
		ID:        proto.ColStr{},
		Timestamp: proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true},
	}
}
