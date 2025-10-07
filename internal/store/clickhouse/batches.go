package clickhouse

type UpdatesBatch struct {
	Blocks         *BlocksBatch
	Txs1           *Txs1Batch
	Txs2           *Txs2Batch
	Txs3           *Txs3Batch
	Txs4           *Txs4Batch
	Txs5           *Txs5Batch
	Txs6           *Txs6Batch
	Txs7           *Txs7Batch
	Txs8           *Txs8Batch
	Txs9           *Txs9Batch
	Txs10          *Txs10Batch
	Txs11          *Txs11Batch
	Txs11Transfers *Txs11TransferBatch
	Txs12          *Txs12Batch
	Txs12Data      *Txs12DataBatch
	Txs13          *Txs13Batch
	Txs14          *Txs14Batch
	Txs15          *Txs15Batch
	Txs16          *Txs16Batch
	Txs16Args      *TxsArgsBatch
	Txs16Payments  *TxsPaymentsBatch
	Txs17          *Txs17Batch
	Txs18          *Txs18Batch
	Txs18Args      *TxsArgsBatch
	Txs18Payments  *TxsPaymentsBatch
	Balances       *BalanceBatch
	DataStorage    *DataStorageBatch
}

func NewUpdatesBatch() *UpdatesBatch {
	return &UpdatesBatch{
		Blocks:         NewBlocksBatch(),
		Txs1:           NewTxs1Batch(),
		Txs2:           NewTxs2Batch(),
		Txs3:           NewTxs3Batch(),
		Txs4:           NewTxs4Batch(),
		Txs5:           NewTxs5Batch(),
		Txs6:           NewTxs6Batch(),
		Txs7:           NewTxs7Batch(),
		Txs8:           NewTxs8Batch(),
		Txs9:           NewTxs9Batch(),
		Txs10:          NewTxs10Batch(),
		Txs11:          NewTxs11Batch(),
		Txs11Transfers: NewTxs11TransferBatch(),
		Txs12:          NewTxs12Batch(),
		Txs12Data:      NewTxs12DataBatch(),
		Txs13:          NewTxs13Batch(),
		Txs14:          NewTxs14Batch(),
		Txs15:          NewTxs15Batch(),
		Txs16:          NewTxs16Batch(),
		Txs16Args:      NewTxsArgsBatch(),
		Txs16Payments:  NewTxsPaymentsBatch(),
		Txs17:          NewTxs17Batch(),
		Txs18:          NewTxs18Batch(),
		Txs18Args:      NewTxsArgsBatch(),
		Txs18Payments:  NewTxsPaymentsBatch(),
		Balances:       NewBalanceBatch(),
		DataStorage:    NewDataStorageBatch(),
	}
}

func (u *UpdatesBatch) Append(updates *Updates) {
	u.Blocks.Height.AppendArr(updates.Blocks.Height)
	u.Blocks.ID.AppendArr(updates.Blocks.ID)
	u.Blocks.Timestamp.AppendArr(updates.Blocks.Timestamp)

	u.Txs1.Append(updates.Txs1)
	u.Txs2.Append(updates.Txs2)
	u.Txs3.Append(updates.Txs3)
	u.Txs4.Append(updates.Txs4)
	u.Txs5.Append(updates.Txs5)
	u.Txs6.Append(updates.Txs6)
	u.Txs7.Append(updates.Txs7)
	u.Txs8.Append(updates.Txs8)
	u.Txs9.Append(updates.Txs9)
	u.Txs10.Append(updates.Txs10)
	u.Txs11.Append(updates.Txs11)
	u.Txs11Transfers.Append(updates.Txs11Transfers)
	u.Txs12.Append(updates.Txs12)
	u.Txs12Data.Append(updates.Txs12Data)
	u.Txs13.Append(updates.Txs13)
	u.Txs14.Append(updates.Txs14)
	u.Txs15.Append(updates.Txs15)
	u.Txs16.Append(updates.Txs16)
	u.Txs16Args.Append(updates.Txs16Args)
	u.Txs16Payments.Append(updates.Txs16Payments)
	u.Txs17.Append(updates.Txs17)
	u.Txs18.Append(updates.Txs18)
	u.Txs18Args.Append(updates.Txs18Args)
	u.Txs18Payments.Append(updates.Txs18Payments)
	u.Balances.Append(updates.Balances)
	u.DataStorage.Append(updates.DataStorage)
}
