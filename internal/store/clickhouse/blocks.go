package clickhouse

import (
	"time"

	"github.com/ClickHouse/ch-go/proto"
)

type UpdatesBatch struct {
	Blocks         *BlocksBatch
	Txs1           *Txs1Batch
	Txs2           *Txs2Batch
	Txs3           *Txs3Batch
	Txs4           *Txs4Batch
	Txs5           *Txs5Batch
	Txs6           *Txs6Batch
	Txs7           *Txs7Batch
	Txs8           *Txs8Batch
	Txs9           *Txs9Batch
	Txs10          *Txs10Batch
	Txs11          *Txs11Batch
	Txs11Transfers *Txs11TransferBatch
	Txs12          *Txs12Batch
	Txs12Data      *Txs12DataBatch
	Txs13          *Txs13Batch
	Txs14          *Txs14Batch
	Txs15          *Txs15Batch
	Txs16          *Txs16Batch
	Txs16Args      *TxsArgsBatch
	Txs16Payments  *TxsPaymentsBatch
	Txs17          *Txs17Batch
	Txs18          *Txs18Batch
	Txs18Args      *TxsArgsBatch
	Txs18Payments  *TxsPaymentsBatch
	Balances       *BalanceBatch
	DataStorage    *DataStorageBatch
}

func NewUpdatesBatch() *UpdatesBatch {
	return &UpdatesBatch{
		Blocks:         NewBlocksBatch(),
		Txs1:           NewTxs1Batch(),
		Txs2:           NewTxs2Batch(),
		Txs3:           NewTxs3Batch(),
		Txs4:           NewTxs4Batch(),
		Txs5:           NewTxs5Batch(),
		Txs6:           NewTxs6Batch(),
		Txs7:           NewTxs7Batch(),
		Txs8:           NewTxs8Batch(),
		Txs9:           NewTxs9Batch(),
		Txs10:          NewTxs10Batch(),
		Txs11:          NewTxs11Batch(),
		Txs11Transfers: NewTxs11TransferBatch(),
		Txs12:          NewTxs12Batch(),
		Txs12Data:      NewTxs12DataBatch(),
		Txs13:          NewTxs13Batch(),
		Txs14:          NewTxs14Batch(),
		Txs15:          NewTxs15Batch(),
		Txs16:          NewTxs16Batch(),
		Txs16Args:      NewTxsArgsBatch(),
		Txs16Payments:  NewTxsPaymentsBatch(),
		Txs17:          NewTxs17Batch(),
		Txs18:          NewTxs18Batch(),
		Txs18Args:      NewTxsArgsBatch(),
		Txs18Payments:  NewTxsPaymentsBatch(),
		Balances:       NewBalanceBatch(),
		DataStorage:    NewDataStorageBatch(),
	}
}

func (u *UpdatesBatch) Append(updates *Updates) {
	u.Blocks.Height.AppendArr(updates.Blocks.Height)
	u.Blocks.Id.AppendArr(updates.Blocks.Id)
	u.Blocks.Timestamp.AppendArr(updates.Blocks.Timestamp)

	u.Txs1.Append(updates.Txs1)
	u.Txs2.Append(updates.Txs2)
	u.Txs3.Append(updates.Txs3)
	u.Txs4.Append(updates.Txs4)
	u.Txs5.Append(updates.Txs5)
	u.Txs6.Append(updates.Txs6)
	u.Txs7.Append(updates.Txs7)
	u.Txs8.Append(updates.Txs8)
	u.Txs9.Append(updates.Txs9)
	u.Txs10.Append(updates.Txs10)
	u.Txs11.Append(updates.Txs11)
	u.Txs11Transfers.Append(updates.Txs11Transfers)
	u.Txs12.Append(updates.Txs12)
	u.Txs12Data.Append(updates.Txs12Data)
	u.Txs13.Append(updates.Txs13)
	u.Txs14.Append(updates.Txs14)
	u.Txs15.Append(updates.Txs15)
	u.Txs16.Append(updates.Txs16)
	u.Txs16Args.Append(updates.Txs16Args)
	u.Txs16Payments.Append(updates.Txs16Payments)
	u.Txs17.Append(updates.Txs17)
	u.Txs18.Append(updates.Txs18)
	u.Txs18Args.Append(updates.Txs18Args)
	u.Txs18Payments.Append(updates.Txs18Payments)
	u.Balances.Append(updates.Balances)
	u.DataStorage.Append(updates.DataStorage)
}

