package clickhouse

import (
	"github.com/mr-tron/base58"
	"github.com/slashfast/wavesconsumer/v2/pkg/utils"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"time"
)

func (c *Connector) parseUpdate(blockUpdate *events.BlockchainUpdated) (*Updates, error) {
	updates := NewUpdates()
	height := blockUpdate.Height
	blockId := base58.Encode(blockUpdate.Id)

	blockAppend := blockUpdate.GetAppend()
	block := blockAppend.GetBlock()
	blockTimestamp := time.UnixMilli(block.Block.Header.Timestamp)

	updates.Blocks.Height = append(updates.Blocks.Height, height)
	updates.Blocks.Id = append(updates.Blocks.Id, blockId)
	updates.Blocks.Timestamp = append(updates.Blocks.Timestamp, blockTimestamp)

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

		signature := ""
		if len(proofs) > 0 {
			signature = proofs[0]
		}

		txUid := (UidBase * int64(height)) + int64(i)
		txId := base58.Encode(transactionIds[i])

		txUids = append(txUids, txUid)
		txIds = append(txIds, txId)

		fillBaseTx := func(txs *Txs, status string) {
			txs.Uid = append(txs.Uid, txUid)
			txs.Id = append(txs.Id, txId)
			txs.Sender = append(txs.Sender, base58.Encode(transactionsMetadata[i].SenderAddress))
			txs.Height = append(txs.Height, height)
			txs.Proofs = append(txs.Proofs, proofs)
			txs.Status = append(txs.Status, status)
			txs.Signature = append(txs.Signature, signature)
		}

		switch signedTx := commonTx.Transaction.(type) {
		case *waves.SignedTransaction_WavesTransaction:
			wavesTx := signedTx.WavesTransaction
			fillBaseWavesTx := func(txs *Txs, txType int8) {
				txs.Type = append(txs.Type, txType)
				txs.SenderPublicKey = append(txs.SenderPublicKey, base58.Encode(wavesTx.SenderPublicKey))
				txs.Timestamp = append(txs.Timestamp, time.UnixMilli(wavesTx.Timestamp))
				txs.Version = append(txs.Version, wavesTx.Version)
				txs.Fee = append(txs.Fee, wavesTx.GetFee().GetAmount())
			}
			switch data := wavesTx.GetData().(type) {
			case *waves.Transaction_Genesis:
				txs := updates.Txs1
				d := data.Genesis
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.GenesisTransaction))

				txs.RecipientAddress = append(txs.RecipientAddress, utils.MustAddressFromBody(d.RecipientAddress))
				txs.Amount = append(txs.Amount, d.Amount)

			case *waves.Transaction_Payment:
				txs := updates.Txs2
				d := data.Payment
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.PaymentTransaction))

				txs.RecipientAddress = append(txs.RecipientAddress, utils.MustAddressFromBody(d.RecipientAddress))
				txs.Amount = append(txs.Amount, d.Amount)

			case *waves.Transaction_Issue:
				txs := updates.Txs3
				d := data.Issue
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.IssueTransaction))

				txs.AssetId = append(txs.AssetId, txId)
				txs.AssetName = append(txs.AssetName, d.Name)
				txs.Description = append(txs.Description, d.Description)
				txs.Amount = append(txs.Amount, d.Amount)
				txs.Decimals = append(txs.Decimals, d.Decimals)
				txs.Reissuable = append(txs.Reissuable, d.Reissuable)
				txs.Script = append(txs.Script, base58.Encode(d.Script))

			case *waves.Transaction_Transfer:
				txs := updates.Txs4
				d := data.Transfer
				metadata := transactionsMetadata[i].GetTransfer()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.TransferTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.Amount.AssetId))
				txs.Amount = append(txs.Amount, d.Amount.Amount)
				txs.RecipientAddress = append(txs.RecipientAddress, utils.MustAddressFromBody(metadata.RecipientAddress))
				txs.RecipientAlias = append(txs.RecipientAlias, d.Recipient.GetAlias())
				txs.FeeAssetId = append(txs.FeeAssetId, utils.ExtractAssetId(wavesTx.Fee.AssetId))
				txs.Attachment = append(txs.Attachment, base58.Encode(d.Attachment))

			case *waves.Transaction_Reissue:
				txs := updates.Txs5
				d := data.Reissue
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.ReissueTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.AssetAmount.AssetId))
				txs.Amount = append(txs.Amount, d.AssetAmount.Amount)
				txs.Reissuable = append(txs.Reissuable, d.Reissuable)

			case *waves.Transaction_Burn:
				txs := updates.Txs6
				d := data.Burn
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.BurnTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.AssetAmount.AssetId))
				txs.Amount = append(txs.Amount, d.AssetAmount.Amount)

			case *waves.Transaction_Exchange:
				txs := updates.Txs7
				txsOrders := updates.Txs7Orders
				d := data.Exchange
				metadata := transactionsMetadata[i].GetExchange()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.ExchangeTransaction))

				txs.Amount = append(txs.Amount, d.Amount)
				txs.Price = append(txs.Price, d.Price)
				txs.AmountAssetId = append(txs.AmountAssetId, utils.ExtractAssetId(d.Orders[0].AssetPair.AmountAssetId))
				txs.PriceAssetId = append(txs.PriceAssetId, utils.ExtractAssetId(d.Orders[0].AssetPair.PriceAssetId))
				txs.BuyMatcherFee = append(txs.BuyMatcherFee, d.BuyMatcherFee)
				txs.SellMatcherFee = append(txs.SellMatcherFee, d.SellMatcherFee)
				txs.FeeAssetId = append(txs.FeeAssetId, utils.ExtractAssetId(wavesTx.Fee.AssetId))

				for i, order := range d.Orders {
					proofs := make([]string, len(order.Proofs))
					for i, proof := range order.Proofs {
						proofs[i] = base58.Encode(proof)
					}

					signature := ""
					if len(proofs) > 0 {
						signature = proofs[0]
					}

					txsOrders.TxUid = append(txsOrders.TxUid, txUid)
					txsOrders.TxId = append(txsOrders.TxId, txId)
					txsOrders.Id = append(txsOrders.Id, base58.Encode(metadata.OrderIds[i]))
					txsOrders.Version = append(txsOrders.Version, order.Version)
					txsOrders.Sender = append(txsOrders.Sender, base58.Encode(metadata.OrderSenderAddresses[i]))
					txsOrders.SenderPublicKey = append(txsOrders.SenderPublicKey, base58.Encode(metadata.OrderSenderPublicKeys[i]))
					txsOrders.MatcherPublicKey = append(txsOrders.MatcherPublicKey, base58.Encode(order.MatcherPublicKey))
					txsOrders.AmountAssetId = append(txsOrders.AmountAssetId, utils.ExtractAssetId(order.AssetPair.AmountAssetId))
					txsOrders.PriceAssetId = append(txsOrders.PriceAssetId, utils.ExtractAssetId(order.AssetPair.PriceAssetId))
					txsOrders.OrderType = append(txsOrders.OrderType, int32(order.OrderSide))
					txsOrders.Amount = append(txsOrders.Amount, order.Amount)
					txsOrders.Price = append(txsOrders.Price, order.Price)
					txsOrders.Timestamp = append(txsOrders.Timestamp, time.UnixMilli(order.Timestamp))
					txsOrders.Expiration = append(txsOrders.Expiration, time.UnixMilli(order.Expiration))
					txsOrders.MatcherFee = append(txsOrders.MatcherFee, order.MatcherFee.Amount)
					txsOrders.MatcherFeeAssetId = append(txsOrders.MatcherFeeAssetId, utils.ExtractAssetId(order.MatcherFee.AssetId))
					txsOrders.Proofs = append(txsOrders.Proofs, proofs)
					txsOrders.Signature = append(txsOrders.Signature, signature)
					txsOrders.Eip712Signature = append(txsOrders.Eip712Signature, base58.Encode(order.GetEip712Signature()))
					txsOrders.PriceMode = append(txsOrders.PriceMode, int32(order.PriceMode))
					txsOrders.Height = append(txsOrders.Height, height)
				}

			case *waves.Transaction_Lease:
				txs := updates.Txs8
				d := data.Lease
				metadata := transactionsMetadata[i].GetLease()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.LeaseTransaction))

				txs.RecipientAddress = append(txs.RecipientAddress, utils.MustAddressFromBody(metadata.RecipientAddress))
				txs.RecipientAlias = append(txs.RecipientAlias, d.Recipient.GetAlias())
				txs.Amount = append(txs.Amount, d.Amount)

			case *waves.Transaction_LeaseCancel:
				txs := updates.Txs9
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.LeaseCancelTransaction))

				txs.LeaseTxUid = append(txs.LeaseTxUid, txUid)

			case *waves.Transaction_CreateAlias:
				txs := updates.Txs10
				d := data.CreateAlias
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.CreateAliasTransaction))

				txs.Alias = append(txs.Alias, d.Alias)

			case *waves.Transaction_MassTransfer:
				txs := updates.Txs11
				txsTransfers := updates.Txs11Transfers
				d := data.MassTransfer
				metadata := transactionsMetadata[i].GetMassTransfer()
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.MassTransferTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.AssetId))
				txs.Attachment = append(txs.Attachment, base58.Encode(d.Attachment))

				for i, transfer := range d.Transfers {
					txsTransfers.TxUid = append(txsTransfers.TxUid, txUid)
					txsTransfers.TxId = append(txsTransfers.TxId, txId)
					txsTransfers.PositionInTx = append(txsTransfers.PositionInTx, int32(i))
					txsTransfers.RecipientAddress = append(txsTransfers.RecipientAddress, utils.MustAddressFromBody(metadata.RecipientsAddresses[i]))
					txsTransfers.RecipientAlias = append(txsTransfers.RecipientAlias, transfer.Recipient.GetAlias())
					txsTransfers.Amount = append(txsTransfers.Amount, transfer.Amount)
					txsTransfers.Height = append(txsTransfers.Height, height)
				}

			case *waves.Transaction_DataTransaction:
				txs := updates.Txs12
				txsData := updates.Txs12Data
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
					txsData.TxUid = append(txsData.TxUid, txUid)
					txsData.TxId = append(txsData.TxId, txId)
					txsData.PositionInTx = append(txsData.PositionInTx, int32(i))
					txsData.DataKey = append(txsData.DataKey, utils.EscapeUnicodeNull(dataEntry.Key))
					txsData.DataType = append(txsData.DataType, dataType)
					txsData.DataValueInteger = append(txsData.DataValueInteger, dataEntry.GetIntValue())
					txsData.DataValueBoolean = append(txsData.DataValueBoolean, dataEntry.GetBoolValue())
					txsData.DataValueBinary = append(txsData.DataValueBinary, utils.ToBase64(dataEntry.GetBinaryValue()))
					txsData.DataValueString = append(txsData.DataValueString, dataEntry.GetStringValue())
					txsData.Height = append(txsData.Height, height)
				}
			case *waves.Transaction_SetScript:
				txs := updates.Txs13
				d := data.SetScript
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SetScriptTransaction))

				txs.Script = append(txs.Script, utils.ToBase64(d.Script))

			case *waves.Transaction_SponsorFee:
				txs := updates.Txs14
				d := data.SponsorFee
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SponsorshipTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.MinFee.AssetId))
				txs.MinSponsoredAssetFee = append(txs.MinSponsoredAssetFee, d.MinFee.Amount)

			case *waves.Transaction_SetAssetScript:
				txs := updates.Txs15
				d := data.SetAssetScript
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.SetAssetScriptTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.AssetId))
				txs.Script = append(txs.Script, utils.ToBase64(d.Script))

			case *waves.Transaction_InvokeScript:
				txs := updates.Txs16
				txsArgs := updates.Txs16Args
				txsPayments := updates.Txs16Payments
				d := data.InvokeScript
				metadata := transactionsMetadata[i].GetInvokeScript()
				txStatus := "success"
				if metadata.Result.GetErrorMessage() != nil {
					txStatus = "script_execution_failed"
				}

				fillBaseTx(txs.Txs, txStatus)
				fillBaseWavesTx(txs.Txs, int8(proto.InvokeScriptTransaction))

				txs.DappAddress = append(txs.DappAddress, base58.Encode(metadata.DAppAddress))
				txs.DappAlias = append(txs.DappAlias, d.DApp.GetAlias())
				txs.FunctionName = append(txs.FunctionName, metadata.FunctionName)
				txs.FeeAssetId = append(txs.FeeAssetId, utils.ExtractAssetId(wavesTx.Fee.AssetId))

				for i, arg := range metadata.Arguments {
					err := AppendTxArg(txsArgs, arg, txId, txUid, int32(i), height)
					if err != nil {
						return nil, err
					}
				}

				for i, payment := range metadata.Payments {
					AppendTxPayment(txsPayments, payment, txId, txUid, int32(i), height)
				}

			case *waves.Transaction_UpdateAssetInfo:
				txs := updates.Txs17
				d := data.UpdateAssetInfo
				fillBaseTx(txs.Txs, "success")
				fillBaseWavesTx(txs.Txs, int8(proto.UpdateAssetInfoTransaction))

				txs.AssetId = append(txs.AssetId, utils.ExtractAssetId(d.AssetId))
				txs.AssetName = append(txs.AssetName, utils.EscapeUnicodeNull(d.Name))
				txs.Description = append(txs.Description, utils.EscapeUnicodeNull(d.Description))

			}
		case *waves.SignedTransaction_EthereumTransaction:
			txs := updates.Txs18
			txStatus := "success"
			funcName := ""
			metadata := transactionsMetadata[i].GetEthereum()
			txs.Txs.Type = append(txs.Txs.Type, int8(proto.EthereumMetamaskTransaction))
			txs.Txs.SenderPublicKey = append(txs.Txs.SenderPublicKey, base58.Encode(metadata.SenderPublicKey))
			txs.Txs.Timestamp = append(txs.Txs.Timestamp, time.UnixMilli(metadata.Timestamp))
			txs.Txs.Version = append(txs.Txs.Version, 1)
			txs.Txs.Fee = append(txs.Txs.Fee, metadata.Fee)

			txs.Bytes = append(txs.Bytes, utils.ToBase64(signedTx.EthereumTransaction))

			switch actionMetadata := metadata.Action.(type) {
			case *events.TransactionMetadata_EthereumMetadata_Invoke:
				if actionMetadata.Invoke.Result.GetErrorMessage() != nil {
					txStatus = "script_execution_failed"
				}
				funcName = actionMetadata.Invoke.FunctionName

				for i, arg := range actionMetadata.Invoke.Arguments {
					err := AppendTxArg(updates.Txs18Args, arg, txId, txUid, int32(i), height)
					if err != nil {
						return nil, err
					}

				}

				for i, payment := range actionMetadata.Invoke.Payments {
					AppendTxPayment(updates.Txs18Payments, payment, txId, txUid, int32(i), height)
				}
			}
			txs.FunctionName = append(txs.FunctionName, funcName)
			fillBaseTx(txs.Txs, txStatus)
		}
	}

	balances := updates.Balances

	for i, su := range txsStateUpdates {
		for _, bu := range su.Balances {
			delta := bu.AmountAfter.Amount - bu.AmountBefore

			balances.Address = append(balances.Address, base58.Encode(bu.Address))
			balances.AssetId = append(balances.AssetId, utils.ExtractAssetId(bu.AmountAfter.AssetId))
			balances.Amount = append(balances.Amount, bu.AmountAfter.Amount)
			balances.Delta = append(balances.Delta, delta)
			balances.TxUid = append(balances.TxUid, txUids[i])
			balances.TxId = append(balances.TxId, txIds[i])
			balances.Height = append(balances.Height, height)
		}
	}

	return updates, nil
}
