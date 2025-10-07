package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs18Batch struct {
	*TxsBatch
	Bytes        proto.ColStr
	FunctionName proto.ColStr
}

func NewTxs18Batch() *Txs18Batch {
	return &Txs18Batch{
		TxsBatch:     NewTxsBatch(),
		Bytes:        proto.ColStr{},
		FunctionName: proto.ColStr{},
	}
}

func (t *Txs18Batch) Append(txs *Txs18) {
	t.TxsBatch.Append(txs.Txs)
	t.Bytes.AppendArr(txs.Bytes)
	t.FunctionName.AppendArr(txs.FunctionName)
}
