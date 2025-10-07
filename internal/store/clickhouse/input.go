package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type UpdatesInput map[string]proto.Input

func NewBaseTxInput(txs *TxsBatch) proto.Input {
	return proto.Input{
		{Name: "uid", Data: &txs.UID},
		{Name: "type", Data: &txs.Type},
		{Name: "sender", Data: &txs.Sender},
		{Name: "sender_public_key", Data: &txs.SenderPublicKey},
		{Name: "timestamp", Data: &txs.Timestamp},
		{Name: "height", Data: &txs.Height},
		{Name: "id", Data: &txs.ID},
		{Name: "signature", Data: &txs.Signature},
		{Name: "proofs", Data: txs.Proofs},
		{Name: "version", Data: &txs.Version},
		{Name: "fee", Data: &txs.Fee},
		{Name: "status", Data: &txs.Status},
	}
}

func NewTxsArgsInput(args *TxsArgsBatch) proto.Input {
	return proto.Input{
		{Name: "tx_uid", Data: &args.TxUID},
		{Name: "tx_id", Data: &args.TxID},
		{Name: "position_in_args", Data: args.PositionInArgs},
		{Name: "arg_type", Data: &args.ArgType},
		{Name: "arg_value_integer", Data: &args.ArgValueInteger},
		{Name: "arg_value_bool", Data: &args.ArgValueBoolean},
		{Name: "arg_value_binary", Data: &args.ArgValueBinary},
		{Name: "arg_value_string", Data: &args.ArgValueString},
		{Name: "arg_value_list", Data: &args.ArgValueList},
		{Name: "height", Data: &args.Height},
	}
}

func NewTxsPaymentsInput(payments *TxsPaymentsBatch) proto.Input {
	return proto.Input{
		{Name: "tx_uid", Data: &payments.TxUID},
		{Name: "tx_id", Data: &payments.TxID},
		{Name: "position_in_payment", Data: payments.PositionInPayment},
		{Name: "amount", Data: &payments.Amount},
		{Name: "asset_id", Data: &payments.AssetID},
		{Name: "height", Data: &payments.Height},
	}
}

