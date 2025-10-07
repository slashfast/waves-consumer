package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs9Batch struct {
	*TxsBatch
	LeaseTxUID proto.ColInt64
}

func NewTxs9Batch() *Txs9Batch {
	return &Txs9Batch{
		TxsBatch:   NewTxsBatch(),
		LeaseTxUID: proto.ColInt64{},
	}
}

func (t *Txs9Batch) Append(txs *Txs9) {
	t.TxsBatch.Append(txs.Txs)
	t.LeaseTxUID.AppendArr(txs.LeaseTxUID)
}