//func (u)

type BlocksBatch struct {
	Height    proto.ColInt32
	Id        proto.ColStr
	Timestamp proto.ColDateTime64
}

func NewBlocksBatch() *BlocksBatch {
	return &BlocksBatch{
		Height:    proto.ColInt32{},
		Id:        proto.ColStr{},
		Timestamp: proto.ColDateTime64{Precision: proto.PrecisionMilli, PrecisionSet: true},
	}
}

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

type Txs1Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs1Batch() *Txs1Batch {
	return &Txs1Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs1Batch) Append(txs *Txs1) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.Amount.AppendArr(txs.Amount)
}

type Txs2Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs2Batch() *Txs2Batch {
	return &Txs2Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs2Batch) Append(txs *Txs2) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.Amount.AppendArr(txs.Amount)
}

type Txs3Batch struct {
	*TxsBatch
	AssetId     proto.ColStr
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
		AssetId:     proto.ColStr{},
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
	t.AssetId.AppendArr(txs.AssetId)
	t.AssetName.AppendArr(txs.AssetName)
	t.Description.AppendArr(txs.Description)
	t.Amount.AppendArr(txs.Amount)
	t.Decimals.AppendArr(txs.Decimals)
	t.Reissuable.AppendArr(txs.Reissuable)
	t.Script.AppendArr(txs.Script)
}

type Txs4Batch struct {
	*TxsBatch
	AssetId          proto.ColStr
	Amount           proto.ColInt64
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	FeeAssetId       proto.ColStr
	Attachment       proto.ColStr
}

func NewTxs4Batch() *Txs4Batch {
	return &Txs4Batch{
		TxsBatch:         NewTxsBatch(),
		AssetId:          proto.ColStr{},
		Amount:           proto.ColInt64{},
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		FeeAssetId:       proto.ColStr{},
		Attachment:       proto.ColStr{},
	}
}

func (t *Txs4Batch) Append(txs *Txs4) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Amount.AppendArr(txs.Amount)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.FeeAssetId.AppendArr(txs.FeeAssetId)
	t.Attachment.AppendArr(txs.Attachment)
}

type Txs5Batch struct {
	*TxsBatch
	AssetId    proto.ColStr
	Amount     proto.ColInt64
	Reissuable proto.ColBool
}

func NewTxs5Batch() *Txs5Batch {
	return &Txs5Batch{
		TxsBatch:   NewTxsBatch(),
		AssetId:    proto.ColStr{},
		Amount:     proto.ColInt64{},
		Reissuable: proto.ColBool{},
	}
}

func (t *Txs5Batch) Append(txs *Txs5) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Amount.AppendArr(txs.Amount)
	t.Reissuable.AppendArr(txs.Reissuable)
}

type Txs6Batch struct {
	*TxsBatch
	AssetId proto.ColStr
	Amount  proto.ColInt64
}

func NewTxs6Batch() *Txs6Batch {
	return &Txs6Batch{
		TxsBatch: NewTxsBatch(),
		AssetId:  proto.ColStr{},
		Amount:   proto.ColInt64{},
	}
}

func (t *Txs6Batch) Append(txs *Txs6) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Amount.AppendArr(txs.Amount)
}

