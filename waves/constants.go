package waves

var TxTypes = map[string]uint8{
	"Transaction_Genesis":         1,
	"Transaction_Payment":         2,
	"Transaction_Issue":           3,
	"Transaction_Transfer":        4,
	"Transaction_Reissue":         5,
	"Transaction_Exchange":        7,
	"Transaction_Lease":           8,
	"Transaction_LeaseCancel":     9,
	"Transaction_Burn":            6,
	"Transaction_CreateAlias":     10,
	"Transaction_MassTransfer":    11,
	"Transaction_DataTransaction": 12,
	"Transaction_SponsorFee":      14,
	"Transaction_SetScript":       13,
	"Transaction_SetAssetScript":  15,
	"Transaction_InvokeScript":    16,
	"Transaction_UpdateAssetInfo": 17,
	"Transaction_Ethereum":        18,
}

var Asset = ReferencedAsset{
	Id:       "WAVES",
	Decimals: 8,
	Name:     "Waves",
}
