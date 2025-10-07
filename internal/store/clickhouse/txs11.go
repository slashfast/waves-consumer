package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs11Batch struct {
	*TxsBatch
	AssetID    proto.ColStr
	Attachment proto.ColStr
}

func NewTxs11Batch() *Txs11Batch {
	return &Txs11Batch{
		TxsBatch:   NewTxsBatch(),
		AssetID:    proto.ColStr{},
		Attachment: proto.ColStr{},
	}
}

func (t *Txs11Batch) Append(txs *Txs11) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.Attachment.AppendArr(txs.Attachment)
}

type Txs11TransferBatch struct {
	TxUID            proto.ColInt64
	TxID             proto.ColStr
	PositionInTx     proto.ColInt32
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	Amount           proto.ColInt64
	Height           proto.ColInt32
}

func NewTxs11TransferBatch() *Txs11TransferBatch {
	return &Txs11TransferBatch{
		TxUID:            proto.ColInt64{},
		TxID:             proto.ColStr{},
		PositionInTx:     proto.ColInt32{},
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		Amount:           proto.ColInt64{},
		Height:           proto.ColInt32{},
	}
}

func (t *Txs11TransferBatch) Append(txs *Txs11Transfer) {
	t.TxUID.AppendArr(txs.TxUID)
	t.TxID.AppendArr(txs.TxID)
	t.PositionInTx.AppendArr(txs.PositionInTx)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.Amount.AppendArr(txs.Amount)
	t.Height.AppendArr(txs.Height)
}
