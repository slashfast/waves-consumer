package waves

import (
	"context"
	"encoding/json"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	blockchainupdates "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"math"
	"net/http"
	"reflect"
	"sync"
	"time"
)

type NodeClientConfig interface {
	GetNodeURL() string
	GetGRPCAddress() string
}

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

type UpdatesClient struct {
	conn             blockchainupdates.BlockchainUpdatesApiClient
	rest             *http.Client
	nodeURL          string
	referencedAssets map[string]ReferencedAsset
	log              *zerolog.Logger
}

func (uc *UpdatesClient) GetCurrentHeight() int32 {
	r, err := uc.rest.Get(uc.nodeURL + "/blocks/height")
	if err != nil {
		uc.log.Fatal().Err(err).Msg("")
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			uc.log.Fatal().Err(err).Msg("")
		}
	}()

	var target BlocksHeightResponse

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		uc.log.Fatal().Err(err).Msg("")
	}

	return target.Height - 1
}

func (uc *UpdatesClient) getAssetById(assetId string) ReferencedAsset {
	r, err := uc.rest.Get(uc.nodeURL + "/assets/details/" + assetId + "?full=true")
	if err != nil {
		uc.log.Fatal().Err(err).Msg("")
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			uc.log.Fatal().Err(err).Msg("")
		}
	}()

	var target AssetDetails

	if err := json.NewDecoder(r.Body).Decode(&target); err != nil {
		uc.log.Fatal().Err(err).Msg("")
	}

	return ReferencedAsset{
		Id:       target.AssetId,
		Decimals: uint8(target.Decimals),
		Name:     target.Name,
	}
}

func (uc *UpdatesClient) decomposeBalanceUpdate(bu *events.StateUpdate_BalanceUpdate) (string, ReferencedAsset, uint64, int64, bool) {
	var asset ReferencedAsset

	if assetId := bu.AmountAfter.AssetId; assetId != nil {
		b58assetId := toBase58String(assetId)
		var ok bool
		asset, ok = uc.referencedAssets[b58assetId]
		if !ok {
			asset = uc.getAssetById(b58assetId)
			uc.referencedAssets[b58assetId] = asset
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

func (uc *UpdatesClient) getUpdate(updateResponse *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
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
		if _, ok := uc.referencedAssets[id]; !ok {
			uc.referencedAssets[id] = ReferencedAsset{
				Id:       id,
				Decimals: uint8(ra.Decimals),
				Name:     ra.Name,
			}
		}
	}

	for _, bu := range stateUpdatesBalances {
		address, asset, subunits, deltaSubunits, isInbound := uc.decomposeBalanceUpdate(bu)
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
			address, asset, subunits, deltaSubunits, isInbound := uc.decomposeBalanceUpdate(bu)

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

func (uc *UpdatesClient) GetUpdates(ch chan *blockchainupdates.GetBlockUpdateResponse) []BalanceUpdate {
	var updates []BalanceUpdate
	for u := range ch {
		updates = append(updates, uc.getUpdate(u)...)
	}
	return updates
}

func (uc *UpdatesClient) GetBlocks(ch chan *blockchainupdates.GetBlockUpdateResponse, from int32, len int32) {
	var wg sync.WaitGroup
	toHeight := from + len - 1
	for height := from; height <= toHeight; height++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				resp, err := uc.conn.GetBlockUpdate(
					context.Background(),
					&blockchainupdates.GetBlockUpdateRequest{Height: height},
				)
				if err == nil {
					ch <- resp
					return
				}
			}
		}()
	}

	wg.Wait()
	close(ch)
}

func NewUpdatesClient(lc fx.Lifecycle, l *zerolog.Logger, config NodeClientConfig) (*UpdatesClient, error) {
	conn, err := grpc.NewClient(
		config.GetGRPCAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32*256)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                1 * time.Minute,
			Timeout:             30 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	client := blockchainupdates.NewBlockchainUpdatesApiClient(conn)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if _, err := client.GetBlockUpdate(
				context.Background(),
				&blockchainupdates.GetBlockUpdateRequest{Height: 1},
			); err != nil {
				return err
			}
			l.Info().Msg("Node GRPC probe is OK")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info().Msg("Closing GRPC client")
			return conn.Close()
		},
	})

	return &UpdatesClient{
		conn:             client,
		rest:             &http.Client{},
		nodeURL:          config.GetNodeURL(),
		referencedAssets: make(map[string]ReferencedAsset),
		log:              l,
	}, nil
}
