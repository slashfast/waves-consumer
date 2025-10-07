package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs3Batch struct {
	*TxsBatch
	AssetID     proto.ColStr
	AssetName   proto.ColStr
	Description proto.ColStr
	Amount      proto.ColInt64
	Decimals    proto.ColInt32
	Reissuable  proto.ColBool
	Script      proto.ColStr
}

func NewTxs3Batch() *Txs3Batch {
	return &Txs3Batch{
		TxsBatch:    NewTxsBatch(),
		AssetID:     proto.ColStr{},
		AssetName:   proto.ColStr{},
		Description: proto.ColStr{},
		Amount:      proto.ColInt64{},
		Decimals:    proto.ColInt32{},
		Reissuable:  proto.ColBool{},
		Script:      proto.ColStr{},
	}
}

func (t *Txs3Batch) Append(txs *Txs3) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.AssetName.AppendArr(txs.AssetName)
	t.Description.AppendArr(txs.Description)
	t.Amount.AppendArr(txs.Amount)
	t.Decimals.AppendArr(txs.Decimals)
	t.Reissuable.AppendArr(txs.Reissuable)
	t.Script.AppendArr(txs.Script)
}
