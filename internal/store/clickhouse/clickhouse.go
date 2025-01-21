package clickhouse

import (
	"context"
	"errors"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/internal/store"
	"go.uber.org/fx"
	"net/url"
	"runtime"
	"sync"
)

type ConnectorConfig interface {
	StoreDsn() string
	StoreBufferSize() int32
}

func dsnIntoOpts(dsn string) (*ch.Options, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	password, _ := u.User.Password()

	return &ch.Options{
		Address:  u.Host,
		User:     u.User.Username(),
		Password: password,
		Database: u.Path[1:],
	}, nil
}

type Config struct {
	bufferSize int32
}

type Connector struct {
	ctx    context.Context
	conn   *ch.Client
	log    zerolog.Logger
	wg     *sync.WaitGroup
	config Config
}

func NewConnector(ctx context.Context, lc fx.Lifecycle, l *zerolog.Logger, config ConnectorConfig) (store.Connector, error) {
	opts, err := dsnIntoOpts(config.StoreDsn())
	opts.Compression = ch.CompressionDisabled
	if err != nil {
		return nil, err
	}

	conn, err := ch.Dial(context.Background(), *opts)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			wg.Wait()
			if !conn.IsClosed() {
				return conn.Close()
			}
			return nil
		},
	})

	return &Connector{
		ctx:    ctx,
		conn:   conn,
		log:    l.With().Str("instance", "store").Logger(),
		wg:     wg,
		config: Config{config.StoreBufferSize()},
	}, nil
}

func (c *Connector) Height() (int32, error) {
	c.wg.Add(1)
	defer c.wg.Done()

	var maxHeightResult proto.ColInt32
	err := c.conn.Do(context.Background(), ch.Query{
		Body:   "SELECT max(height) FROM blocks",
		Result: proto.ResultColumn{Data: &maxHeightResult},
	})

	if err != nil {
		return 0, err
	}

	return maxHeightResult[0], nil
}

func (c *Connector) RunLoader(updatesCh store.UpdatesChannel, currentHeightCh chan int32, startHeight int32) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	parsedCh := make(chan *Updates, 1000)
	defer close(parsedCh)

	go c.runParser(parsedCh, currentHeightCh, startHeight)

	sem := make(chan struct{}, runtime.NumCPU())
	defer close(sem)

	worker := func(update store.Update) {
		wg.Add(1)
		defer func() { <-sem }()
		defer wg.Done()
		parsedUpdate, err := c.parseUpdate(update)
		if err != nil {
			c.log.Fatal().Err(err).Msg("failed to parse update")
		}
		parsedCh <- parsedUpdate
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-updatesCh:
			sem <- struct{}{}
			go worker(update)
		}
	}
}

func (c *Connector) runParser(parsedBlocksCh chan *Updates, currentHeightCh chan int32, startHeight int32) {
	targetHeight := <-currentHeightCh
	storeHeight := startHeight
	expectedHeight := storeHeight
	isQuietMode := false
	buffCap := c.config.bufferSize

	parsedBlocksBuffer := make(map[int32]*Updates)
	updatesBuffer := NewUpdatesBatch()

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-parsedBlocksCh:
			parsedBlocksBuffer[update.Blocks.Height[0]] = update

			if update, exists := parsedBlocksBuffer[expectedHeight]; exists {
				updatesBuffer.Append(update)
				delete(parsedBlocksBuffer, expectedHeight)
				expectedHeight++
			}

			heightDelta := targetHeight - storeHeight
			if !isQuietMode {
				isQuietMode = heightDelta <= buffCap
			}

			if isQuietMode {
				buffCap = heightDelta
			}

			buffLen := len(updatesBuffer.Blocks.Height)
			isNeedToFlush := int32(buffLen) == buffCap

			if isNeedToFlush {
				err := c.flushBuffer(updatesBuffer)
				if err != nil {
					c.log.Fatal().Err(err).Msg("failed to flush buffer")
				}
				storeHeight += buffCap
			}
		}
	}
}

func (c *Connector) flushBuffer(buff *UpdatesBatch) error {
	c.wg.Add(1)
	defer c.wg.Done()
	inputs := NewUpdatesInput(buff)

	for table, input := range inputs {
		err := c.conn.Do(context.Background(), ch.Query{
			Body:  "INSERT INTO " + table + " VALUES",
			Input: input,
			Settings: []ch.Setting{
				{Key: "async_insert", Value: "1"},
				{Key: "wait_for_async_insert", Value: "0"},
			},
		})
		if err != nil {
			return errors.Join(err, errors.New(table))
		}
	}

	*buff = *NewUpdatesBatch()
	c.log.Debug().Msg("buffer flushed")

	return nil
}
