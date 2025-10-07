package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type TxsBatch struct {
	Uid             proto.ColInt64
	Type            proto.ColInt8
	Sender          proto.ColStr
	SenderPublicKey proto.ColStr
	Timestamp       proto.ColDateTime64
	Height          proto.ColInt32
	Id              proto.ColStr
	Signature       proto.ColStr
	Proofs          *proto.ColArr[string]
	Version         proto.ColInt32
	Fee             proto.ColInt64
	Status          proto.ColStr
}

func NewTxsBatch() *TxsBatch {
	return &TxsBatch{
		Uid:             proto.ColInt64{},
		Type:            proto.ColInt8{},
		Sender:          proto.ColStr{},
		SenderPublicKey: proto.ColStr{},
		Timestamp:       proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true},
		Height:          proto.ColInt32{},
		Id:              proto.ColStr{},
		Signature:       proto.ColStr{},
		Proofs:          new(proto.ColStr).Array(),
		Version:         proto.ColInt32{},
		Fee:             proto.ColInt64{},
		Status:          proto.ColStr{},
	}
}

func (t *TxsBatch) Append(txs *Txs) {
	t.Uid.AppendArr(txs.Uid)
	t.Type.AppendArr(txs.Type)
	t.Sender.AppendArr(txs.Sender)
	t.SenderPublicKey.AppendArr(txs.SenderPublicKey)
	t.Timestamp.AppendArr(txs.Timestamp)
	t.Height.AppendArr(txs.Height)
	t.Id.AppendArr(txs.Id)
	t.Signature.AppendArr(txs.Signature)
	t.Proofs.AppendArr(txs.Proofs)
	t.Version.AppendArr(txs.Version)
	t.Fee.AppendArr(txs.Fee)
	t.Status.AppendArr(txs.Status)
}
