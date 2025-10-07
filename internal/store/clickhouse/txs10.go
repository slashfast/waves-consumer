package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs10Batch struct {
	*TxsBatch
	Alias proto.ColStr
}

func NewTxs10Batch() *Txs10Batch {
	return &Txs10Batch{
		TxsBatch: NewTxsBatch(),
		Alias:    proto.ColStr{},
	}
}

func (t *Txs10Batch) Append(txs *Txs10) {
	t.TxsBatch.Append(txs.Txs)
	t.Alias.AppendArr(txs.Alias)
}
