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
	OrdersPriceAssetId      *proto.ColArr[string]
	OrdersMatcherFeeAssetId *proto.ColArr[string]
	OrdersId                *proto.ColArr[string]
	OrdersVersion           *proto.ColArr[int32]
	OrdersSender            *proto.ColArr[string]
	OrdersSenderPublicKey   *proto.ColArr[string]
	OrdersMatcherPublicKey  *proto.ColArr[string]
	OrdersOrderType         *proto.ColArr[int32]
	OrdersExpiration        *proto.ColArr[time.Time]
	OrdersAmountAssetId     *proto.ColArr[string]
	OrdersAmount            *proto.ColArr[int64]
	OrdersPrice             *proto.ColArr[int64]
	*TxsBatch
	FeeAssetId     proto.ColStr
	PriceAssetId   proto.ColStr
	AmountAssetId  proto.ColStr
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
		AmountAssetId:  proto.ColStr{},
		PriceAssetId:   proto.ColStr{},
		BuyMatcherFee:  proto.ColInt64{},
		SellMatcherFee: proto.ColInt64{},
		FeeAssetId:     proto.ColStr{},

		// Nested orders
		OrdersId:                new(proto.ColStr).Array(),
		OrdersVersion:           new(proto.ColInt32).Array(),
		OrdersSender:            new(proto.ColStr).Array(),
		OrdersSenderPublicKey:   new(proto.ColStr).Array(),
		OrdersMatcherPublicKey:  new(proto.ColStr).Array(),
		OrdersAmountAssetId:     new(proto.ColStr).Array(),
		OrdersPriceAssetId:      new(proto.ColStr).Array(),
		OrdersOrderType:         new(proto.ColInt32).Array(),
		OrdersAmount:            new(proto.ColInt64).Array(),
		OrdersPrice:             new(proto.ColInt64).Array(),
		OrdersTimestamp:         &proto.ColArr[time.Time]{Data: &proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}},
		OrdersExpiration:        &proto.ColArr[time.Time]{Data: &proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true}},
		OrdersMatcherFee:        new(proto.ColInt64).Array(),
		OrdersMatcherFeeAssetId: new(proto.ColStr).Array(),
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
	t.AmountAssetId.AppendArr(txs.AmountAssetId)
	t.PriceAssetId.AppendArr(txs.PriceAssetId)
	t.BuyMatcherFee.AppendArr(txs.BuyMatcherFee)
	t.SellMatcherFee.AppendArr(txs.SellMatcherFee)
	t.FeeAssetId.AppendArr(txs.FeeAssetId)

	// Nested orders
	t.OrdersId.AppendArr(txs.OrdersId)
	t.OrdersVersion.AppendArr(txs.OrdersVersion)
	t.OrdersSender.AppendArr(txs.OrdersSender)
	t.OrdersSenderPublicKey.AppendArr(txs.OrdersSenderPublicKey)
	t.OrdersMatcherPublicKey.AppendArr(txs.OrdersMatcherPublicKey)
	t.OrdersAmountAssetId.AppendArr(txs.OrdersAmountAssetId)
	t.OrdersPriceAssetId.AppendArr(txs.OrdersPriceAssetId)
	t.OrdersOrderType.AppendArr(txs.OrdersOrderType)
	t.OrdersAmount.AppendArr(txs.OrdersAmount)
	t.OrdersPrice.AppendArr(txs.OrdersPrice)
	t.OrdersTimestamp.AppendArr(txs.OrdersTimestamp)
	t.OrdersExpiration.AppendArr(txs.OrdersExpiration)
	t.OrdersMatcherFee.AppendArr(txs.OrdersMatcherFee)
	t.OrdersMatcherFeeAssetId.AppendArr(txs.OrdersMatcherFeeAssetId)
	t.OrdersProofs.AppendArr(txs.OrdersProofs)
	t.OrdersSignature.AppendArr(txs.OrdersSignature)
	t.OrdersEip712Signature.AppendArr(txs.OrdersEip712Signature)
	t.OrdersPriceMode.AppendArr(txs.OrdersPriceMode)
}
