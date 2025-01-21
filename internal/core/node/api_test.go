package node

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestStartHeightPoller(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	v := viper.New()
	v.SetConfigType("dotenv")
	v.SetConfigFile("../../../.env")
	globalConfig := config.NewConfig(v)
	config := NewConfigFromProvider(globalConfig)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	node, err := NewApi(ctx, &logger, config)

	if err != nil {
		t.Fatal(err)
	}

	currentHeight := node.CurrentHeight() - 1

	heightCh := make(chan int32)
	defer close(heightCh)
	go node.StartHeightPoller(heightCh)

	assert.Equal(t, currentHeight, <-heightCh, "Height should be equal")
}