type Txs7Batch struct {
	*TxsBatch
	Amount         proto.ColInt64
	Price          proto.ColInt64
	AmountAssetId  proto.ColStr
	PriceAssetId   proto.ColStr
	BuyMatcherFee  proto.ColInt64
	SellMatcherFee proto.ColInt64
	FeeAssetId     proto.ColStr

	// Nested orders columns
	OrdersId                *proto.ColArr[string]
	OrdersVersion           *proto.ColArr[int32]
	OrdersSender            *proto.ColArr[string]
	OrdersSenderPublicKey   *proto.ColArr[string]
	OrdersMatcherPublicKey  *proto.ColArr[string]
	OrdersAmountAssetId     *proto.ColArr[string]
	OrdersPriceAssetId      *proto.ColArr[string]
	OrdersOrderType         *proto.ColArr[int32]
	OrdersAmount            *proto.ColArr[int64]
	OrdersPrice             *proto.ColArr[int64]
	OrdersTimestamp         *proto.ColArr[time.Time]
	OrdersExpiration        *proto.ColArr[time.Time]
	OrdersMatcherFee        *proto.ColArr[int64]
	OrdersMatcherFeeAssetId *proto.ColArr[string]
	OrdersProofs            proto.ColArr[[]string]
	OrdersSignature         *proto.ColArr[string]
	OrdersEip712Signature   *proto.ColArr[string]
	OrdersPriceMode         *proto.ColArr[int32]
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

type Txs8Batch struct {
	*TxsBatch
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	Amount           proto.ColInt64
}

func NewTxs8Batch() *Txs8Batch {
	return &Txs8Batch{
		TxsBatch:         NewTxsBatch(),
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		Amount:           proto.ColInt64{},
	}
}

func (t *Txs8Batch) Append(txs *Txs8) {
	t.TxsBatch.Append(txs.Txs)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.Amount.AppendArr(txs.Amount)
}

type Txs9Batch struct {
	*TxsBatch
	LeaseTxUid proto.ColInt64
}

func NewTxs9Batch() *Txs9Batch {
	return &Txs9Batch{
		TxsBatch:   NewTxsBatch(),
		LeaseTxUid: proto.ColInt64{},
	}
}

func (t *Txs9Batch) Append(txs *Txs9) {
	t.TxsBatch.Append(txs.Txs)
	t.LeaseTxUid.AppendArr(txs.LeaseTxUid)
}

type Txs10Batch struct {
	*TxsBatch
	Alias proto.ColStr
}

func NewTxs10Batch() *Txs10Batch {
	return &Txs10Batch{
		TxsBatch: NewTxsBatch(),
		Alias:    proto.ColStr{},
	}
}

func (t *Txs10Batch) Append(txs *Txs10) {
	t.TxsBatch.Append(txs.Txs)
	t.Alias.AppendArr(txs.Alias)
}

type Txs11Batch struct {
	*TxsBatch
	AssetId    proto.ColStr
	Attachment proto.ColStr
}

func NewTxs11Batch() *Txs11Batch {
	return &Txs11Batch{
		TxsBatch:   NewTxsBatch(),
		AssetId:    proto.ColStr{},
		Attachment: proto.ColStr{},
	}
}

func (t *Txs11Batch) Append(txs *Txs11) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Attachment.AppendArr(txs.Attachment)
}

type Txs11TransferBatch struct {
	TxUid            proto.ColInt64
	TxId             proto.ColStr
	PositionInTx     proto.ColInt32
	RecipientAddress proto.ColStr
	RecipientAlias   proto.ColStr
	Amount           proto.ColInt64
	Height           proto.ColInt32
}

func NewTxs11TransferBatch() *Txs11TransferBatch {
	return &Txs11TransferBatch{
		TxUid:            proto.ColInt64{},
		TxId:             proto.ColStr{},
		PositionInTx:     proto.ColInt32{},
		RecipientAddress: proto.ColStr{},
		RecipientAlias:   proto.ColStr{},
		Amount:           proto.ColInt64{},
		Height:           proto.ColInt32{},
	}
}

