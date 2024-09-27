package main

import "time"

type BalanceUpdate struct {
	height         uint32
	asset          string
	address        string
	txVersion      uint8
	txType         uint8
	tx             string
	balance        uint64
	assetDecimals  uint8
	assetName      string
	isInbound      bool
	isStateUpdate  bool
	timestamp      time.Time
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