func NewUpdatesInput(updates *UpdatesBatch) UpdatesInput {
	return UpdatesInput{
		"blocks": {
			{Name: "height", Data: &updates.Blocks.Height},
			{Name: "id", Data: &updates.Blocks.ID},
			{Name: "timestamp", Data: &updates.Blocks.Timestamp},
		},
		"txs_1": append(
			NewBaseTxInput(updates.Txs1.TxsBatch),
			proto.InputColumn{Name: "recipient_address", Data: &updates.Txs1.RecipientAddress},
			proto.InputColumn{Name: "amount", Data: &updates.Txs1.Amount},
		),
		"txs_2": append(
			NewBaseTxInput(updates.Txs2.TxsBatch),
			proto.InputColumn{Name: "recipient_address", Data: &updates.Txs2.RecipientAddress},
			proto.InputColumn{Name: "amount", Data: &updates.Txs2.Amount},
		),
		"txs_3": append(
			NewBaseTxInput(updates.Txs3.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs3.AssetID},
			proto.InputColumn{Name: "asset_name", Data: &updates.Txs3.AssetName},
			proto.InputColumn{Name: "description", Data: &updates.Txs3.Description},
			proto.InputColumn{Name: "amount", Data: &updates.Txs3.Amount},
			proto.InputColumn{Name: "decimals", Data: &updates.Txs3.Decimals},
			proto.InputColumn{Name: "reissuable", Data: &updates.Txs3.Reissuable},
			proto.InputColumn{Name: "script", Data: &updates.Txs3.Script},
		),
		"txs_4": append(
			NewBaseTxInput(updates.Txs4.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs4.AssetID},
			proto.InputColumn{Name: "amount", Data: &updates.Txs4.Amount},
			proto.InputColumn{Name: "recipient_address", Data: &updates.Txs4.RecipientAddress},
			proto.InputColumn{Name: "recipient_alias", Data: &updates.Txs4.RecipientAlias},
			proto.InputColumn{Name: "fee_asset_id", Data: &updates.Txs4.FeeAssetID},
			proto.InputColumn{Name: "attachment", Data: &updates.Txs4.Attachment},
		),
		"txs_5": append(
			NewBaseTxInput(updates.Txs5.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs5.AssetID},
			proto.InputColumn{Name: "amount", Data: &updates.Txs5.Amount},
			proto.InputColumn{Name: "reissuable", Data: &updates.Txs5.Reissuable},
		),
		"txs_6": append(
			NewBaseTxInput(updates.Txs6.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs6.AssetID},
			proto.InputColumn{Name: "amount", Data: &updates.Txs6.Amount},
		),
		"txs_7": append(
			NewBaseTxInput(updates.Txs7.TxsBatch),
			proto.InputColumn{Name: "amount", Data: &updates.Txs7.Amount},
			proto.InputColumn{Name: "price", Data: &updates.Txs7.Price},
			proto.InputColumn{Name: "amount_asset_id", Data: &updates.Txs7.AmountAssetID},
			proto.InputColumn{Name: "price_asset_id", Data: &updates.Txs7.PriceAssetID},
			proto.InputColumn{Name: "buy_matcher_fee", Data: &updates.Txs7.BuyMatcherFee},
			proto.InputColumn{Name: "sell_matcher_fee", Data: &updates.Txs7.SellMatcherFee},
			proto.InputColumn{Name: "fee_asset_id", Data: &updates.Txs7.FeeAssetID},
			// Nested orders
			proto.InputColumn{Name: "orders.id", Data: updates.Txs7.OrdersID},
			proto.InputColumn{Name: "orders.version", Data: updates.Txs7.OrdersVersion},
			proto.InputColumn{Name: "orders.sender", Data: updates.Txs7.OrdersSender},
			proto.InputColumn{Name: "orders.sender_public_key", Data: updates.Txs7.OrdersSenderPublicKey},
			proto.InputColumn{Name: "orders.matcher_public_key", Data: updates.Txs7.OrdersMatcherPublicKey},
			proto.InputColumn{Name: "orders.amount_asset_id", Data: updates.Txs7.OrdersAmountAssetID},
			proto.InputColumn{Name: "orders.price_asset_id", Data: updates.Txs7.OrdersPriceAssetID},
			proto.InputColumn{Name: "orders.order_type", Data: updates.Txs7.OrdersOrderType},
			proto.InputColumn{Name: "orders.amount", Data: updates.Txs7.OrdersAmount},
			proto.InputColumn{Name: "orders.price", Data: updates.Txs7.OrdersPrice},
			proto.InputColumn{Name: "orders.timestamp", Data: updates.Txs7.OrdersTimestamp},
			proto.InputColumn{Name: "orders.expiration", Data: updates.Txs7.OrdersExpiration},
			proto.InputColumn{Name: "orders.matcher_fee", Data: updates.Txs7.OrdersMatcherFee},
			proto.InputColumn{Name: "orders.matcher_fee_asset_id", Data: updates.Txs7.OrdersMatcherFeeAssetID},
			proto.InputColumn{Name: "orders.proofs", Data: &updates.Txs7.OrdersProofs},
			proto.InputColumn{Name: "orders.signature", Data: updates.Txs7.OrdersSignature},
			proto.InputColumn{Name: "orders.eip712_signature", Data: updates.Txs7.OrdersEip712Signature},
			proto.InputColumn{Name: "orders.price_mode", Data: updates.Txs7.OrdersPriceMode},
		),
		"txs_8": append(
			NewBaseTxInput(updates.Txs8.TxsBatch),
			proto.InputColumn{Name: "recipient_address", Data: &updates.Txs8.RecipientAddress},
			proto.InputColumn{Name: "recipient_alias", Data: &updates.Txs8.RecipientAlias},
			proto.InputColumn{Name: "amount", Data: &updates.Txs8.Amount},
		),
		"txs_9": append(
			NewBaseTxInput(updates.Txs9.TxsBatch),
			proto.InputColumn{Name: "lease_tx_uid", Data: &updates.Txs9.LeaseTxUID},
		),
		"txs_10": append(
			NewBaseTxInput(updates.Txs10.TxsBatch),
			proto.InputColumn{Name: "alias", Data: &updates.Txs10.Alias},
		),
		"txs_11": append(
			NewBaseTxInput(updates.Txs11.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs11.AssetID},
			proto.InputColumn{Name: "attachment", Data: &updates.Txs11.Attachment},
		),
		"txs_11_transfers": {
			{Name: "tx_uid", Data: &updates.Txs11Transfers.TxUID},
			{Name: "tx_id", Data: &updates.Txs11Transfers.TxID},
			{Name: "position_in_tx", Data: &updates.Txs11Transfers.PositionInTx},
			{Name: "recipient_address", Data: &updates.Txs11Transfers.RecipientAddress},
			{Name: "recipient_alias", Data: &updates.Txs11Transfers.RecipientAlias},
			{Name: "amount", Data: &updates.Txs11Transfers.Amount},
			{Name: "height", Data: &updates.Txs11Transfers.Height},
		},
		"txs_12": NewBaseTxInput(updates.Txs12.TxsBatch),
		"txs_12_data": {
			{Name: "tx_uid", Data: &updates.Txs12Data.TxUID},
			{Name: "tx_id", Data: &updates.Txs12Data.TxID},
			{Name: "position_in_tx", Data: &updates.Txs12Data.PositionInTx},
			{Name: "data_key", Data: &updates.Txs12Data.DataKey},
			{Name: "data_type", Data: &updates.Txs12Data.DataType},
			{Name: "data_value_integer", Data: &updates.Txs12Data.DataValueInteger},
			{Name: "data_value_bool", Data: &updates.Txs12Data.DataValueBoolean},
			{Name: "data_value_binary", Data: &updates.Txs12Data.DataValueBinary},
			{Name: "data_value_string", Data: &updates.Txs12Data.DataValueString},
			{Name: "height", Data: &updates.Txs12Data.Height},
		},
		"txs_13": append(
			NewBaseTxInput(updates.Txs13.TxsBatch),
			proto.InputColumn{Name: "script", Data: &updates.Txs13.Script},
		),
		"txs_14": append(
			NewBaseTxInput(updates.Txs14.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs14.AssetID},
			proto.InputColumn{Name: "min_sponsored_asset_fee", Data: &updates.Txs14.MinSponsoredAssetFee},
		),
		"txs_15": append(
			NewBaseTxInput(updates.Txs15.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs15.AssetID},
			proto.InputColumn{Name: "script", Data: &updates.Txs15.Script},
		),
		"txs_16": append(
			NewBaseTxInput(updates.Txs16.TxsBatch),
			proto.InputColumn{Name: "dapp_address", Data: &updates.Txs16.DappAddress},
			proto.InputColumn{Name: "dapp_alias", Data: &updates.Txs16.DappAlias},
			proto.InputColumn{Name: "function_name", Data: &updates.Txs16.FunctionName},
			proto.InputColumn{Name: "fee_asset_id", Data: &updates.Txs16.FeeAssetID},
		),
		"txs_16_args":     NewTxsArgsInput(updates.Txs16Args),
		"txs_16_payments": NewTxsPaymentsInput(updates.Txs16Payments),
		"txs_17": append(
			NewBaseTxInput(updates.Txs17.TxsBatch),
			proto.InputColumn{Name: "asset_id", Data: &updates.Txs17.AssetID},
			proto.InputColumn{Name: "asset_name", Data: &updates.Txs17.AssetName},
			proto.InputColumn{Name: "description", Data: &updates.Txs17.Description},
		),
		"txs_18": append(
			NewBaseTxInput(updates.Txs18.TxsBatch),
			proto.InputColumn{Name: "bytes", Data: &updates.Txs18.Bytes},
			proto.InputColumn{Name: "function_name", Data: &updates.Txs18.FunctionName},
		),
		"txs_18_args":     NewTxsArgsInput(updates.Txs18Args),
		"txs_18_payments": NewTxsPaymentsInput(updates.Txs18Payments),
		"balances": {
			{Name: "address", Data: &updates.Balances.Address},
			{Name: "asset_id", Data: &updates.Balances.AssetID},
			{Name: "amount", Data: &updates.Balances.Amount},
			{Name: "delta", Data: &updates.Balances.Delta},
			{Name: "tx_uid", Data: &updates.Balances.TxUID},
			{Name: "tx_id", Data: &updates.Balances.TxID},
			{Name: "height", Data: &updates.Balances.Height},
		},
		"data_storage": {
			{Name: "address", Data: &updates.DataStorage.Address},
			{Name: "data_key", Data: &updates.DataStorage.DataKey},
			{Name: "data_type", Data: &updates.DataStorage.DataType},
			{Name: "data_value_integer", Data: &updates.DataStorage.DataValueInteger},
			{Name: "data_value_bool", Data: &updates.DataStorage.DataValueBoolean},
			{Name: "data_value_binary", Data: &updates.DataStorage.DataValueBinary},
			{Name: "data_value_string", Data: &updates.DataStorage.DataValueString},
			{Name: "tx_uid", Data: &updates.DataStorage.TxUID},
			{Name: "tx_id", Data: &updates.DataStorage.TxID},
			{Name: "height", Data: &updates.DataStorage.Height},
		},
	}
}
