package clickhouse

import (
	"encoding/json"
	"time"

	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves"

	"github.com/slashfast/wavesconsumer/v2/pkg/utils"
)

const (
	UIDBase int64 = 100_000
)

type Updates struct {
	Blocks         *Blocks
	Txs1           *Txs1
	Txs2           *Txs2
	Txs3           *Txs3
	Txs4           *Txs4
	Txs5           *Txs5
	Txs6           *Txs6
	Txs7           *Txs7
	Txs8           *Txs8
	Txs9           *Txs9
	Txs10          *Txs10
	Txs11          *Txs11
	Txs11Transfers *Txs11Transfer
	Txs12          *Txs12
	Txs12Data      *Txs12Data
	Txs13          *Txs13
	Txs14          *Txs14
	Txs15          *Txs15
	Txs16          *Txs16
	Txs16Args      *TxsArgs
	Txs16Payments  *TxsPayments
	Txs17          *Txs17
	Txs18          *Txs18
	Txs18Args      *TxsArgs
	Txs18Payments  *TxsPayments
	Balances       *Balance
	DataStorage    *DataStorage
}

func NewUpdates() *Updates {
	return &Updates{
		NewBlocks(),
		NewTxs1(),
		NewTxs2(),
		NewTxs3(),
		NewTxs4(),
		NewTxs5(),
		NewTxs6(),
		NewTxs7(),
		NewTxs8(),
		NewTxs9(),
		NewTxs10(),
		NewTxs11(),
		NewTxs11Transfer(),
		NewTxs12(),
		NewTxs12Data(),
		NewTxs13(),
		NewTxs14(),
		NewTxs15(),
		NewTxs16(),
		NewTxsArgs(),
		NewTxsPayments(),
		NewTxs17(),
		NewTxs18(),
		NewTxsArgs(),
		NewTxsPayments(),
		NewBalance(),
		NewDataStorage(),
	}
}

type Blocks struct {
	Height    []int32
	ID        []string
	Timestamp []time.Time
}

func NewBlocks() *Blocks {
	return &Blocks{
		[]int32{},
		[]string{},
		[]time.Time{},
	}
}

type Txs struct {
	UID             []int64
	Type            []int8
	Sender          []string
	SenderPublicKey []string
	Timestamp       []time.Time
	Height          []int32
	ID              []string
	Signature       []string
	Proofs          [][]string
	Version         []int32
	Fee             []int64
	Status          []string
}

func NewTxs() *Txs {
	return &Txs{
		[]int64{},
		[]int8{},
		[]string{},
		[]string{},
		[]time.Time{},
		[]int32{},
		[]string{},
		[]string{},
		[][]string{},
		[]int32{},
		[]int64{},
		[]string{},
	}
}

type Txs1 struct {
	*Txs
	RecipientAddress []string
	Amount           []int64
}

func NewTxs1() *Txs1 {
	return &Txs1{
		NewTxs(),
		[]string{},
		[]int64{},
	}
}

type Txs2 struct {
	*Txs
	RecipientAddress []string
	Amount           []int64
}

func NewTxs2() *Txs2 {
	return &Txs2{
		NewTxs(),
		[]string{},
		[]int64{},
	}
}

type Txs3 struct {
	*Txs
	AssetID     []string
	AssetName   []string
	Description []string
	Amount      []int64
	Decimals    []int32
	Reissuable  []bool
	Script      []string
}

func NewTxs3() *Txs3 {
	return &Txs3{
		NewTxs(),
		[]string{},
		[]string{},
		[]string{},
		[]int64{},
		[]int32{},
		[]bool{},
		[]string{},
	}
}

type Txs4 struct {
	*Txs
	AssetID          []string
	Amount           []int64
	RecipientAddress []string
	RecipientAlias   []string
	FeeAssetID       []string
	Attachment       []string
}

