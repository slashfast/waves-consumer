package utils

import (
	"encoding/base64"
	"github.com/mr-tron/base58"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"strings"
)

const (
	MainNet = 'W'
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

func MustAddressFromBody(b []byte) string {
	w, err := proto.RebuildAddress(MainNet, b)
	if err != nil {
		panic(err)
	}
	return w.String()
}
