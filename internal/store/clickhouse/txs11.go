package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs11Batch struct {
	*TxsBatch
	AssetId    proto.ColStr
	Attachment proto.ColStr
}

func NewTxs11Batch() *Txs11Batch {
	return &Txs11Batch{
		TxsBatch:   NewTxsBatch(),
		AssetId:    proto.ColStr{},
		Attachment: proto.ColStr{},
	}
}

func (t *Txs11Batch) Append(txs *Txs11) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Attachment.AppendArr(txs.Attachment)
}

type Txs11TransferBatch struct {
	TxUid            proto.ColInt64
	TxId             proto.ColStr
	PositionInTx     proto.ColInt32
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	Amount           proto.ColInt64
	Height           proto.ColInt32
}

func NewTxs11TransferBatch() *Txs11TransferBatch {
	return &Txs11TransferBatch{
		TxUid:            proto.ColInt64{},
		TxId:             proto.ColStr{},
		PositionInTx:     proto.ColInt32{},
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		Amount:           proto.ColInt64{},
		Height:           proto.ColInt32{},
	}
}

func (t *Txs11TransferBatch) Append(txs *Txs11Transfer) {
	t.TxUid.AppendArr(txs.TxUid)
	t.TxId.AppendArr(txs.TxId)
	t.PositionInTx.AppendArr(txs.PositionInTx)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.Amount.AppendArr(txs.Amount)
	t.Height.AppendArr(txs.Height)
}
