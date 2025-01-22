package node

import "time"

type ConfigProvider interface {
	NodeAddress() string
	NodeSecure() bool
	NodeNewHeightDelay() time.Duration
}

type Config struct {
	Address        string
	Secure         bool
	NewHeightDelay time.Duration
}

func NewConfigFromProvider(provider ConfigProvider) Config {
	return Config{
		Address:        provider.NodeAddress(),
		Secure:         provider.NodeSecure(),
		NewHeightDelay: provider.NodeNewHeightDelay(),
	}
}
