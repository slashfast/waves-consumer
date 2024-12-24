package main

import (
	"time"
)

type BalanceUpdate struct {
	height         uint32
	assetId        string
	address        string
	txVersion      uint8
	txType         uint8
	txId           string
	counterparty   string
	subunits       uint64
	deltaSubunits  int64
	assetDecimals  uint8
	assetName      string
	isInbound      bool
	isStateUpdate  bool
	txTimestamp    time.Time
	blockTimestamp time.Time
}

type ReferencedAsset struct {
	id       string
	decimals uint8
	name     string
}

type BlocksHeightResponse struct {
	Height int32 `json:"height"`
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
