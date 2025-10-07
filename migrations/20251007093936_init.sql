-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blocks
(
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY height
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_1
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `recipient_address` FixedString(35) CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_2
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `recipient_address` FixedString(35) CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_3
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `asset_name` String CODEC(LZ4),
    `description` String CODEC(ZSTD(3)),
    `amount` Int64 CODEC(T64, LZ4),
    `decimals` Int32 CODEC(T64, LZ4),
    `reissuable` Bool,
    `script` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_4
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `recipient_address` FixedString(35) CODEC(LZ4),
    `recipient_alias` String CODEC(LZ4),
    `fee_asset_id` String CODEC(LZ4),
    `attachment` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_5
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `reissuable` Bool
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_6
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_7
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `price` Int64 CODEC(T64, LZ4),
    `amount_asset_id` String CODEC(LZ4),
    `price_asset_id` String CODEC(LZ4),
    `buy_matcher_fee` Int64 CODEC(T64, LZ4),
    `sell_matcher_fee` Int64 CODEC(T64, LZ4),
    `fee_asset_id` String CODEC(LZ4),
    `orders` Nested(id String, version Int32, sender FixedString(35), sender_public_key String, matcher_public_key String, amount_asset_id String, price_asset_id String, order_type Int32, amount Int64, price Int64, timestamp DateTime64(3), expiration DateTime64(3), matcher_fee Int64, matcher_fee_asset_id String, proofs Array(String), signature String, eip712_signature String, price_mode Int32)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (amount_asset_id, price_asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_8
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `recipient_address` FixedString(35) CODEC(LZ4),
    `recipient_alias` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_9
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `lease_tx_uid` Int64 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_10
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `alias` String CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (alias, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_11
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `attachment` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (uid, timestamp)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_11_transfers
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_tx` Int32 CODEC(T64, LZ4),
    `recipient_address` FixedString(35) CODEC(LZ4),
    `recipient_alias` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_tx)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_12
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (sender, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_12_data
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_tx` Int32 CODEC(T64, LZ4),
    `data_key` String CODEC(LZ4),
    `data_type` LowCardinality(String) CODEC(LZ4),
    `data_value_integer` Int64 CODEC(T64, LZ4),
    `data_value_bool` Bool,
    `data_value_binary` String CODEC(ZSTD(3)),
    `data_value_string` String CODEC(LZ4),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_tx)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_13
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `script` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (sender, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_14
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `min_sponsored_asset_fee` Int64 CODEC(T64, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_15
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `script` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_16
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `dapp_address` FixedString(35) CODEC(LZ4),
    `dapp_alias` String CODEC(LZ4),
    `function_name` String CODEC(LZ4),
    `fee_asset_id` String CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (dapp_address, function_name, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_16_args
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_args` Int32 CODEC(T64, LZ4),
    `arg_type` LowCardinality(String) CODEC(LZ4),
    `arg_value_integer` Int64 CODEC(T64, LZ4),
    `arg_value_bool` Bool,
    `arg_value_binary` String CODEC(ZSTD(3)),
    `arg_value_string` String CODEC(LZ4),
    `arg_value_list` String CODEC(ZSTD(3)),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_args)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_16_payments
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_payment` Int32 CODEC(T64, LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `asset_id` String CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_payment)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_17
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `asset_name` String CODEC(LZ4),
    `description` String CODEC(ZSTD(3))
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (asset_id, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_18
(
    `uid` Int64 CODEC(Delta, LZ4),
    `type` Int8 CODEC(T64, LZ4),
    `sender` FixedString(35) CODEC(LZ4),
    `sender_public_key` String CODEC(LZ4),
    `timestamp` DateTime64(3) CODEC(Delta, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `signature` String CODEC(ZSTD(1)),
    `proofs` Array(String) CODEC(ZSTD(3)),
    `version` Int32 CODEC(T64, LZ4),
    `fee` Int64 CODEC(T64, LZ4),
    `status` LowCardinality(String) CODEC(LZ4),
    `bytes` String CODEC(ZSTD(3)),
    `function_name` String CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (function_name, uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_18_args
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_args` Int32 CODEC(T64, LZ4),
    `arg_type` LowCardinality(String) CODEC(LZ4),
    `arg_value_integer` Int64 CODEC(T64, LZ4),
    `arg_value_bool` Bool,
    `arg_value_binary` String CODEC(ZSTD(3)),
    `arg_value_string` String CODEC(LZ4),
    `arg_value_list` String CODEC(ZSTD(3)),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_args)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS txs_18_payments
(
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `position_in_payment` Int32 CODEC(T64, LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `height` Int32 CODEC(Delta, LZ4),
    `asset_id` String CODEC(LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (tx_uid, position_in_payment)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS balances
(
    `address` FixedString(35) CODEC(LZ4),
    `asset_id` String CODEC(LZ4),
    `amount` Int64 CODEC(T64, LZ4),
    `delta` Int64 CODEC(T64, LZ4),
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (address, asset_id, tx_uid)
PARALLEL WITH
CREATE TABLE IF NOT EXISTS data_storage
(
    `address` FixedString(35) CODEC(LZ4),
    `data_key` String CODEC(LZ4),
    `data_type` LowCardinality(String) CODEC(LZ4),
    `data_value_integer` Int64 CODEC(T64, LZ4),
    `data_value_bool` Bool,
    `data_value_binary` String CODEC(ZSTD(3)),
    `data_value_string` String CODEC(LZ4),
    `tx_uid` Int64 CODEC(Delta, LZ4),
    `tx_id` String CODEC(ZSTD(1)),
    `height` Int32 CODEC(Delta, LZ4)
)
ENGINE = MergeTree
PARTITION BY intDiv(height, 100000) * 100000
ORDER BY (address, data_key, tx_uid)
;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop all tables in reverse order
DROP TABLE IF EXISTS data_storage
PARALLEL WITH
DROP TABLE IF EXISTS balances
PARALLEL WITH
DROP TABLE IF EXISTS txs_18_payments
PARALLEL WITH
DROP TABLE IF EXISTS txs_18_args
PARALLEL WITH
DROP TABLE IF EXISTS txs_18
PARALLEL WITH
DROP TABLE IF EXISTS txs_17
PARALLEL WITH
DROP TABLE IF EXISTS txs_16_payments
PARALLEL WITH
DROP TABLE IF EXISTS txs_16_args
PARALLEL WITH
DROP TABLE IF EXISTS txs_16
PARALLEL WITH
DROP TABLE IF EXISTS txs_15
PARALLEL WITH
DROP TABLE IF EXISTS txs_14
PARALLEL WITH
DROP TABLE IF EXISTS txs_13
PARALLEL WITH
DROP TABLE IF EXISTS txs_12_data
PARALLEL WITH
DROP TABLE IF EXISTS txs_12
PARALLEL WITH
DROP TABLE IF EXISTS txs_11_transfers
PARALLEL WITH
DROP TABLE IF EXISTS txs_11
PARALLEL WITH
DROP TABLE IF EXISTS txs_10
PARALLEL WITH
DROP TABLE IF EXISTS txs_9
PARALLEL WITH
DROP TABLE IF EXISTS txs_8
PARALLEL WITH
DROP TABLE IF EXISTS txs_7
PARALLEL WITH
DROP TABLE IF EXISTS txs_6
PARALLEL WITH
DROP TABLE IF EXISTS txs_5
PARALLEL WITH
DROP TABLE IF EXISTS txs_4
PARALLEL WITH
DROP TABLE IF EXISTS txs_3
PARALLEL WITH
DROP TABLE IF EXISTS txs_2
PARALLEL WITH
DROP TABLE IF EXISTS txs_1
PARALLEL WITH
DROP TABLE IF EXISTS blocks
;
-- +goose StatementEnd
