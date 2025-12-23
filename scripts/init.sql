-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- 1. Trades Table
CREATE TABLE IF NOT EXISTS trades (
    trade_id TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    exchange TEXT NOT NULL,
    price NUMERIC NOT NULL,
    amount NUMERIC NOT NULL,
    side TEXT, -- 'buy' or 'sell'
    PRIMARY KEY (symbol, exchange, trade_id, time)
);

-- Convert to hypertable (with exception handling for "already exists")
DO $$
BEGIN
    BEGIN
        PERFORM create_hypertable('trades', 'time');
    EXCEPTION
        WHEN others THEN
            RAISE NOTICE 'Table trades is already a hypertable or failed to convert: %', SQLERRM;
    END;
END $$;

CREATE INDEX IF NOT EXISTS idx_trades_symbol_time ON trades (symbol, time DESC);

-- 2. K-Line Table
CREATE TABLE IF NOT EXISTS klines (
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    exchange TEXT NOT NULL,
    period TEXT NOT NULL, -- '1m', '5m', '1h'
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (symbol, exchange, period, time)
);

-- Convert to hypertable (with exception handling for "already exists")
DO $$
BEGIN
    BEGIN
        PERFORM create_hypertable('klines', 'time');
    EXCEPTION
        WHEN others THEN
            RAISE NOTICE 'Table klines is already a hypertable or failed to convert: %', SQLERRM;
    END;
END $$;

CREATE INDEX IF NOT EXISTS idx_klines_symbol_period_time ON klines (symbol, period, time DESC);

-- 3. Business Tables
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW (),
    updated_at TIMESTAMPTZ DEFAULT NOW ()
);

CREATE TABLE IF NOT EXISTS user_exchange_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    exchange VARCHAR(50) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_secret VARCHAR(255) NOT NULL,
    label VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW ()
);

CREATE TABLE IF NOT EXISTS strategies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users (id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW ()
);

CREATE TABLE IF NOT EXISTS backtest_runs (
    id BIGSERIAL PRIMARY KEY,
    strategy_id BIGINT REFERENCES strategies (id),
    symbol VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    report JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW ()
);
