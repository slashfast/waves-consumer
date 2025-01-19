package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type UpdatesInput map[string]proto.Input

func NewBaseTxInput(txs *Txs) proto.Input {
	return proto.Input{
		{"uid", &txs.Uid},
		{"type", &txs.Type},
		{"sender", &txs.Sender},
		{"sender_public_key", &txs.SenderPublicKey},
		{"timestamp", &txs.Timestamp},
		{"height", &txs.Height},
		{"id", &txs.Id},
		{"signature", &txs.Signature},
		{"proofs", txs.Proofs},
		{"version", &txs.Version},
		{"fee", &txs.Fee},
		{"status", &txs.Status},
	}
}

func NewTxsArgsInput(args *TxsArgs) proto.Input {
	return proto.Input{
		{"tx_uid", &args.TxUid},
		{"tx_id", &args.TxId},
		{"position_in_args", args.PositionInArgs},
		{"arg_type", &args.ArgType},
		{"arg_value_integer", &args.ArgValueInteger},
		{"arg_value_bool", &args.ArgValueBoolean},
		{"arg_value_binary", &args.ArgValueBinary},
		{"arg_value_string", &args.ArgValueString},
		{"arg_value_list", &args.ArgValueList},
		{"height", &args.Height},
	}
}

func NewTxsPaymentsInput(payments *TxsPayments) proto.Input {
	return proto.Input{
		{"tx_uid", &payments.TxUid},
		{"tx_id", &payments.TxId},
		{"position_in_payment", payments.PositionInPayment},
		{"amount", &payments.Amount},
		{"asset_id", &payments.AssetId},
		{"height", &payments.Height},
	}
}

