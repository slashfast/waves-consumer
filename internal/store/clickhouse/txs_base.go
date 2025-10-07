package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type TxsBatch struct {
	UID             proto.ColInt64
	Type            proto.ColInt8
	Sender          proto.ColStr
	SenderPublicKey proto.ColStr
	Timestamp       proto.ColDateTime64
	Height          proto.ColInt32
	ID              proto.ColStr
	Signature       proto.ColStr
	Proofs          *proto.ColArr[string]
	Version         proto.ColInt32
	Fee             proto.ColInt64
	Status          proto.ColStr
}

func NewTxsBatch() *TxsBatch {
	return &TxsBatch{
		UID:             proto.ColInt64{},
		Type:            proto.ColInt8{},
		Sender:          proto.ColStr{},
		SenderPublicKey: proto.ColStr{},
		Timestamp:       proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true},
		Height:          proto.ColInt32{},
		ID:              proto.ColStr{},
		Signature:       proto.ColStr{},
		Proofs:          new(proto.ColStr).Array(),
		Version:         proto.ColInt32{},
		Fee:             proto.ColInt64{},
		Status:          proto.ColStr{},
	}
}

func (t *TxsBatch) Append(txs *Txs) {
	t.UID.AppendArr(txs.UID)
	t.Type.AppendArr(txs.Type)
	t.Sender.AppendArr(txs.Sender)
	t.SenderPublicKey.AppendArr(txs.SenderPublicKey)
	t.Timestamp.AppendArr(txs.Timestamp)
	t.Height.AppendArr(txs.Height)
	t.ID.AppendArr(txs.ID)
	t.Signature.AppendArr(txs.Signature)
	t.Proofs.AppendArr(txs.Proofs)
	t.Version.AppendArr(txs.Version)
	t.Fee.AppendArr(txs.Fee)
	t.Status.AppendArr(txs.Status)
}
