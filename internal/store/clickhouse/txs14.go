package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs14Batch struct {
	*TxsBatch
	AssetID              proto.ColStr
	MinSponsoredAssetFee proto.ColInt64
}

func NewTxs14Batch() *Txs14Batch {
	return &Txs14Batch{
		TxsBatch:             NewTxsBatch(),
		AssetID:              proto.ColStr{},
		MinSponsoredAssetFee: proto.ColInt64{},
	}
}

func (t *Txs14Batch) Append(txs *Txs14) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.MinSponsoredAssetFee.AppendArr(txs.MinSponsoredAssetFee)
}
