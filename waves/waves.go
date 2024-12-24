package waves

import (
	"encoding/json"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"log"
	"net/http"
	"reflect"
	"time"
)

var referencedAssets = make(map[string]ReferencedAsset)
var RestApiHost string
var httpClient = http.Client{}

type BalanceUpdate struct {
	Height         uint32
	AssetId        string
	Address        string
	TxVersion      uint8
	TxType         uint8
	TxId           string
	Counterparty   string
	Subunits       uint64
	DeltaSubunits  int64
	AssetDecimals  uint8
	AssetName      string
	IsInbound      bool
	IsStateUpdate  bool
	TxTimestamp    time.Time
	BlockTimestamp time.Time
}

type AssetDetails struct {
	AssetId              string `json:"assetId"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Quantity             int64  `json:"amount"`
	Decimals             int    `json:"decimals"`
	IssueHeight          int    `json:"issueHeight"`
	IssueTimestamp       int64  `json:"issueTimestamp"`
	Issuer               string `json:"issuer"`
	IssuerPublicKey      string `json:"issuerPublicKey"`
	Reissuable           bool   `json:"reissuable"`
	Scripted             bool   `json:"scripted"`
	SequenceInBlock      int    `json:"sequenceInBlock"`
	MinSponsoredAssetFee int64  `json:"minSponsoredAssetFee"`
	OriginTransactionId  string `json:"originTransactionId"`
	ScriptDetails        struct {
		ScriptComplexity int    `json:"scriptComplexity"`
		Script           string `json:"script"`
		ScriptText       string `json:"scriptText"`
	} `json:"scriptDetails"`
}

type ReferencedAsset struct {
	Id       string
	Decimals uint8
	Name     string
}

type BlocksHeightResponse struct {
	Height int32 `json:"height"`
}

func toBase58String(addr []byte) string {
	return base58.Encode(addr)
}

func decomposeBalanceUpdate(
	bu *events.StateUpdate_BalanceUpdate,
) (string, ReferencedAsset, uint64, int64, bool) {
	var asset ReferencedAsset

	if assetId := bu.AmountAfter.AssetId; assetId != nil {
		b58assetId := toBase58String(assetId)
		var ok bool
		asset, ok = referencedAssets[b58assetId]
		if !ok {
			asset = getAssetById(b58assetId)
			referencedAssets[b58assetId] = asset
		}
	} else {
		asset = Asset
	}

	address := toBase58String(bu.Address)
	beforeSubunits := bu.AmountBefore
	afterSubunits := bu.AmountAfter.Amount
	deltaSubunits := afterSubunits - beforeSubunits
	isInbound := deltaSubunits > 0

	return address, asset, uint64(afterSubunits), deltaSubunits, isInbound
}

func getUpdate(updateResponse *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
	var insertions []BalanceUpdate
	update := updateResponse.GetUpdate()
	appendField := update.GetAppend()
	block := appendField.GetBlock()
	height := uint32(update.Height)
	blockTimestamp := time.Unix(0, block.Block.Header.Timestamp*int64(time.Millisecond)).UTC()
	stateUpdatesBalances := appendField.StateUpdate.GetBalances()
	transactionIds := appendField.TransactionIds
	transactionStateUpdates := appendField.TransactionStateUpdates
	orderedReferencedAssets := update.ReferencedAssets

	for _, ra := range orderedReferencedAssets {
		id := toBase58String(ra.Id)
		if _, ok := referencedAssets[id]; !ok {
			referencedAssets[id] = ReferencedAsset{
				Id:       id,
				Decimals: uint8(ra.Decimals),
				Name:     ra.Name,
			}
		}
	}

	for _, bu := range stateUpdatesBalances {
		address, asset, subunits, deltaSubunits, isInbound := decomposeBalanceUpdate(bu)
		insertions = append(insertions, BalanceUpdate{
			Height:         height,
			AssetId:        asset.Id,
			Address:        address,
			TxVersion:      0,
			TxType:         0,
			TxId:           "",
			Counterparty:   "",
			Subunits:       subunits,
			DeltaSubunits:  deltaSubunits,
			AssetDecimals:  asset.Decimals,
			AssetName:      asset.Name,
			IsInbound:      isInbound,
			IsStateUpdate:  true,
			TxTimestamp:    blockTimestamp,
			BlockTimestamp: blockTimestamp,
		})
	}

	for i, tsu := range transactionStateUpdates {
		txId := toBase58String(transactionIds[i])

		var timestamp time.Time
		var txType uint8
		var txVersion uint8
		var counterparty string

		if txData := block.Block.Transactions[i].GetWavesTransaction(); txData != nil {
			timestamp = time.Unix(0, txData.Timestamp*int64(time.Millisecond)).UTC()
			txVersion = uint8(txData.Version)
			txType = TxTypes[reflect.TypeOf(txData.Data).Elem().Name()]
			txMetadata := appendField.TransactionsMetadata[i]
			counterparty = toBase58String(txMetadata.SenderAddress)
		} else {
			txMetadata := appendField.TransactionsMetadata[i]
			txData := txMetadata.GetEthereum()
			timestamp = time.Unix(0, txData.Timestamp*int64(time.Millisecond)).UTC()
			txVersion = 0
			txType = TxTypes["Transaction_Ethereum"]
			counterparty = toBase58String(txMetadata.SenderAddress)
		}

		for _, bu := range tsu.Balances {
			address, asset, subunits, deltaSubunits, isInbound := decomposeBalanceUpdate(bu)

			insertions = append(insertions, BalanceUpdate{
				Height:         height,
				AssetId:        asset.Id,
				Address:        address,
				TxVersion:      txVersion,
				TxType:         txType,
				TxId:           txId,
				Counterparty:   counterparty,
				Subunits:       subunits,
				DeltaSubunits:  deltaSubunits,
				AssetDecimals:  asset.Decimals,
				AssetName:      asset.Name,
				IsInbound:      isInbound,
				IsStateUpdate:  false,
				TxTimestamp:    timestamp,
				BlockTimestamp: blockTimestamp,
			})
		}
	}

	return insertions
}

func GetUpdates(ch chan *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
	var updates []BalanceUpdate
	for u := range ch {
		updates = append(updates, getUpdate(u)...)
	}
	return updates
}

func getAssetById(assetId string) ReferencedAsset {
	r, err := httpClient.Get(RestApiHost + "/assets/details/" + assetId + "?full=true")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var target AssetDetails

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		log.Fatal(err)
	}

	return ReferencedAsset{
		Id:       target.AssetId,
		Decimals: uint8(target.Decimals),
		Name:     target.Name,
	}
}

func GetCurrentHeight() int32 {
	r, err := httpClient.Get(RestApiHost + "/blocks/height")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var target BlocksHeightResponse

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		log.Fatal(err)
	}

	return target.Height - 1
}