func (t *Txs11TransferBatch) Append(txs *Txs11Transfer) {
	t.TxUid.AppendArr(txs.TxUid)
	t.TxId.AppendArr(txs.TxId)
	t.PositionInTx.AppendArr(txs.PositionInTx)
	t.RecipientAddress.AppendArr(txs.RecipientAddress)
	t.RecipientAlias.AppendArr(txs.RecipientAlias)
	t.Amount.AppendArr(txs.Amount)
	t.Height.AppendArr(txs.Height)
}

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
	TxUid            proto.ColInt64
	TxId             proto.ColStr
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
		TxUid:            proto.ColInt64{},
		TxId:             proto.ColStr{},
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
	t.TxUid.AppendArr(txs.TxUid)
	t.TxId.AppendArr(txs.TxId)
	t.PositionInTx.AppendArr(txs.PositionInTx)
	t.DataKey.AppendArr(txs.DataKey)
	t.DataType.AppendArr(txs.DataType)
	t.DataValueInteger.AppendArr(txs.DataValueInteger)
	t.DataValueBoolean.AppendArr(txs.DataValueBoolean)
	t.DataValueBinary.AppendArr(txs.DataValueBinary)
	t.DataValueString.AppendArr(txs.DataValueString)
	t.Height.AppendArr(txs.Height)
}

type Txs13Batch struct {
	*TxsBatch
	Script proto.ColStr
}

func NewTxs13Batch() *Txs13Batch {
	return &Txs13Batch{
		TxsBatch: NewTxsBatch(),
		Script:   proto.ColStr{},
	}
}

func (t *Txs13Batch) Append(txs *Txs13) {
	t.TxsBatch.Append(txs.Txs)
	t.Script.AppendArr(txs.Script)
}

type Txs14Batch struct {
	*TxsBatch
	AssetId              proto.ColStr
	MinSponsoredAssetFee proto.ColInt64
}

func NewTxs14Batch() *Txs14Batch {
	return &Txs14Batch{
		TxsBatch:             NewTxsBatch(),
		AssetId:              proto.ColStr{},
		MinSponsoredAssetFee: proto.ColInt64{},
	}
}

func (t *Txs14Batch) Append(txs *Txs14) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.MinSponsoredAssetFee.AppendArr(txs.MinSponsoredAssetFee)
}

type Txs15Batch struct {
	*TxsBatch
	AssetId proto.ColStr
	Script  proto.ColStr
}

func NewTxs15Batch() *Txs15Batch {
	return &Txs15Batch{
		TxsBatch: NewTxsBatch(),
		AssetId:  proto.ColStr{},
		Script:   proto.ColStr{},
	}
}

func (t *Txs15Batch) Append(txs *Txs15) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.Script.AppendArr(txs.Script)
}

type Txs16Batch struct {
	*TxsBatch
	DappAddress  proto.ColStr
	DappAlias    proto.ColStr
	FunctionName proto.ColStr
	FeeAssetId   proto.ColStr
}

func NewTxs16Batch() *Txs16Batch {
	return &Txs16Batch{
		TxsBatch:     NewTxsBatch(),
		DappAddress:  proto.ColStr{},
		DappAlias:    proto.ColStr{},
		FunctionName: proto.ColStr{},
		FeeAssetId:   proto.ColStr{},
	}
}

func (t *Txs16Batch) Append(txs *Txs16) {
	t.TxsBatch.Append(txs.Txs)
	t.DappAddress.AppendArr(txs.DappAddress)
	t.DappAlias.AppendArr(txs.DappAlias)
	t.FunctionName.AppendArr(txs.FunctionName)
	t.FeeAssetId.AppendArr(txs.FeeAssetId)
}

type TxsArgsBatch struct {
	TxUid           proto.ColInt64
	TxId            proto.ColStr
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
		TxUid:           proto.ColInt64{},
		TxId:            proto.ColStr{},
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
	t.TxUid.AppendArr(txs.TxUid)
	t.TxId.AppendArr(txs.TxId)
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
	TxUid             proto.ColInt64
	TxId              proto.ColStr
	PositionInPayment proto.ColInt32
	Amount            proto.ColInt64
	Height            proto.ColInt32
	AssetId           proto.ColStr
}

