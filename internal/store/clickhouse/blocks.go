package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type BlocksBatch struct {
	Id        proto.ColStr
	Timestamp proto.ColDateTime64
	Height    proto.ColInt32
}

func NewBlocksBatch() *BlocksBatch {
	return &BlocksBatch{
		Height:    proto.ColInt32{},
		Id:        proto.ColStr{},
		Timestamp: proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true},
	}
}
