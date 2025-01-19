package utils

import (
	"encoding/base64"
	"github.com/mr-tron/base58"
	"strings"
)

func ExtractAssetId(bytes []byte) string {
	if bytes == nil {
		return "WAVES"
	}
	return base58.Encode(bytes)
}

func ToPrefixedBase64(bytes []byte) string {
	return "base64:" + base64.StdEncoding.EncodeToString(bytes)
}

func EscapeUnicodeNull(s string) string {
	return strings.ReplaceAll(s, "\x00", "\\0")
}
