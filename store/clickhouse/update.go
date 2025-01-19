package clickhouse

import (
	"github.com/mr-tron/base58"
	"github.com/slashfast/wavesconsumer/v2/utils"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"time"
)

func (c *Connector) parseUpdate(blockUpdate *events.BlockchainUpdated) error {

	height := blockUpdate.Height
	blockId := base58.Encode(blockUpdate.Id)

	blockAppend := blockUpdate.GetAppend()
	block := blockAppend.GetBlock()
	blockTimestamp := time.UnixMilli(block.Block.Header.Timestamp)

	c.buff.Blocks.Height.Append(height)
	c.buff.Blocks.Id.Append(blockId)
	c.buff.Blocks.Timestamp.Append(blockTimestamp)

	transactions := block.Block.Transactions
	transactionsMetadata := blockAppend.TransactionsMetadata
	txsStateUpdates := blockAppend.TransactionStateUpdates

	var txUids []int64
	var txIds []string
	// txs
	transactionIds := blockAppend.TransactionIds
	for i, commonTx := range transactions {
		proofs := make([]string, len(commonTx.Proofs))
		for i, proof := range commonTx.Proofs {
			proofs[i] = base58.Encode(proof)
		}
		txUid := (UidBase * int64(height)) + int64(i)
		txId := base58.Encode(transactionIds[i])

		txUids = append(txUids, txUid)
		txIds = append(txIds, txId)

		fillBaseTx := func(txs *Txs, status string) {
			txs.Uid.Append(txUid)
			txs.Id.Append(txId)
			txs.Sender.Append(base58.Encode(transactionsMetadata[i].SenderAddress))
			txs.Height.Append(height)
			txs.Proofs.Append(proofs)
			txs.Status.Append(status)
			txs.Signature.Append(proofs[0])
		}

		switch signedTx := commonTx.Transaction.(type) {
		case *waves.SignedTransaction_WavesTransaction:
			wavesTx := signedTx.WavesTransaction
			fillBaseWavesTx := func(txs *Txs, txType int8) {
				txs.Type.Append(txType)
				txs.SenderPublicKey.Append(base58.Encode(wavesTx.SenderPublicKey))
				txs.Timestamp.Append(time.UnixMilli(wavesTx.Timestamp))
				txs.Version.Append(wavesTx.Version)
				txs.Fee.Append(wavesTx.GetFee().GetAmount())
			}
			switch data := wavesTx.GetData().(type) {
			case *waves.Transaction_Genesis:
				txs := c.buff.Txs1
				d := data.Genesis
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.GenesisTransaction))

				txs.RecipientAddress.Append(base58.Encode(d.RecipientAddress))
				txs.Amount.Append(d.Amount)

			case *waves.Transaction_Payment:
				txs := c.buff.Txs2
				d := data.Payment
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.PaymentTransaction))

				txs.RecipientAddress.Append(base58.Encode(d.RecipientAddress))
				txs.Amount.Append(d.Amount)

			case *waves.Transaction_Issue:
				txs := c.buff.Txs3
				d := data.Issue
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.IssueTransaction))

				txs.AssetId.Append(txId)
				txs.AssetName.Append(d.Name)
				txs.Description.Append(d.Description)
				txs.Amount.Append(d.Amount)
				txs.Decimals.Append(d.Decimals)
				txs.Reissuable.Append(d.Reissuable)
				txs.Script.Append(base58.Encode(d.Script))

			case *waves.Transaction_Transfer:
				txs := c.buff.Txs4
				d := data.Transfer
				metadata := transactionsMetadata[i].GetTransfer()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.TransferTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.Amount.AssetId))
				txs.Amount.Append(d.Amount.Amount)
				txs.RecipientAddress.Append(base58.Encode(metadata.RecipientAddress))
				txs.RecipientAlias.Append(d.Recipient.GetAlias())
				txs.FeeAssetId.Append(utils.ExtractAssetId(wavesTx.Fee.AssetId))
				txs.Attachment.Append(base58.Encode(d.Attachment))

			case *waves.Transaction_Reissue:
				txs := c.buff.Txs5
				d := data.Reissue
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.ReissueTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.AssetAmount.AssetId))
				txs.Amount.Append(d.AssetAmount.Amount)
				txs.Reissuable.Append(d.Reissuable)

			case *waves.Transaction_Burn:
				txs := c.buff.Txs6
				d := data.Burn
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.BurnTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.AssetAmount.AssetId))
				txs.Amount.Append(d.AssetAmount.Amount)

			case *waves.Transaction_Exchange:
				txs := c.buff.Txs7
				txsOrders := c.buff.Txs7Orders
				d := data.Exchange
				metadata := transactionsMetadata[i].GetExchange()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.ExchangeTransaction))

				txs.Amount.Append(d.Amount)
				txs.Price.Append(d.Price)
				txs.AmountAssetId.Append(utils.ExtractAssetId(d.Orders[0].AssetPair.AmountAssetId))
				txs.PriceAssetId.Append(utils.ExtractAssetId(d.Orders[0].AssetPair.PriceAssetId))
				txs.BuyMatcherFee.Append(d.BuyMatcherFee)
				txs.SellMatcherFee.Append(d.SellMatcherFee)
				txs.FeeAssetId.Append(utils.ExtractAssetId(wavesTx.Fee.AssetId))

				for i, order := range d.Orders {
					proofs := make([]string, len(order.Proofs))
					for i, proof := range order.Proofs {
						proofs[i] = base58.Encode(proof)
					}
					txsOrders.TxUid.Append(txUid)
					txsOrders.TxId.Append(txId)
					txsOrders.Id.Append(base58.Encode(metadata.OrderIds[i]))
					txsOrders.Version.Append(order.Version)
					txsOrders.Sender.Append(base58.Encode(metadata.OrderSenderAddresses[i]))
					txsOrders.SenderPublicKey.Append(base58.Encode(metadata.OrderSenderPublicKeys[i]))
					txsOrders.MatcherPublicKey.Append(base58.Encode(order.MatcherPublicKey))
					txsOrders.AmountAssetId.Append(utils.ExtractAssetId(order.AssetPair.AmountAssetId))
					txsOrders.PriceAssetId.Append(utils.ExtractAssetId(order.AssetPair.PriceAssetId))
					txsOrders.OrderType.Append(int32(order.OrderSide))
					txsOrders.Amount.Append(order.Amount)
					txsOrders.Price.Append(order.Price)
					txsOrders.Timestamp.Append(time.UnixMilli(order.Timestamp))
					txsOrders.Expiration.Append(time.UnixMilli(order.Expiration))
					txsOrders.MatcherFee.Append(order.MatcherFee.Amount)
					txsOrders.MatcherFeeAssetId.Append(utils.ExtractAssetId(order.MatcherFee.AssetId))
					txsOrders.Proofs.Append(proofs)
					txsOrders.Signature.Append(proofs[0])
					txsOrders.Eip712Signature.Append(base58.Encode(order.GetEip712Signature()))
					txsOrders.PriceMode.Append(int32(order.PriceMode))
				}

			case *waves.Transaction_Lease:
				txs := c.buff.Txs8
				d := data.Lease
				metadata := transactionsMetadata[i].GetLease()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.LeaseTransaction))

				txs.RecipientAddress.Append(base58.Encode(metadata.RecipientAddress))
				txs.RecipientAlias.Append(d.Recipient.GetAlias())
				txs.Amount.Append(d.Amount)

			case *waves.Transaction_LeaseCancel:
				txs := c.buff.Txs9
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.LeaseCancelTransaction))

				txs.LeaseTxUid.Append(txUid)

			case *waves.Transaction_CreateAlias:
				txs := c.buff.Txs10
				d := data.CreateAlias
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.CreateAliasTransaction))

				txs.Alias.Append(d.Alias)

			case *waves.Transaction_MassTransfer:
				txs := c.buff.Txs11
				txsTransfers := c.buff.Txs11Transfers
				d := data.MassTransfer
				metadata := transactionsMetadata[i].GetMassTransfer()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.MassTransferTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.AssetId))
				txs.Attachment.Append(base58.Encode(d.Attachment))

				for i, transfer := range d.Transfers {
					txsTransfers.TxUid.Append(txUid)
					txsTransfers.TxId.Append(txId)
					txsTransfers.PositionInTx.Append(int32(i))
					txsTransfers.RecipientAddress.Append(base58.Encode(metadata.RecipientsAddresses[i]))
					txsTransfers.RecipientAlias.Append(transfer.Recipient.GetAlias())
					txsTransfers.Amount.Append(transfer.Amount)
					txsTransfers.Height.Append(height)
				}

			case *waves.Transaction_DataTransaction:
				txs := c.buff.Txs12
				txsData := c.buff.Txs12Data
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.DataTransaction))

				d := data.DataTransaction
				for i, dataEntry := range d.Data {
					dataType := ""
					switch dataEntry.Value.(type) {
					case *waves.DataTransactionData_DataEntry_IntValue:
						dataType = "integer"
					case *waves.DataTransactionData_DataEntry_BoolValue:
						dataType = "boolean"
					case *waves.DataTransactionData_DataEntry_BinaryValue:
						dataType = "binary"
					case *waves.DataTransactionData_DataEntry_StringValue:
						dataType = "string"
					}
					txsData.TxUid.Append(txUid)
					txsData.TxId.Append(txId)
					txsData.PositionInTx.Append(int32(i))
					txsData.DataKey.Append(utils.EscapeUnicodeNull(dataEntry.Key))
					txsData.DataType.Append(dataType)
					txsData.DataValueInteger.Append(dataEntry.GetIntValue())
					txsData.DataValueBoolean.Append(dataEntry.GetBoolValue())
					txsData.DataValueBinary.Append(utils.ToPrefixedBase64(dataEntry.GetBinaryValue()))
					txsData.DataValueString.Append(dataEntry.GetStringValue())
					txsData.Height.Append(height)
				}
			case *waves.Transaction_SetScript:
				txs := c.buff.Txs13
				d := data.SetScript
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SetScriptTransaction))

				txs.Script.Append(utils.ToPrefixedBase64(d.Script))

			case *waves.Transaction_SponsorFee:
				txs := c.buff.Txs14
				d := data.SponsorFee
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SponsorshipTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.MinFee.AssetId))
				txs.MinSponsoredAssetFee.Append(d.MinFee.Amount)

			case *waves.Transaction_SetAssetScript:
				txs := c.buff.Txs15
				d := data.SetAssetScript
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SetAssetScriptTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.AssetId))
				txs.Script.Append(utils.ToPrefixedBase64(d.Script))

			case *waves.Transaction_InvokeScript:
				txs := c.buff.Txs16
				txsArgs := c.buff.Txs16Args
				txsPayments := c.buff.Txs16Payments
				d := data.InvokeScript
				metadata := transactionsMetadata[i].GetInvokeScript()
				txStatus := "success"
				if metadata.Result.GetErrorMessage() != nil {
					txStatus = "script_execution_failed"
				}

				fillBaseTx(txs.Txs, txStatus)
				fillBaseWavesTx(txs.Txs, int8(proto.InvokeScriptTransaction))

				txs.DappAddress.Append(base58.Encode(metadata.DAppAddress))
				txs.DappAlias.Append(d.DApp.GetAlias())
				txs.FunctionName.Append(metadata.FunctionName)
				txs.FeeAssetId.Append(utils.ExtractAssetId(wavesTx.Fee.AssetId))

				for i, arg := range metadata.Arguments {
					err := AppendTxArg(txsArgs, arg, txId, txUid, int32(i), height)
					if err != nil {
						return err
					}
				}

				for i, payment := range metadata.Payments {
					AppendTxPayment(txsPayments, payment, txId, txUid, int32(i), height)
				}

			case *waves.Transaction_UpdateAssetInfo:
				txs := c.buff.Txs17
				d := data.UpdateAssetInfo
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.UpdateAssetInfoTransaction))

				txs.AssetId.Append(utils.ExtractAssetId(d.AssetId))
				txs.AssetName.Append(utils.EscapeUnicodeNull(d.Name))
				txs.Description.Append(utils.EscapeUnicodeNull(d.Description))

			}
		case *waves.SignedTransaction_EthereumTransaction:
			txs := c.buff.Txs18
			metadata := transactionsMetadata[i].GetEthereum()

			txs.Txs.Type.Append(int8(proto.EthereumMetamaskTransaction))
			txs.Txs.SenderPublicKey.Append(base58.Encode(metadata.SenderPublicKey))
			txs.Txs.Timestamp.Append(time.UnixMilli(metadata.Timestamp))
			txs.Txs.Version.Append(1)
			txs.Txs.Fee.Append(metadata.Fee)

			txs.Bytes.Append(utils.ToPrefixedBase64(signedTx.EthereumTransaction))

			switch actionMetadata := metadata.Action.(type) {
			case *events.TransactionMetadata_EthereumMetadata_Invoke:
				if actionMetadata.Invoke.Result.GetErrorMessage() != nil {
					txs.Txs.Status.Append("script_execution_failed")
				} else {
					txs.Txs.Status.Append("success")
				}

				for i, arg := range actionMetadata.Invoke.Arguments {
					err := AppendTxArg(c.buff.Txs18Args, arg, txId, txUid, int32(i), height)
					if err != nil {
						return err
					}

				}

				for i, payment := range actionMetadata.Invoke.Payments {
					AppendTxPayment(c.buff.Txs18Payments, payment, txId, txUid, int32(i), height)
				}
			}

		}
	}

	balances := c.buff.Balances

	for i, su := range txsStateUpdates {
		for _, bu := range su.Balances {
			delta := bu.AmountAfter.Amount - bu.AmountBefore

			balances.Address.Append(base58.Encode(bu.Address))
			balances.AssetId.Append(utils.ExtractAssetId(bu.AmountAfter.AssetId))
			balances.Amount.Append(bu.AmountAfter.Amount)
			balances.Delta.Append(delta)
			balances.TxUid.Append(txUids[i])
			balances.TxId.Append(txIds[i])
			balances.Height.Append(height)
		}
	}

	return nil
}
