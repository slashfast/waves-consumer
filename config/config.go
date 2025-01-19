package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	nodeApiAddress  string
	nodeIsApiSecure bool

	UpdatesApiAddress     string
	UpdatesIsApiSecure    bool
	UpdatesBufferSize     int
	UpdatesMaxWorkers     int
	UpdatesAfterSyncDelay time.Duration

	storeDsn        string
	storeBufferSize int32
}

func (c Config) GetNodeApiAddress() string {
	return c.nodeApiAddress
}

func (c Config) IsNodeApiSecure() bool {
	return c.nodeIsApiSecure
}

func (c Config) GetUpdatesApiAddress() string {
	return c.UpdatesApiAddress
}

func (c Config) IsUpdatesApiSecure() bool {
	return c.UpdatesIsApiSecure
}

func (c Config) GetUpdatesBufferSize() int {
	return c.UpdatesBufferSize
}

func (c Config) GetUpdatesMaxWorkers() int {
	return c.UpdatesMaxWorkers
}

func (c Config) GetUpdatesAfterSyncDelay() time.Duration { return c.UpdatesAfterSyncDelay }

func (c Config) GetStoreDsn() string {
	return c.storeDsn
}

func (c Config) GetStoreBufferSize() int32 {
	return c.storeBufferSize
}

type EnvParser struct {
	log *zerolog.Logger
}

func NewEnvParser(logger *zerolog.Logger) *EnvParser {
	if logger == nil {
		logger = &log.Logger
	}
	return &EnvParser{logger}
}

func (c EnvParser) parseBool(val string) bool {
	switch strings.ToLower(val) {
	case "true":
		return true
	case "false":
		return false
	default:
		c.log.Fatal().Msgf("failed to parse a '%s' value as bool", val)
		return false
	}
}

func (c EnvParser) parseInt(val string) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		c.log.Fatal().Err(err).Msgf("failed to parse a '%s' value as int", val)
	}
	return n
}

func (c EnvParser) parseInt32(val string) int32 {
	n, err := strconv.ParseInt(val, 10, 0)
	if err != nil {
		c.log.Fatal().Err(err).Msgf("failed to parse a '%s' value as int32", val)
	}
	return int32(n)
}

func (c EnvParser) parseDuration(val string) time.Duration {
	d, err := time.ParseDuration(val)
	if err != nil {
		c.log.Fatal().Err(err).Msgf("failed to parse a '%s' value as duration", val)
	}
	return d
}

func (c EnvParser) get(key string, fallback any, required bool) any {
	if value, ok := os.LookupEnv(key); ok {
		switch fallback.(type) {
		case string:
			return any(value)
		case bool:
			return c.parseBool(value)
		case int:
			return c.parseInt(value)
		case int32:
			return c.parseInt32(value)
		case time.Duration:
			return c.parseDuration(value)
		}
	} else if required {
		c.log.Fatal().Msg(key + " is required")
	}
	return fallback
}

func LoadConfigFromEnv(l *zerolog.Logger) (*Config, error) {
	logger := l.With().Str("instance", "config_parser").Logger()
	if godotenv.Load() != nil {
		logger.Info().Msg("The dotenv file is missing")
	}

	p := NewEnvParser(&logger)

	return &Config{
		nodeApiAddress:  p.get("NODE_API_ADDRESS", "", true).(string),
		nodeIsApiSecure: p.get("NODE_API_SECURE", true, false).(bool),

		UpdatesApiAddress:     p.get("UPDATES_API_ADDRESS", "", true).(string),
		UpdatesIsApiSecure:    p.get("UPDATES_API_SECURE", true, false).(bool),
		UpdatesBufferSize:     p.get("UPDATES_BUFFER_SIZE", 1000, false).(int),
		UpdatesMaxWorkers:     p.get("UPDATES_MAX_WORKERS", 256*4, false).(int),
		UpdatesAfterSyncDelay: p.get("UPDATES_AFTER_SYNC_DELAY", time.Minute, false).(time.Duration),

		storeDsn:        p.get("STORE_DSN", "", true).(string),
		storeBufferSize: p.get("STORE_BUFFER_SIZE", int32(1000), false).(int32),
	}, nil
}
