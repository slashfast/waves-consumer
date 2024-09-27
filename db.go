package main

import (
	"context"
	"github.com/ClickHouse/ch-go"
)

func createTables(conn *ch.Client, ctx context.Context, dbName string) error {
	if err := conn.Do(ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbName + `.balances_updates
		(
			height          UInt32,
			asset           String,
			address         String,
			tx_version      UInt8,
			tx_type         UInt8,
			tx              String,
			balance         UInt64,
			asset_decimals  UInt8,
			asset_name      String,
			is_inbound      Bool,
			is_state_update Bool,
			timestamp       DateTime64(3),
			block_timestamp DateTime64(3)
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (height, asset, address, tx)
		`,
	}); err != nil {
		return err
	}

	if err := conn.Do(ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbName + `.blocks
		(
			height          UInt32,
		)
		ENGINE = ReplacingMergeTree
		ORDER BY (height)
		`,
	}); err != nil {
		return err
	}

	return nil
}
