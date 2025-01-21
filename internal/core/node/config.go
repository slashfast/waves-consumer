package node

type ConfigProvider interface {
	NodeAddress() string
	NodeSecure() bool
}

type Config struct {
	Address string
	Secure  bool
}

func NewConfigFromProvider(provider ConfigProvider) Config {
	return Config{
		Address: provider.NodeAddress(),
		Secure:  provider.NodeSecure(),
	}
}
