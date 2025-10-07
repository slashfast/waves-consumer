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
	TxUID            proto.ColInt64
	TxID             proto.ColStr
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
		TxUID:            proto.ColInt64{},
		TxID:             proto.ColStr{},
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
	d.TxUID.AppendArr(dataStorage.TxUID)
	d.TxID.AppendArr(dataStorage.TxID)
	d.Height.AppendArr(dataStorage.Height)
}
