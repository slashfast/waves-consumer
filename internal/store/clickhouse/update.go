package clickhouse

import (
	"time"

	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	"github.com/wavesplatform/gowaves/pkg/proto"

	"github.com/slashfast/wavesconsumer/v2/pkg/utils"
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

				// Parse nested orders (always 2: buy and sell)
				orderId := make([]string, len(d.Orders))
				orderVersion := make([]int32, len(d.Orders))
				orderSender := make([]string, len(d.Orders))
				orderSenderPublicKey := make([]string, len(d.Orders))
				orderMatcherPublicKey := make([]string, len(d.Orders))
				orderAmountAssetId := make([]string, len(d.Orders))
				orderPriceAssetId := make([]string, len(d.Orders))
				orderOrderType := make([]int32, len(d.Orders))
				orderAmount := make([]int64, len(d.Orders))
				orderPrice := make([]int64, len(d.Orders))
				orderTimestamp := make([]time.Time, len(d.Orders))
				orderExpiration := make([]time.Time, len(d.Orders))
				orderMatcherFee := make([]int64, len(d.Orders))
				orderMatcherFeeAssetId := make([]string, len(d.Orders))
				orderProofs := make([][]string, len(d.Orders))
				orderSignature := make([]string, len(d.Orders))
				orderEip712Signature := make([]string, len(d.Orders))
				orderPriceMode := make([]int32, len(d.Orders))

				for i, order := range d.Orders {
					proofs := make([]string, len(order.Proofs))
					for j, proof := range order.Proofs {
						proofs[j] = base58.Encode(proof)
					}

					signature := ""
					if len(proofs) > 0 {
						signature = proofs[0]
					}

					orderId[i] = base58.Encode(metadata.OrderIds[i])
					orderVersion[i] = order.Version
					orderSender[i] = base58.Encode(metadata.OrderSenderAddresses[i])
					orderSenderPublicKey[i] = base58.Encode(metadata.OrderSenderPublicKeys[i])
					orderMatcherPublicKey[i] = base58.Encode(order.MatcherPublicKey)
					orderAmountAssetId[i] = utils.ExtractAssetId(order.AssetPair.AmountAssetId)
					orderPriceAssetId[i] = utils.ExtractAssetId(order.AssetPair.PriceAssetId)
					orderOrderType[i] = int32(order.OrderSide)
					orderAmount[i] = order.Amount
					orderPrice[i] = order.Price
					orderTimestamp[i] = time.UnixMilli(order.Timestamp)
					orderExpiration[i] = time.UnixMilli(order.Expiration)
					orderMatcherFee[i] = order.MatcherFee.Amount
					orderMatcherFeeAssetId[i] = utils.ExtractAssetId(order.MatcherFee.AssetId)
					orderProofs[i] = proofs
					orderSignature[i] = signature
					orderEip712Signature[i] = base58.Encode(order.GetEip712Signature())
					orderPriceMode[i] = int32(order.PriceMode)
				}

				txs.OrdersId = append(txs.OrdersId, orderId)
				txs.OrdersVersion = append(txs.OrdersVersion, orderVersion)
				txs.OrdersSender = append(txs.OrdersSender, orderSender)
				txs.OrdersSenderPublicKey = append(txs.OrdersSenderPublicKey, orderSenderPublicKey)
				txs.OrdersMatcherPublicKey = append(txs.OrdersMatcherPublicKey, orderMatcherPublicKey)
				txs.OrdersAmountAssetId = append(txs.OrdersAmountAssetId, orderAmountAssetId)
				txs.OrdersPriceAssetId = append(txs.OrdersPriceAssetId, orderPriceAssetId)
				txs.OrdersOrderType = append(txs.OrdersOrderType, orderOrderType)
				txs.OrdersAmount = append(txs.OrdersAmount, orderAmount)
				txs.OrdersPrice = append(txs.OrdersPrice, orderPrice)
				txs.OrdersTimestamp = append(txs.OrdersTimestamp, orderTimestamp)
				txs.OrdersExpiration = append(txs.OrdersExpiration, orderExpiration)
				txs.OrdersMatcherFee = append(txs.OrdersMatcherFee, orderMatcherFee)
				txs.OrdersMatcherFeeAssetId = append(txs.OrdersMatcherFeeAssetId, orderMatcherFeeAssetId)
				txs.OrdersProofs = append(txs.OrdersProofs, orderProofs)
				txs.OrdersSignature = append(txs.OrdersSignature, orderSignature)
				txs.OrdersEip712Signature = append(txs.OrdersEip712Signature, orderEip712Signature)
				txs.OrdersPriceMode = append(txs.OrdersPriceMode, orderPriceMode)

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
	dataStorage := updates.DataStorage

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
		for _, deu := range su.DataEntries {
			de := deu.DataEntry
			dataType := ""
			switch de.Value.(type) {
			case *waves.DataTransactionData_DataEntry_IntValue:
				dataType = "integer"
			case *waves.DataTransactionData_DataEntry_BoolValue:
				dataType = "boolean"
			case *waves.DataTransactionData_DataEntry_BinaryValue:
				dataType = "binary"
			case *waves.DataTransactionData_DataEntry_StringValue:
				dataType = "string"
			}

			dataStorage.Address = append(dataStorage.Address, base58.Encode(deu.Address))
			dataStorage.DataKey = append(dataStorage.DataKey, utils.EscapeUnicodeNull(de.Key))
			dataStorage.DataType = append(dataStorage.DataType, dataType)
			dataStorage.DataValueInteger = append(dataStorage.DataValueInteger, de.GetIntValue())
			dataStorage.DataValueBoolean = append(dataStorage.DataValueBoolean, de.GetBoolValue())
			dataStorage.DataValueBinary = append(dataStorage.DataValueBinary, utils.ToBase64(de.GetBinaryValue()))
			dataStorage.DataValueString = append(dataStorage.DataValueString, de.GetStringValue())
			dataStorage.TxUid = append(dataStorage.TxUid, txUids[i])
			dataStorage.TxId = append(dataStorage.TxId, txIds[i])
			dataStorage.Height = append(dataStorage.Height, height)
		}
	}

	return updates, nil
}
