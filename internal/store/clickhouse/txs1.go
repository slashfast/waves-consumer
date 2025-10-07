package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs1Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs1Batch() *Txs1Batch {
	return &Txs1Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs1Batch) Append(txs *Txs1) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.Amount.AppendArr(txs.Amount)
}
