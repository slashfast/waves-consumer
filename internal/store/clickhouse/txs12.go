package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs12Batch struct {
	*TxsBatch
}

func NewTxs12Batch() *Txs12Batch {
	return &Txs12Batch{
		TxsBatch: NewTxsBatch(),
	}
}

func (t *Txs12Batch) Append(txs *Txs12) {
	t.TxsBatch.Append(txs.Txs)
}

type Txs12DataBatch struct {
	TxUID            proto.ColInt64
	TxID             proto.ColStr
	PositionInTx     proto.ColInt32
	DataKey          proto.ColStr
	DataType         proto.ColStr
	DataValueInteger proto.ColInt64
	DataValueBoolean proto.ColBool
	DataValueBinary  proto.ColStr
	DataValueString  proto.ColStr
	Height           proto.ColInt32
}

func NewTxs12DataBatch() *Txs12DataBatch {
	return &Txs12DataBatch{
		TxUID:            proto.ColInt64{},
		TxID:             proto.ColStr{},
		PositionInTx:     proto.ColInt32{},
		DataKey:          proto.ColStr{},
		DataType:         proto.ColStr{},
		DataValueInteger: proto.ColInt64{},
		DataValueBoolean: proto.ColBool{},
		DataValueBinary:  proto.ColStr{},
		DataValueString:  proto.ColStr{},
		Height:           proto.ColInt32{},
	}
}

func (t *Txs12DataBatch) Append(txs *Txs12Data) {
	t.TxUID.AppendArr(txs.TxUID)
	t.TxID.AppendArr(txs.TxID)
	t.PositionInTx.AppendArr(txs.PositionInTx)
	t.DataKey.AppendArr(txs.DataKey)
	t.DataType.AppendArr(txs.DataType)
	t.DataValueInteger.AppendArr(txs.DataValueInteger)
	t.DataValueBoolean.AppendArr(txs.DataValueBoolean)
	t.DataValueBinary.AppendArr(txs.DataValueBinary)
	t.DataValueString.AppendArr(txs.DataValueString)
	t.Height.AppendArr(txs.Height)
}
