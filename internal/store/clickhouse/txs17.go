package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type Txs17Batch struct {
	*TxsBatch
	AssetID     proto.ColStr
	AssetName   proto.ColStr
	Description proto.ColStr
}

func NewTxs17Batch() *Txs17Batch {
	return &Txs17Batch{
		TxsBatch:    NewTxsBatch(),
		AssetID:     proto.ColStr{},
		AssetName:   proto.ColStr{},
		Description: proto.ColStr{},
	}
}

func (t *Txs17Batch) Append(txs *Txs17) {
	t.TxsBatch.Append(txs.Txs)
	t.AssetID.AppendArr(txs.AssetID)
	t.AssetName.AppendArr(txs.AssetName)
	t.Description.AppendArr(txs.Description)
}