func NewTxsPaymentsBatch() *TxsPaymentsBatch {
	return &TxsPaymentsBatch{
		TxUid:             proto.ColInt64{},
		TxId:              proto.ColStr{},
		PositionInPayment: proto.ColInt32{},
		Amount:            proto.ColInt64{},
		Height:            proto.ColInt32{},
		AssetId:           proto.ColStr{},
	}
}

func (t *TxsPaymentsBatch) Append(txs *TxsPayments) {
	t.TxUid.AppendArr(txs.TxUid)
	t.TxId.AppendArr(txs.TxId)
	t.PositionInPayment.AppendArr(txs.PositionInPayment)
	t.Amount.AppendArr(txs.Amount)
	t.Height.AppendArr(txs.Height)
	t.AssetId.AppendArr(txs.AssetId)
}

type Txs17Batch struct {
	*TxsBatch
	AssetId     proto.ColStr
	AssetName   proto.ColStr
	Description proto.ColStr
}

func NewTxs17Batch() *Txs17Batch {
	return &Txs17Batch{
		TxsBatch:    NewTxsBatch(),
		AssetId:     proto.ColStr{},
		AssetName:   proto.ColStr{},
		Description: proto.ColStr{},
	}
}

func (t *Txs17Batch) Append(txs *Txs17) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetId.AppendArr(txs.AssetId)
	t.AssetName.AppendArr(txs.AssetName)
	t.Description.AppendArr(txs.Description)
}

type Txs18Batch struct {
	*TxsBatch
	Bytes        proto.ColStr
	FunctionName proto.ColStr
}

func NewTxs18Batch() *Txs18Batch {
	return &Txs18Batch{
		TxsBatch:     NewTxsBatch(),
		Bytes:        proto.ColStr{},
		FunctionName: proto.ColStr{},
	}
}

func (t *Txs18Batch) Append(txs *Txs18) {
	t.TxsBatch.Append(txs.Txs)
	t.Bytes.AppendArr(txs.Bytes)
	t.FunctionName.AppendArr(txs.FunctionName)
}

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

type DataStorageBatch struct {
	Address          proto.ColStr
	DataKey          proto.ColStr
	DataType         proto.ColStr
	DataValueInteger proto.ColInt64
	DataValueBoolean proto.ColBool
	DataValueBinary  proto.ColStr
	DataValueString  proto.ColStr
	TxUid            proto.ColInt64
	TxId             proto.ColStr
	Height           proto.ColInt32
}

func NewDataStorageBatch() *DataStorageBatch {
	return &DataStorageBatch{
		Address:          proto.ColStr{},
		DataKey:          proto.ColStr{},
		DataType:         proto.ColStr{},
		DataValueInteger: proto.ColInt64{},
		DataValueBoolean: proto.ColBool{},
		DataValueBinary:  proto.ColStr{},
		DataValueString:  proto.ColStr{},
		TxUid:            proto.ColInt64{},
		TxId:             proto.ColStr{},
		Height:           proto.ColInt32{},
	}
}

func (d *DataStorageBatch) Append(dataStorage *DataStorage) {
	d.Address.AppendArr(dataStorage.Address)
	d.DataKey.AppendArr(dataStorage.DataKey)
	d.DataType.AppendArr(dataStorage.DataType)
	d.DataValueInteger.AppendArr(dataStorage.DataValueInteger)
	d.DataValueBoolean.AppendArr(dataStorage.DataValueBoolean)
	d.DataValueBinary.AppendArr(dataStorage.DataValueBinary)
	d.DataValueString.AppendArr(dataStorage.DataValueString)
	d.TxUid.AppendArr(dataStorage.TxUid)
	d.TxId.AppendArr(dataStorage.TxId)
	d.Height.AppendArr(dataStorage.Height)
}