func NewUpdatesInput(updates *Updates) UpdatesInput {
	return UpdatesInput{
		"blocks": {
			{"height", &updates.Blocks.Height},
			{"id", &updates.Blocks.Id},
			{"timestamp", &updates.Blocks.Timestamp},
		},
		"txs_1": append(
			NewBaseTxInput(updates.Txs1.Txs),
			proto.InputColumn{"recipient_address", &updates.Txs1.RecipientAddress},
			proto.InputColumn{"amount", &updates.Txs1.Amount},
		),
		"txs_2": append(
			NewBaseTxInput(updates.Txs2.Txs),
			proto.InputColumn{"recipient_address", &updates.Txs2.RecipientAddress},
			proto.InputColumn{"amount", &updates.Txs2.Amount},
		),
		"txs_3": append(
			NewBaseTxInput(updates.Txs3.Txs),
			proto.InputColumn{"asset_id", &updates.Txs3.AssetId},
			proto.InputColumn{"asset_name", &updates.Txs3.AssetName},
			proto.InputColumn{"description", &updates.Txs3.Description},
			proto.InputColumn{"amount", &updates.Txs3.Amount},
			proto.InputColumn{"decimals", &updates.Txs3.Decimals},
			proto.InputColumn{"reissuable", &updates.Txs3.Reissuable},
			proto.InputColumn{"script", &updates.Txs3.Script},
		),
		"txs_4": append(
			NewBaseTxInput(updates.Txs4.Txs),
			proto.InputColumn{"asset_id", &updates.Txs4.AssetId},
			proto.InputColumn{"amount", &updates.Txs4.Amount},
			proto.InputColumn{"recipient_address", &updates.Txs4.RecipientAddress},
			proto.InputColumn{"recipient_alias", &updates.Txs4.RecipientAlias},
			proto.InputColumn{"fee_asset_id", &updates.Txs4.FeeAssetId},
			proto.InputColumn{"attachment", &updates.Txs4.Attachment},
		),
		"txs_5": append(
			NewBaseTxInput(updates.Txs5.Txs),
			proto.InputColumn{"asset_id", &updates.Txs5.AssetId},
			proto.InputColumn{"amount", &updates.Txs5.Amount},
			proto.InputColumn{"reissuable", &updates.Txs5.Reissuable},
		),
		"txs_6": append(
			NewBaseTxInput(updates.Txs6.Txs),
			proto.InputColumn{"asset_id", &updates.Txs6.AssetId},
			proto.InputColumn{"amount", &updates.Txs6.Amount},
		),
		"txs_7": append(
			NewBaseTxInput(updates.Txs7.Txs),
			proto.InputColumn{"amount", &updates.Txs7.Amount},
			proto.InputColumn{"price", &updates.Txs7.Price},
			proto.InputColumn{"amount_asset_id", &updates.Txs7.AmountAssetId},
			proto.InputColumn{"price_asset_id", &updates.Txs7.PriceAssetId},
			proto.InputColumn{"buy_matcher_fee", &updates.Txs7.BuyMatcherFee},
			proto.InputColumn{"sell_matcher_fee", &updates.Txs7.SellMatcherFee},
			proto.InputColumn{"fee_asset_id", &updates.Txs7.FeeAssetId},
		),
		"txs_7_orders": {
			{"tx_uid", &updates.Txs7Orders.TxUid},
			{"tx_id", &updates.Txs7Orders.TxId},
			{"id", &updates.Txs7Orders.Id},
			{"version", &updates.Txs7Orders.Version},
			{"sender", &updates.Txs7Orders.Sender},
			{"sender_public_key", &updates.Txs7Orders.SenderPublicKey},
			{"matcher_public_key", &updates.Txs7Orders.MatcherPublicKey},
			{"amount_asset_id", &updates.Txs7Orders.AmountAssetId},
			{"price_asset_id", &updates.Txs7Orders.PriceAssetId},
			{"order_type", &updates.Txs7Orders.OrderType},
			{"amount", &updates.Txs7Orders.Amount},
			{"price", &updates.Txs7Orders.Price},
			{"timestamp", &updates.Txs7Orders.Timestamp},
			{"expiration", &updates.Txs7Orders.Expiration},
			{"matcher_fee", &updates.Txs7Orders.MatcherFee},
			{"matcher_fee_asset_id", &updates.Txs7Orders.MatcherFeeAssetId},
			{"proofs", updates.Txs7Orders.Proofs},
			{"signature", &updates.Txs7Orders.Signature},
			{"eip712_signature", &updates.Txs7Orders.Eip712Signature},
			{"price_mode", &updates.Txs7Orders.PriceMode},
		},
		"txs_8": append(
			NewBaseTxInput(updates.Txs8.Txs),
			proto.InputColumn{"recipient_address", &updates.Txs8.RecipientAddress},
			proto.InputColumn{"recipient_alias", &updates.Txs8.RecipientAlias},
			proto.InputColumn{"amount", &updates.Txs8.Amount},
		),
		"txs_9": append(
			NewBaseTxInput(updates.Txs9.Txs),
			proto.InputColumn{"lease_tx_uid", &updates.Txs9.LeaseTxUid},
		),
		"txs_10": append(
			NewBaseTxInput(updates.Txs10.Txs),
			proto.InputColumn{"alias", &updates.Txs10.Alias},
		),
		"txs_11": append(
			NewBaseTxInput(updates.Txs11.Txs),
			proto.InputColumn{"asset_id", &updates.Txs11.AssetId},
			proto.InputColumn{"attachment", &updates.Txs11.Attachment},
		),
		"txs_11_transfers": {
			{"tx_uid", &updates.Txs11Transfers.TxUid},
			{"tx_id", &updates.Txs11Transfers.TxId},
			{"position_in_tx", &updates.Txs11Transfers.PositionInTx},
			{"recipient_address", &updates.Txs11Transfers.RecipientAddress},
			{"recipient_alias", &updates.Txs11Transfers.RecipientAlias},
			{"amount", &updates.Txs11Transfers.Amount},
			{"height", &updates.Txs11Transfers.Height},
		},
		"txs_12": NewBaseTxInput(updates.Txs12.Txs),
		"txs_12_data": {
			{"tx_uid", &updates.Txs12Data.TxUid},
			{"tx_id", &updates.Txs12Data.TxId},
			{"position_in_tx", &updates.Txs12Data.PositionInTx},
			{"data_key", &updates.Txs12Data.DataKey},
			{"data_type", &updates.Txs12Data.DataType},
			{"data_value_integer", &updates.Txs12Data.DataValueInteger},
			{"data_value_bool", &updates.Txs12Data.DataValueBoolean},
			{"data_value_binary", &updates.Txs12Data.DataValueBinary},
			{"data_value_string", &updates.Txs12Data.DataValueString},
			{"height", &updates.Txs12Data.Height},
		},
		"txs_13": append(
			NewBaseTxInput(updates.Txs13.Txs),
			proto.InputColumn{"script", &updates.Txs13.Script},
		),
		"txs_14": append(
			NewBaseTxInput(updates.Txs14.Txs),
			proto.InputColumn{"asset_id", &updates.Txs14.AssetId},
			proto.InputColumn{"min_sponsored_asset_fee", &updates.Txs14.MinSponsoredAssetFee},
		),
		"txs_15": append(
			NewBaseTxInput(updates.Txs15.Txs),
			proto.InputColumn{"script", &updates.Txs15.Script},
		),
		"txs_16": append(
			NewBaseTxInput(updates.Txs16.Txs),
			proto.InputColumn{"dapp_address", &updates.Txs16.DappAddress},
			proto.InputColumn{"dapp_alias", &updates.Txs16.DappAlias},
			proto.InputColumn{"function_name", &updates.Txs16.FunctionName},
			proto.InputColumn{"fee_asset_id", &updates.Txs16.FeeAssetId},
		),
		"txs_16_args":     NewTxsArgsInput(updates.Txs16Args),
		"txs_16_payments": NewTxsPaymentsInput(updates.Txs16Payments),
		"txs_17": append(
			NewBaseTxInput(updates.Txs17.Txs),
			proto.InputColumn{"asset_id", &updates.Txs17.AssetId},
			proto.InputColumn{"asset_name", &updates.Txs17.AssetName},
			proto.InputColumn{"description", &updates.Txs17.Description},
		),
		"txs_18": append(
			NewBaseTxInput(updates.Txs18.Txs),
			proto.InputColumn{"bytes", &updates.Txs18.Bytes},
			proto.InputColumn{"function_name", &updates.Txs18.FunctionName},
		),
		"txs_18_args":     NewTxsArgsInput(updates.Txs18Args),
		"txs_18_payments": NewTxsPaymentsInput(updates.Txs18Payments),
		"balances": {
			{"address", &updates.Balances.Address},
			{"asset_id", &updates.Balances.AssetId},
			{"amount", &updates.Balances.Amount},
			{"delta", &updates.Balances.Delta},
			{"tx_uid", &updates.Balances.TxUid},
			{"tx_id", &updates.Balances.TxId},
			{"height", &updates.Balances.Height},
		},
	}
}
