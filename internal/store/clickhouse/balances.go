package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type BalanceBatch struct {
	Address proto.ColStr
	AssetID proto.ColStr
	Amount  proto.ColInt64
	Delta   proto.ColInt64
	TxUID   proto.ColInt64
	TxID    proto.ColStr
	Height  proto.ColInt32
}

func NewBalanceBatch() *BalanceBatch {
	return &BalanceBatch{
		Address: proto.ColStr{},
		AssetID: proto.ColStr{},
		Amount:  proto.ColInt64{},
		Delta:   proto.ColInt64{},
		TxUID:   proto.ColInt64{},
		TxID:    proto.ColStr{},
		Height:  proto.ColInt32{},
	}
}

func (b *BalanceBatch) Append(balance *Balance) {
	b.Address.AppendArr(balance.Address)
	b.AssetID.AppendArr(balance.AssetID)
	b.Amount.AppendArr(balance.Amount)
	b.Delta.AppendArr(balance.Delta)
	b.TxUID.AppendArr(balance.TxUID)
	b.TxID.AppendArr(balance.TxID)
	b.Height.AppendArr(balance.Height)
}
