package config

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"time"
)

const (
	nodeAddress = "NODE_ADDRESS"
	nodeSecure  = "NODE_SECURE"

	updatesAddress        = "UPDATES_ADDRESS"
	updatesSecure         = "UPDATES_SECURE"
	updatesBufferSize     = "UPDATES_BUFFER_SIZE"
	updatesMaxWorkers     = "UPDATES_MAX_WORKERS"
	updatesAfterSyncDelay = "UPDATES_AFTER_SYNC_DELAY"

	storeDsn        = "STORE_DSN"
	storeBufferSize = "STORE_BUFFER_SIZE"
)

type Config struct {
	nodeAddress string
	nodeSecure  bool

	updatesAddress        string
	updatesSecure         bool
	updatesBufferSize     int32
	updatesMaxWorkers     int32
	updatesAfterSyncDelay time.Duration

	storeDsn        string
	storeBufferSize int32
}

func (c Config) NodeAddress() string {
	return c.nodeAddress
}

func (c Config) NodeSecure() bool {
	return c.nodeSecure
}

func (c Config) UpdatesAddress() string {
	return c.updatesAddress
}

func (c Config) UpdatesSecure() bool {
	return c.updatesSecure
}

func (c Config) UpdatesBufferSize() int32 {
	return c.updatesBufferSize
}

func (c Config) UpdatesMaxWorkers() int32 {
	return c.updatesMaxWorkers
}

func (c Config) UpdatesAfterSyncDelay() time.Duration {
	return c.updatesAfterSyncDelay
}

func (c Config) StoreDsn() string {
	return c.storeDsn
}

func (c Config) StoreBufferSize() int32 {
	return c.storeBufferSize
}

func NewConfig() Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigType("dotenv")
	v.SetConfigFile(".env")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	v.SetDefault(nodeSecure, true)
	v.SetDefault(updatesSecure, true)
	v.SetDefault(updatesBufferSize, 1000)
	v.SetDefault(updatesMaxWorkers, 1024)
	v.SetDefault(updatesAfterSyncDelay, "1m")
	v.SetDefault(storeBufferSize, 1000)

	return Config{
		nodeAddress:           v.GetString(nodeAddress),
		nodeSecure:            v.GetBool(nodeSecure),
		updatesAddress:        v.GetString(updatesAddress),
		updatesSecure:         v.GetBool(updatesSecure),
		updatesBufferSize:     v.GetInt32(updatesBufferSize),
		updatesMaxWorkers:     v.GetInt32(updatesMaxWorkers),
		updatesAfterSyncDelay: v.GetDuration(updatesAfterSyncDelay),
		storeDsn:              v.GetString(storeDsn),
		storeBufferSize:       v.GetInt32(storeBufferSize),
	}
}

func Provider(l *zerolog.Logger) (Config, error) {
	logger := l.With().Str("instance", "config_parser").Logger()

	cfg := NewConfig()

	if cfg.nodeAddress == "" {
		logger.Fatal().Msgf("%s is required", nodeAddress)
	}
	if cfg.updatesAddress == "" {
		logger.Fatal().Msgf("%s is required", updatesAddress)
	}
	if cfg.storeDsn == "" {
		logger.Fatal().Msgf("%s is required", storeDsn)
	}

	return cfg, nil
}
