package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs2Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs2Batch() *Txs2Batch {
	return &Txs2Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs2Batch) Append(txs *Txs2) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.Amount.AppendArr(txs.Amount)
}
