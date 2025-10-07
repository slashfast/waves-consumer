package clickhouse

import (
	"github.com/ClickHouse/ch-go/proto"
)

type DataStorageBatch struct {
	Address          proto.ColStr
	DataKey          proto.ColStr
	DataType         proto.ColStr
	DataValueInteger proto.ColInt64
	DataValueBoolean proto.ColBool
	DataValueBinary  proto.ColStr
	DataValueString  proto.ColStr
	TxUid            proto.ColInt64
	TxId             proto.ColStr
	Height           proto.ColInt32
}

func NewDataStorageBatch() *DataStorageBatch {
	return &DataStorageBatch{
		Address:          proto.ColStr{},
		DataKey:          proto.ColStr{},
		DataType:         proto.ColStr{},
		DataValueInteger: proto.ColInt64{},
		DataValueBoolean: proto.ColBool{},
		DataValueBinary:  proto.ColStr{},
		DataValueString:  proto.ColStr{},
		TxUid:            proto.ColInt64{},
		TxId:             proto.ColStr{},
		Height:           proto.ColInt32{},
	}
}

func (d *DataStorageBatch) Append(dataStorage *DataStorage) {
	d.Address.AppendArr(dataStorage.Address)
	d.DataKey.AppendArr(dataStorage.DataKey)
	d.DataType.AppendArr(dataStorage.DataType)
	d.DataValueInteger.AppendArr(dataStorage.DataValueInteger)
	d.DataValueBoolean.AppendArr(dataStorage.DataValueBoolean)
	d.DataValueBinary.AppendArr(dataStorage.DataValueBinary)
	d.DataValueString.AppendArr(dataStorage.DataValueString)
	d.TxUid.AppendArr(dataStorage.TxUid)
	d.TxId.AppendArr(dataStorage.TxId)
	d.Height.AppendArr(dataStorage.Height)
}
