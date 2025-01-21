package updates

import "time"

type ConfigProvider interface {
	UpdatesAddress() string
	UpdatesSecure() bool
	UpdatesBufferSize() int32
	UpdatesMaxWorkers() int32
	UpdatesAfterSyncDelay() time.Duration
}

type Config struct {
	Address        string
	Secure         bool
	BufferSize     int32
	MaxWorkers     int32
	AfterSyncDelay time.Duration
}

func NewConfigFromProvider(provider ConfigProvider) Config {
	return Config{
		Address:        provider.UpdatesAddress(),
		Secure:         provider.UpdatesSecure(),
		BufferSize:     provider.UpdatesBufferSize(),
		MaxWorkers:     provider.UpdatesMaxWorkers(),
		AfterSyncDelay: provider.UpdatesAfterSyncDelay(),
	}
}