func NewTxs4() *Txs4 {
	return &Txs4{
		NewTxs(),
		[]string{},
		[]int64{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
	}
}

type Txs5 struct {
	*Txs
	AssetID    []string
	Amount     []int64
	Reissuable []bool
}

func NewTxs5() *Txs5 {
	return &Txs5{
		NewTxs(),
		[]string{},
		[]int64{},
		[]bool{},
	}
}

type Txs6 struct {
	*Txs
	AssetID []string
	Amount  []int64
}

func NewTxs6() *Txs6 {
	return &Txs6{
		NewTxs(),
		[]string{},
		[]int64{},
	}
}

type Txs7 struct {
	*Txs
	Amount         []int64
	Price          []int64
	AmountAssetID  []string
	PriceAssetID   []string
	BuyMatcherFee  []int64
	SellMatcherFee []int64
	FeeAssetID     []string

	// Nested orders (always 2: buy and sell)
	OrdersID                [][]string
	OrdersVersion           [][]int32
	OrdersSender            [][]string
	OrdersSenderPublicKey   [][]string
	OrdersMatcherPublicKey  [][]string
	OrdersAmountAssetID     [][]string
	OrdersPriceAssetID      [][]string
	OrdersOrderType         [][]int32
	OrdersAmount            [][]int64
	OrdersPrice             [][]int64
	OrdersTimestamp         [][]time.Time
	OrdersExpiration        [][]time.Time
	OrdersMatcherFee        [][]int64
	OrdersMatcherFeeAssetID [][]string
	OrdersProofs            [][][]string
	OrdersSignature         [][]string
	OrdersEip712Signature   [][]string
	OrdersPriceMode         [][]int32
}

func NewTxs7() *Txs7 {
	return &Txs7{
		NewTxs(),
		[]int64{},
		[]int64{},
		[]string{},
		[]string{},
		[]int64{},
		[]int64{},
		[]string{},

		// Nested orders
		[][]string{},
		[][]int32{},
		[][]string{},
		[][]string{},
		[][]string{},
		[][]string{},
		[][]string{},
		[][]int32{},
		[][]int64{},
		[][]int64{},
		[][]time.Time{},
		[][]time.Time{},
		[][]int64{},
		[][]string{},
		[][][]string{},
		[][]string{},
		[][]string{},
		[][]int32{},
	}
}

type Txs8 struct {
	*Txs
	RecipientAddress []string
	RecipientAlias   []string
	Amount           []int64
}

func NewTxs8() *Txs8 {
	return &Txs8{
		NewTxs(),
		[]string{},
		[]string{},
		[]int64{},
	}
}

type Txs9 struct {
	*Txs
	LeaseTxUID []int64
}

func NewTxs9() *Txs9 {
	return &Txs9{
		NewTxs(),
		[]int64{},
	}
}

type Txs10 struct {
	*Txs
	Alias []string
}

func NewTxs10() *Txs10 {
	return &Txs10{
		NewTxs(),
		[]string{},
	}
}

type Txs11 struct {
	*Txs
	AssetID    []string
	Attachment []string
}

func NewTxs11() *Txs11 {
	return &Txs11{
		NewTxs(),
		[]string{},
		[]string{},
	}
}

type Txs11Transfer struct {
	TxUID            []int64
	TxID             []string
	PositionInTx     []int32
	RecipientAddress []string
	RecipientAlias   []string
	Amount           []int64
	Height           []int32
}

func NewTxs11Transfer() *Txs11Transfer {
	return &Txs11Transfer{
		[]int64{},
		[]string{},
		[]int32{},
		[]string{},
		[]string{},
		[]int64{},
		[]int32{},
	}
}

type Txs12 struct {
	*Txs
}

func NewTxs12() *Txs12 {
	return &Txs12{
		NewTxs(),
	}
}

type Txs12Data struct {
	TxUID            []int64
	TxID             []string
	PositionInTx     []int32
	DataKey          []string
	DataType         []string
	DataValueInteger []int64
	DataValueBoolean []bool
	DataValueBinary  []string
	DataValueString  []string
	Height           []int32
}

func NewTxs12Data() *Txs12Data {
	return &Txs12Data{
		[]int64{},
		[]string{},
		[]int32{},
		[]string{},
		[]string{},
		[]int64{},
		[]bool{},
		[]string{},
		[]string{},
		[]int32{},
	}
}

type Txs13 struct {
	*Txs
	Script []string
}

func NewTxs13() *Txs13 {
	return &Txs13{
		NewTxs(),
		[]string{},
	}
}

type Txs14 struct {
	*Txs
	AssetID              []string
	MinSponsoredAssetFee []int64
}

func NewTxs14() *Txs14 {
	return &Txs14{
		NewTxs(),
		[]string{},
		[]int64{},
	}
}

type Txs15 struct {
	*Txs
	AssetID []string
	Script  []string
}

func NewTxs15() *Txs15 {
	return &Txs15{
		NewTxs(),
		[]string{},
		[]string{},
	}
}

type Txs16 struct {
	*Txs
	DappAddress  []string
	DappAlias    []string
	FunctionName []string
	FeeAssetID   []string
}

func NewTxs16() *Txs16 {
	return &Txs16{
		NewTxs(),
		[]string{},
		[]string{},
		[]string{},
		[]string{},
	}
}

type TxsArgs struct {
	TxUID           []int64
	TxID            []string
	PositionInArgs  []int32
	ArgType         []string
	ArgValueInteger []int64
	ArgValueBoolean []bool
	ArgValueBinary  []string
	ArgValueString  []string
	ArgValueList    []string
	Height          []int32
}

func NewTxsArgs() *TxsArgs {
	return &TxsArgs{
		[]int64{},
		[]string{},
		[]int32{},
		[]string{},
		[]int64{},
		[]bool{},
		[]string{},
		[]string{},
		[]string{},
		[]int32{},
	}
}

func AppendTxArg(argsPtr *TxsArgs, arg *waves.InvokeScriptResult_Call_Argument, txID string, txUID int64, pos, height int32) error {
	var argType string
	switch arg.Value.(type) {
	case *waves.InvokeScriptResult_Call_Argument_IntegerValue:
		argType = "integer"
	case *waves.InvokeScriptResult_Call_Argument_BooleanValue:
		argType = "boolean"
	case *waves.InvokeScriptResult_Call_Argument_BinaryValue:
		argType = "binary"
	case *waves.InvokeScriptResult_Call_Argument_StringValue:
		argType = "string"
	case *waves.InvokeScriptResult_Call_Argument_List_:
		argType = "list"
	default:
		argType = ""
	}

	argValueList, err := json.Marshal(arg.GetList())
	if err != nil {
		return err
	}

	argsPtr.TxUID = append(argsPtr.TxUID, txUID)
	argsPtr.TxID = append(argsPtr.TxID, txID)
	argsPtr.PositionInArgs = append(argsPtr.PositionInArgs, pos)
	argsPtr.ArgType = append(argsPtr.ArgType, argType)
	argsPtr.ArgValueInteger = append(argsPtr.ArgValueInteger, arg.GetIntegerValue())
	argsPtr.ArgValueBoolean = append(argsPtr.ArgValueBoolean, arg.GetBooleanValue())
	argsPtr.ArgValueBinary = append(argsPtr.ArgValueBinary, utils.ToPrefixedBase64(arg.GetBinaryValue()))
	argsPtr.ArgValueString = append(argsPtr.ArgValueString, utils.EscapeUnicodeNull(arg.GetStringValue()))
	argsPtr.ArgValueList = append(argsPtr.ArgValueList, string(argValueList))
	argsPtr.Height = append(argsPtr.Height, height)

	return nil
}

type TxsPayments struct {
	TxUID             []int64
	TxID              []string
	PositionInPayment []int32
	Amount            []int64
	Height            []int32
	AssetID           []string
}

func NewTxsPayments() *TxsPayments {
	return &TxsPayments{
		[]int64{},
		[]string{},
		[]int32{},
		[]int64{},
		[]int32{},
		[]string{},
	}
}

func AppendTxPayment(paymentsPtr *TxsPayments, payment *waves.Amount, txID string, txUID int64, pos, height int32) {
	paymentsPtr.TxUID = append(paymentsPtr.TxUID, txUID)
	paymentsPtr.TxID = append(paymentsPtr.TxID, txID)
	paymentsPtr.PositionInPayment = append(paymentsPtr.PositionInPayment, pos)
	paymentsPtr.Amount = append(paymentsPtr.Amount, payment.Amount)
	paymentsPtr.Height = append(paymentsPtr.Height, height)
	paymentsPtr.AssetID = append(paymentsPtr.AssetID, utils.ExtractAssetID(payment.AssetId))
}

type Txs17 struct {
	*Txs
	AssetID     []string
	AssetName   []string
	Description []string
}

func NewTxs17() *Txs17 {
	return &Txs17{
		NewTxs(),
		[]string{},
		[]string{},
		[]string{},
	}
}

type Txs18 struct {
	*Txs
	Bytes        []string
	FunctionName []string
}

func NewTxs18() *Txs18 {
	return &Txs18{
		NewTxs(),
		[]string{},
		[]string{},
	}
}

type Balance struct {
	Address []string
	AssetID []string
	Amount  []int64
	Delta   []int64
	TxUID   []int64
	TxID    []string
	Height  []int32
}

func NewBalance() *Balance {
	return &Balance{
		[]string{},
		[]string{},
		[]int64{},
		[]int64{},
		[]int64{},
		[]string{},
		[]int32{},
	}
}

type DataStorage struct {
	Address          []string
	DataKey          []string
	DataType         []string
	DataValueInteger []int64
	DataValueBoolean []bool
	DataValueBinary  []string
	DataValueString  []string
	TxUID            []int64
	TxID             []string
	Height           []int32
}

func NewDataStorage() *DataStorage {
	return &DataStorage{
		[]string{},
		[]string{},
		[]string{},
		[]int64{},
		[]bool{},
		[]string{},
		[]string{},
		[]int64{},
		[]string{},
		[]int32{},
	}
}
