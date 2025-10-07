package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs6Batch struct {
	*TxsBatch
	AssetID proto.ColStr
	Amount  proto.ColInt64
}

func NewTxs6Batch() *Txs6Batch {
	return &Txs6Batch{
		TxsBatch: NewTxsBatch(),
		AssetID:  proto.ColStr{},
		Amount:   proto.ColInt64{},
	}
}

func (t *Txs6Batch) Append(txs *Txs6) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.Amount.AppendArr(txs.Amount)
}
