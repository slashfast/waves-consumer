package clickhouse

import (
	"context"
	"errors"
	"github.com/ClickHouse/ch-go"
	"github.com/ClickHouse/ch-go/proto"
	"github.com/rs/zerolog"
	"github.com/slashfast/wavesconsumer/v2/store"
	"go.uber.org/fx"
	"net/url"
	"sync"
)

type ConnectorConfig interface {
	GetStoreDsn() string
	GetStoreBufferSize() int32
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
	buff   *Updates
}

func NewConnector(ctx context.Context, lc fx.Lifecycle, l *zerolog.Logger, config ConnectorConfig) (store.Connector, error) {
	opts, err := dsnIntoOpts(config.GetStoreDsn())
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
		log:    l.With().Str("instance", "clickhouse").Logger(),
		wg:     wg,
		config: Config{config.GetStoreBufferSize()},
		buff:   NewUpdates(),
	}, nil
}

func (c *Connector) GetHeight() (int32, error) {
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
	targetHeight := <-currentHeightCh
	storeHeight := startHeight
	isRtMode := false
	buffCap := c.config.bufferSize
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for {
		select {
		case <-c.ctx.Done():
			return
		case update := <-updatesCh:
			err := c.parseUpdate(update)
			if err != nil {
				c.log.Fatal().Err(err).Msg("failed to parse update")
			}

			heightDelta := targetHeight - storeHeight
			if !isRtMode {
				isRtMode = heightDelta <= buffCap
			}

			if isRtMode {
				buffCap = heightDelta
			}

			isNeedToFlush := int32(len(c.buff.Blocks.Height)) == buffCap

			if isRtMode || isNeedToFlush {
				err := c.flushBuffer()
				if err != nil {
					c.log.Fatal().Err(err).Msg("failed to flush buffer")
				}
			}
			storeHeight++
		}
	}
}

func (c *Connector) flushBuffer() error {
	c.wg.Add(1)
	defer c.wg.Done()
	inputs := NewUpdatesInput(c.buff)
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
	c.resetBuffer()
	return nil
}

func (c *Connector) resetBuffer() {
	c.buff = NewUpdates()
}
