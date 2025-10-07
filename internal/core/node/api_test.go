package node

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"

	"github.com/slashfast/wavesconsumer/v2/internal/config"
)

func NewTestNode(ctx context.Context) (*API, error) {
	v := viper.New()
	v.SetConfigType("dotenv")
	v.SetConfigFile("../../../.env")
	globalConfig := config.NewConfig(v)
	cfg := NewConfigFromProvider(globalConfig)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return NewAPI(ctx, &logger, cfg)
}

func TestStartHeightPoller(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, err := NewTestNode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	currentHeight := node.CurrentHeight() - 1

	heightCh := make(chan int32)
	defer close(heightCh)
	go node.StartHeightPoller(heightCh)

	assert.Equal(t, currentHeight, <-heightCh, "Height should be equal")
}

func TestStartHeightProvider(t *testing.T) {
	const (
		delta      = 100
		maxWorkers = 1000
		timeout    = 5 * time.Minute
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	node, err := NewTestNode(ctx)
	if err != nil {
		t.Fatal(err)
	}

	currentHeight := node.CurrentHeight()
	startHeight := currentHeight - delta

	heightCh := make(chan int32)
	defer close(heightCh)

	go node.StartHeightProvider(heightCh, startHeight)

	sem := make(chan struct{}, maxWorkers)
	defer close(sem)

	wg := sync.WaitGroup{}
	defer wg.Wait()

	counter := atomic.Int32{}

	// consume blocks and wait for a new block until the timeout
	for range delta + 1 {
		select {
		case <-ctx.Done():
			t.Fatal("Test timed out")
		case height := <-heightCh:
			sem <- struct{}{}
			counter.Add(1)
			wg.Add(1)

			go func() {
				defer wg.Done()
				defer counter.Add(-1)
				defer func() { <-sem }()

				// delay simulation
				time.Sleep(time.Duration(rand.Intn(300)+100) * time.Millisecond)

				node.log.Debug().Int32("consumed", height).Send()
				node.log.Debug().Int32("workers_count", counter.Load()).Send()
			}()
		}
	}
}
