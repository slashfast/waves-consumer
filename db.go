package main

import (
	"context"
	"github.com/ClickHouse/ch-go"
)

func createTables(conn *ch.Client, ctx context.Context, dbName string) error {
	if err := conn.Do(ctx, ch.Query{
		Body: "CREATE TABLE IF NOT EXISTS " + dbName + `.balances_updates
		(
			asset_id        String,
			address         FixedString(35),
			height          UInt32,
			tx_version      UInt8,
			tx_type         UInt8,
			tx_id           String,
			balance         UInt64,
			delta           Int64,
    		counterparty    String,
			asset_decimals  UInt8,
			asset_name      String,
			is_inbound      Bool,
			is_state_update Bool,
			tx_timestamp    DateTime64(3),
			block_timestamp DateTime64(3)
		)
		ENGINE = ReplacingMergeTree
		PRIMARY KEY (asset_id, address)
		ORDER BY (asset_id, address, height, tx_id)
        SETTINGS non_replicated_deduplication_window = 1000;
		`,
	}); err != nil {
		return err
	}

	//if err := conn.Do(ctx, ch.Query{
	//	Body: "CREATE TABLE IF NOT EXISTS " + dbName + `.blocks
	//	(
	//		height          UInt32,
	//	)
	//	ENGINE = ReplacingMergeTree
	//	ORDER BY (height)
	//	`,
	//}); err != nil {
	//	return err
	//}

	return nil
}
