package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"os"
)

type Config struct {
	NodeURL     string
	GRPCAddress string
	DSN         string
	NatsDSN     string
}

func (c Config) GetNodeURL() string {
	return c.NodeURL
}

func (c Config) GetGRPCAddress() string {
	return c.GRPCAddress
}

func (c Config) GetClickHouseDSN() string {
	return c.DSN
}

func (c Config) GetNatsDSN() string {
	return c.NatsDSN
}

func getEnv(key, fallback string, required bool) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else if required {
		log.Fatal().Msg(key + " is required")
	}
	return fallback
}

func LoadConfigFromEnv(_ fx.Lifecycle, l *zerolog.Logger) (*Config, error) {
	if godotenv.Load() != nil {
		l.Info().Msg("The dotenv file is missing")
	}

	return &Config{
		NodeURL:     getEnv("NODE_URL", "", true),
		GRPCAddress: getEnv("GRPC_ADDRESS", "", true),
		DSN:         getEnv("DSN", "", true),
	}, nil
}
