package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs15Batch struct {
	*TxsBatch
	AssetId proto.ColStr
	Script  proto.ColStr
}

func NewTxs15Batch() *Txs15Batch {
	return &Txs15Batch{
		TxsBatch: NewTxsBatch(),
		AssetId:  proto.ColStr{},
		Script:   proto.ColStr{},
	}
}

func (t *Txs15Batch) Append(txs *Txs15) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Script.AppendArr(txs.Script)
}
