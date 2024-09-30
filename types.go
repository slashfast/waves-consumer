package main

import "time"

type BalanceUpdate struct {
	height         uint32
	assetId        string
	address        string
	txVersion      uint8
	txType         uint8
	txId           string
	balance        uint64
	delta          int64
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
