package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs16Batch struct {
	*TxsBatch
	DappAddress  proto.ColStr
	DappAlias    proto.ColStr
	FunctionName proto.ColStr
	FeeAssetID   proto.ColStr
}

func NewTxs16Batch() *Txs16Batch {
	return &Txs16Batch{
		TxsBatch:     NewTxsBatch(),
		DappAddress:  proto.ColStr{},
		DappAlias:    proto.ColStr{},
		FunctionName: proto.ColStr{},
		FeeAssetID:   proto.ColStr{},
	}
}

func (t *Txs16Batch) Append(txs *Txs16) {
	t.TxsBatch.Append(txs.Txs)
	t.DappAddress.AppendArr(txs.DappAddress)
	t.DappAlias.AppendArr(txs.DappAlias)
	t.FunctionName.AppendArr(txs.FunctionName)
	t.FeeAssetID.AppendArr(txs.FeeAssetID)
}

type TxsArgsBatch struct {
	TxUID           proto.ColInt64
	TxID            proto.ColStr
	PositionInArgs  proto.ColInt32
	ArgType         proto.ColStr
	ArgValueInteger proto.ColInt64
	ArgValueBoolean proto.ColBool
	ArgValueBinary  proto.ColStr
	ArgValueString  proto.ColStr
	ArgValueList    proto.ColStr
	Height          proto.ColInt32
}

func NewTxsArgsBatch() *TxsArgsBatch {
	return &TxsArgsBatch{
		TxUID:           proto.ColInt64{},
		TxID:            proto.ColStr{},
		PositionInArgs:  proto.ColInt32{},
		ArgType:         proto.ColStr{},
		ArgValueInteger: proto.ColInt64{},
		ArgValueBoolean: proto.ColBool{},
		ArgValueBinary:  proto.ColStr{},
		ArgValueString:  proto.ColStr{},
		ArgValueList:    proto.ColStr{},
		Height:          proto.ColInt32{},
	}
}

func (t *TxsArgsBatch) Append(txs *TxsArgs) {
	t.TxUID.AppendArr(txs.TxUID)
	t.TxID.AppendArr(txs.TxID)
	t.PositionInArgs.AppendArr(txs.PositionInArgs)
	t.ArgType.AppendArr(txs.ArgType)
	t.ArgValueInteger.AppendArr(txs.ArgValueInteger)
	t.ArgValueBoolean.AppendArr(txs.ArgValueBoolean)
	t.ArgValueBinary.AppendArr(txs.ArgValueBinary)
	t.ArgValueString.AppendArr(txs.ArgValueString)
	t.ArgValueList.AppendArr(txs.ArgValueList)
	t.Height.AppendArr(txs.Height)
}

type TxsPaymentsBatch struct {
	TxUID             proto.ColInt64
	TxID              proto.ColStr
	PositionInPayment proto.ColInt32
	Amount            proto.ColInt64
	Height            proto.ColInt32
	AssetID           proto.ColStr
}

func NewTxsPaymentsBatch() *TxsPaymentsBatch {
	return &TxsPaymentsBatch{
		TxUID:             proto.ColInt64{},
		TxID:              proto.ColStr{},
		PositionInPayment: proto.ColInt32{},
		Amount:            proto.ColInt64{},
		Height:            proto.ColInt32{},
		AssetID:           proto.ColStr{},
	}
}

func (t *TxsPaymentsBatch) Append(txs *TxsPayments) {
	t.TxUID.AppendArr(txs.TxUID)
	t.TxID.AppendArr(txs.TxID)
	t.PositionInPayment.AppendArr(txs.PositionInPayment)
	t.Amount.AppendArr(txs.Amount)
	t.Height.AppendArr(txs.Height)
	t.AssetID.AppendArr(txs.AssetID)
}
