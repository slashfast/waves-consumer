package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs13Batch struct {
	*TxsBatch
	Script proto.ColStr
}

func NewTxs13Batch() *Txs13Batch {
	return &Txs13Batch{
		TxsBatch: NewTxsBatch(),
		Script:   proto.ColStr{},
	}
}

func (t *Txs13Batch) Append(txs *Txs13) {
	t.TxsBatch.Append(txs.Txs)
	t.Script.AppendArr(txs.Script)
}
