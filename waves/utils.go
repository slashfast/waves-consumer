package waves

import (
	"encoding/base64"
	"github.com/mr-tron/base58"
)

func toBase58(bytes []byte) string {
	return base58.Encode(bytes)
}

func extractAssetID(bytes []byte) string {
	if bytes == nil {
		return "WAVES"
	}
	return base58.Encode(bytes)
}
