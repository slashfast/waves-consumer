package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs5Batch struct {
	*TxsBatch
	AssetID    proto.ColStr
	Amount     proto.ColInt64
	Reissuable proto.ColBool
}

func NewTxs5Batch() *Txs5Batch {
	return &Txs5Batch{
		TxsBatch:   NewTxsBatch(),
		AssetID:    proto.ColStr{},
		Amount:     proto.ColInt64{},
		Reissuable: proto.ColBool{},
	}
}

func (t *Txs5Batch) Append(txs *Txs5) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.Amount.AppendArr(txs.Amount)
	t.Reissuable.AppendArr(txs.Reissuable)
}
