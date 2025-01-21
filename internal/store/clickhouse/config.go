package clickhouse

type ConfigProvider interface {
	StoreDsn() string
	StoreBufferSize() int32
}

type Config struct {
	Dsn        string
	BufferSize int32
}

func NewConfigFromProvider(provider ConfigProvider) Config {
	return Config{
		Dsn:        provider.StoreDsn(),
		BufferSize: provider.StoreBufferSize(),
	}
}
