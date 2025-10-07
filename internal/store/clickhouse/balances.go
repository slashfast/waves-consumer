package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type BalanceBatch struct {
	Address proto.ColStr
	AssetId proto.ColStr
	Amount  proto.ColInt64
	Delta   proto.ColInt64
	TxUid   proto.ColInt64
	TxId    proto.ColStr
	Height  proto.ColInt32
}

func NewBalanceBatch() *BalanceBatch {
	return &BalanceBatch{
		Address: proto.ColStr{},
		AssetId: proto.ColStr{},
		Amount:  proto.ColInt64{},
		Delta:   proto.ColInt64{},
		TxUid:   proto.ColInt64{},
		TxId:    proto.ColStr{},
		Height:  proto.ColInt32{},
	}
}

func (b *BalanceBatch) Append(balance *Balance) {
	b.Address.AppendArr(balance.Address)
	b.AssetId.AppendArr(balance.AssetId)
	b.Amount.AppendArr(balance.Amount)
	b.Delta.AppendArr(balance.Delta)
	b.TxUid.AppendArr(balance.TxUid)
	b.TxId.AppendArr(balance.TxId)
	b.Height.AppendArr(balance.Height)
}
