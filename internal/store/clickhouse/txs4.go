package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs4Batch struct {
	*TxsBatch
	AssetId          proto.ColStr
	Amount           proto.ColInt64
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	FeeAssetId       proto.ColStr
	Attachment       proto.ColStr
}

func NewTxs4Batch() *Txs4Batch {
	return &Txs4Batch{
		TxsBatch:         NewTxsBatch(),
		AssetId:          proto.ColStr{},
		Amount:           proto.ColInt64{},
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		FeeAssetId:       proto.ColStr{},
		Attachment:       proto.ColStr{},
	}
}

func (t *Txs4Batch) Append(txs *Txs4) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Amount.AppendArr(txs.Amount)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.FeeAssetId.AppendArr(txs.FeeAssetId)
	t.Attachment.AppendArr(txs.Attachment)
}
