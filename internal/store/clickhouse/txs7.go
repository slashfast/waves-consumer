package clickhouse

import (
	"time"

	"github.com/ClickHouse/ch-go/proto"
)

type Txs7Batch struct {
	OrdersProofs            proto.ColArr[[]string]
	OrdersTimestamp         *proto.ColArr[time.Time]
	OrdersMatcherFee        *proto.ColArr[int64]
	OrdersPriceMode         *proto.ColArr[int32]
	OrdersEip712Signature   *proto.ColArr[string]
	OrdersSignature         *proto.ColArr[string]
	OrdersPriceAssetID      *proto.ColArr[string]
	OrdersMatcherFeeAssetID *proto.ColArr[string]
	OrdersID                *proto.ColArr[string]
	OrdersVersion           *proto.ColArr[int32]
	OrdersSender            *proto.ColArr[string]
	OrdersSenderPublicKey   *proto.ColArr[string]
	OrdersMatcherPublicKey  *proto.ColArr[string]
	OrdersOrderType         *proto.ColArr[int32]
	OrdersExpiration        *proto.ColArr[time.Time]
	OrdersAmountAssetID     *proto.ColArr[string]
	OrdersAmount            *proto.ColArr[int64]
	OrdersPrice             *proto.ColArr[int64]
	*TxsBatch
	FeeAssetID     proto.ColStr
	PriceAssetID   proto.ColStr
	AmountAssetID  proto.ColStr
	Amount         proto.ColInt64
	Price          proto.ColInt64
	SellMatcherFee proto.ColInt64
	BuyMatcherFee  proto.ColInt64
}

func NewTxs7Batch() *Txs7Batch {
	return &Txs7Batch{
		TxsBatch:       NewTxsBatch(),
		Amount:         proto.ColInt64{},
		Price:          proto.ColInt64{},
		AmountAssetID:  proto.ColStr{},
		PriceAssetID:   proto.ColStr{},
		BuyMatcherFee:  proto.ColInt64{},
		SellMatcherFee: proto.ColInt64{},
		FeeAssetID:     proto.ColStr{},

		// Nested orders
		OrdersID:                new(proto.ColStr).Array(),
		OrdersVersion:           new(proto.ColInt32).Array(),
		OrdersSender:            new(proto.ColStr).Array(),
		OrdersSenderPublicKey:   new(proto.ColStr).Array(),
		OrdersMatcherPublicKey:  new(proto.ColStr).Array(),
		OrdersAmountAssetID:     new(proto.ColStr).Array(),
		OrdersPriceAssetID:      new(proto.ColStr).Array(),
		OrdersOrderType:         new(proto.ColInt32).Array(),
		OrdersAmount:            new(proto.ColInt64).Array(),
		OrdersPrice:             new(proto.ColInt64).Array(),
		OrdersTimestamp:         &proto.ColArr[time.Time]{Data: &proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}},
		OrdersExpiration:        &proto.ColArr[time.Time]{Data: &proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}},
		OrdersMatcherFee:        new(proto.ColInt64).Array(),
		OrdersMatcherFeeAssetID: new(proto.ColStr).Array(),
		OrdersProofs:            proto.ColArr[[]string]{Data: new(proto.ColStr).Array()},
		OrdersSignature:         new(proto.ColStr).Array(),
		OrdersEip712Signature:   new(proto.ColStr).Array(),
		OrdersPriceMode:         new(proto.ColInt32).Array(),
	}
}

func (t *Txs7Batch) Append(txs *Txs7) {
	t.TxsBatch.Append(txs.Txs)
	t.Amount.AppendArr(txs.Amount)
	t.Price.AppendArr(txs.Price)
	t.AmountAssetID.AppendArr(txs.AmountAssetID)
	t.PriceAssetID.AppendArr(txs.PriceAssetID)
	t.BuyMatcherFee.AppendArr(txs.BuyMatcherFee)
	t.SellMatcherFee.AppendArr(txs.SellMatcherFee)
	t.FeeAssetID.AppendArr(txs.FeeAssetID)

	// Nested orders
	t.OrdersID.AppendArr(txs.OrdersID)
	t.OrdersVersion.AppendArr(txs.OrdersVersion)
	t.OrdersSender.AppendArr(txs.OrdersSender)
	t.OrdersSenderPublicKey.AppendArr(txs.OrdersSenderPublicKey)
	t.OrdersMatcherPublicKey.AppendArr(txs.OrdersMatcherPublicKey)
	t.OrdersAmountAssetID.AppendArr(txs.OrdersAmountAssetID)
	t.OrdersPriceAssetID.AppendArr(txs.OrdersPriceAssetID)
	t.OrdersOrderType.AppendArr(txs.OrdersOrderType)
	t.OrdersAmount.AppendArr(txs.OrdersAmount)
	t.OrdersPrice.AppendArr(txs.OrdersPrice)
	t.OrdersTimestamp.AppendArr(txs.OrdersTimestamp)
	t.OrdersExpiration.AppendArr(txs.OrdersExpiration)
	t.OrdersMatcherFee.AppendArr(txs.OrdersMatcherFee)
	t.OrdersMatcherFeeAssetID.AppendArr(txs.OrdersMatcherFeeAssetID)
	t.OrdersProofs.AppendArr(txs.OrdersProofs)
	t.OrdersSignature.AppendArr(txs.OrdersSignature)
	t.OrdersEip712Signature.AppendArr(txs.OrdersEip712Signature)
	t.OrdersPriceMode.AppendArr(txs.OrdersPriceMode)
}
