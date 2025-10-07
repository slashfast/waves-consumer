package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs8Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs8Batch() *Txs8Batch {
	return &Txs8Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs8Batch) Append(txs *Txs8) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.Amount.AppendArr(txs.Amount)
}
