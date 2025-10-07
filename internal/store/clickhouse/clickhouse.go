package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"sync"

	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/slashfast/wavesconsumer/v2/internal/store"
)

func dsnIntoOpts(dsn string) (*ch.Options, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	password, _ := u.User.Password()

	return &ch.Options{
		Address:  u.Host,
		User:     u.User.Username(),
		Password: password,
		Database: u.Path[1:],
	}, nil
}

type Connector struct {
	ctx  context.Context
	conn *ch.Client
	log  zerolog.Logger
	wg   *sync.WaitGroup
	cfg  Config
}

func NewConnector(ctx context.Context, lc fx.Lifecycle, l *zerolog.Logger, configProvider ConfigProvider) (store.Connector, error) {
	cfg := NewConfigFromProvider(configProvider)
	opts, err := dsnIntoOpts(cfg.Dsn)
	if err != nil {
		return nil, err
	}
	opts.Compression = ch.CompressionDisabled

	conn, err := ch.Dial(context.Background(), *opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse: %w", err)
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
		ctx:  ctx,
		conn: conn,
		log:  l.With().Str("instance", "store").Logger(),
		wg:   wg,
		cfg:  cfg,
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
		return 0, fmt.Errorf("failed to query max height: %w", err)
	}

	return maxHeightResult[0], nil
}

func (c *Connector) StartLoader(updatesCh store.UpdatesChannel, currentHeightCh chan int32, startHeight int32) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	parsedCh := make(chan *Updates, 1000)
	defer close(parsedCh)

	go c.runParser(parsedCh, currentHeightCh, startHeight)

	sem := make(chan struct{}, runtime.NumCPU())
	defer close(sem)

	worker := func(update store.Update) {
		defer func() { <-sem }()
		defer wg.Done()
		parsedUpdate, err := c.parseUpdate(update)
		if err != nil {
			c.log.Fatal().Err(err).Msg("failed to parse update")
		} else {
			c.log.Trace().Int32("height", parsedUpdate.Blocks.Height[0]).Msg("parsed")
		}
		parsedCh <- parsedUpdate
	}

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-updatesCh:
			sem <- struct{}{}
			wg.Add(1)
			go worker(update)
		}
	}
}

func (c *Connector) runParser(parsedBlocksCh chan *Updates, currentHeightCh chan int32, startHeight int32) {
	targetHeight := <-currentHeightCh
	storeHeight := startHeight
	expectedHeight := storeHeight
	buffCap := c.cfg.BufferSize

	parsedBlocksBuffer := make(map[int32]*Updates)
	updatesBuffer := NewUpdatesBatch()

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-parsedBlocksCh:
			height := update.Blocks.Height[0]
			parsedBlocksBuffer[height] = update
		default:
			if update, exists := parsedBlocksBuffer[expectedHeight]; exists {
				updatesBuffer.Append(update)
				c.log.Trace().Int32("height", expectedHeight).Msg("appended to buffer")
				delete(parsedBlocksBuffer, expectedHeight)
				expectedHeight++
			}

			heightDelta := targetHeight - storeHeight

			if heightDelta <= 0 {
				buffCap = 1
			} else if heightDelta <= buffCap {
				buffCap = heightDelta
			}

			buffLen := len(updatesBuffer.Blocks.Height)

			if buffLen == int(buffCap) {
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
	buffLen := len(buff.Blocks.Height)
	startHeight := buff.Blocks.Height[0]
	endHeight := buff.Blocks.Height[buffLen-1]
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
			return fmt.Errorf("failed to insert into table %s: %w", table, errors.Join(err, errors.New(table)))
		}
	}

	*buff = *NewUpdatesBatch()
	c.log.Info().
		Int32("from", startHeight).
		Int32("to", endHeight).
		Msgf("%d block(s) flushed from buffer", buffLen)

	return nil
}
