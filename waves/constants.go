package waves

const (
	Genesis         = 1
	Payment         = 2
	Issue           = 3
	Transfer        = 4
	Reissue         = 5
	Burn            = 6
	Exchange        = 7
	Lease           = 8
	LeaseCancel     = 9
	CreateAlias     = 10
	MassTransfer    = 11
	DataTransaction = 12
	SetScript       = 13
	SponsorFee      = 14
	SetAssetScript  = 15
	InvokeScript    = 16
	UpdateAssetInfo = 17
	Ethereum        = 18
)

var TxTypes = map[string]uint8{
	"Transaction_Genesis":         Genesis,
	"Transaction_Payment":         Payment,
	"Transaction_Issue":           Issue,
	"Transaction_Transfer":        Transfer,
	"Transaction_Reissue":         Reissue,
	"Transaction_Burn":            Burn,
	"Transaction_Exchange":        Exchange,
	"Transaction_Lease":           Lease,
	"Transaction_LeaseCancel":     LeaseCancel,
	"Transaction_CreateAlias":     CreateAlias,
	"Transaction_MassTransfer":    MassTransfer,
	"Transaction_DataTransaction": DataTransaction,
	"Transaction_SponsorFee":      SponsorFee,
	"Transaction_SetScript":       SetScript,
	"Transaction_SetAssetScript":  SetAssetScript,
	"Transaction_InvokeScript":    InvokeScript,
	"Transaction_UpdateAssetInfo": UpdateAssetInfo,
	"Transaction_Ethereum":        Ethereum,
}

var AssetWaves = ReferencedAsset{
	Id:       "WAVES",
	Decimals: 8,
	Name:     "Waves",
}
