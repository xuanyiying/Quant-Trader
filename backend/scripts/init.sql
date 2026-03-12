-- Quant Trader Database Schema
-- Single initialization file - run this to create all tables

-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- =====================
-- Market Data Tables
-- =====================

CREATE TABLE IF NOT EXISTS trades (
    trade_id TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    exchange TEXT NOT NULL,
    price NUMERIC NOT NULL,
    amount NUMERIC NOT NULL,
    side TEXT,
    PRIMARY KEY (symbol, exchange, trade_id, time)
);

DO $$
BEGIN
    PERFORM create_hypertable('trades', 'time');
EXCEPTION
    WHEN others THEN RAISE NOTICE 'Table trades is already a hypertable: %', SQLERRM;
END $$;

CREATE INDEX IF NOT EXISTS idx_trades_symbol_time ON trades (symbol, time DESC);

CREATE TABLE IF NOT EXISTS klines (
    id BIGSERIAL,
    time TIMESTAMPTZ NOT NULL,
    symbol TEXT NOT NULL,
    exchange TEXT NOT NULL,
    period TEXT NOT NULL,
    open NUMERIC NOT NULL,
    high NUMERIC NOT NULL,
    low NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    volume NUMERIC NOT NULL,
    PRIMARY KEY (symbol, exchange, period, time)
);

DO $$
BEGIN
    PERFORM create_hypertable('klines', 'time');
EXCEPTION
    WHEN others THEN RAISE NOTICE 'Table klines is already a hypertable: %', SQLERRM;
END $$;

CREATE INDEX IF NOT EXISTS idx_klines_symbol_period_time ON klines (symbol, period, time DESC);

-- =====================
-- User Management
-- =====================

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(255),
    tier_id BIGINT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_exchange_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    exchange VARCHAR(50) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_secret VARCHAR(255) NOT NULL,
    label VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    key_id VARCHAR(50) UNIQUE NOT NULL,
    key_secret VARCHAR(100) NOT NULL,
    name VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user ON api_keys(user_id);

-- =====================
-- Strategies
-- =====================

CREATE TABLE IF NOT EXISTS strategies (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS backtest_runs (
    id BIGSERIAL PRIMARY KEY,
    strategy_id BIGINT REFERENCES strategies(id),
    symbol VARCHAR(50) NOT NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    report JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- =====================
-- Subscription System
-- =====================

CREATE TABLE IF NOT EXISTS subscription_tiers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    max_symbols INT NOT NULL,
    realtime_enabled BOOLEAN DEFAULT FALSE,
    price_monthly NUMERIC(10, 2),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE,
    tier_id INT REFERENCES subscription_tiers(id),
    status VARCHAR(20) DEFAULT 'active',
    expires_at TIMESTAMPTZ,
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_user_subs_stripe_customer ON user_subscriptions(stripe_customer_id);
CREATE INDEX IF NOT EXISTS idx_user_subs_stripe_sub ON user_subscriptions(stripe_subscription_id);

INSERT INTO subscription_tiers (name, max_symbols, realtime_enabled, price_monthly) 
VALUES ('Free', 1, FALSE, 0.00) ON CONFLICT DO NOTHING;
INSERT INTO subscription_tiers (name, max_symbols, realtime_enabled, price_monthly) 
VALUES ('Pro', 10, TRUE, 29.00) ON CONFLICT DO NOTHING;
INSERT INTO subscription_tiers (name, max_symbols, realtime_enabled, price_monthly) 
VALUES ('Enterprise', 1000, TRUE, 99.00) ON CONFLICT DO NOTHING;

-- =====================
-- Alert System
-- =====================

CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    symbol VARCHAR(20) NOT NULL,
    condition_type VARCHAR(50) NOT NULL,
    target_value NUMERIC(20, 8),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_symbol_active ON alerts(symbol) WHERE is_active = TRUE;

-- =====================
-- Paper Trading
-- =====================

CREATE TABLE IF NOT EXISTS paper_accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id),
    balance NUMERIC(20, 8) DEFAULT 100000.0,
    initial_balance NUMERIC(20, 8) DEFAULT 100000.0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS paper_orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    price NUMERIC(20, 8),
    qty NUMERIC(20, 8) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',
    filled_price NUMERIC(20, 8),
    filled_time TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS paper_positions (
    user_id BIGINT REFERENCES users(id),
    symbol VARCHAR(20) NOT NULL,
    qty NUMERIC(20, 8) DEFAULT 0,
    avg_price NUMERIC(20, 8) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, symbol)
);

CREATE INDEX IF NOT EXISTS idx_paper_orders_user ON paper_orders(user_id);

-- =====================
-- Portfolio Management
-- =====================

CREATE TABLE IF NOT EXISTS portfolios (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS portfolio_assets (
    id BIGSERIAL PRIMARY KEY,
    portfolio_id BIGINT NOT NULL REFERENCES portfolios(id),
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL DEFAULT 'long',
    qty NUMERIC(20, 8) DEFAULT 0,
    entry_price NUMERIC(20, 8) DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(portfolio_id, symbol)
);

CREATE INDEX IF NOT EXISTS idx_portfolios_user ON portfolios(user_id);

-- =====================
-- Strategy Marketplace
-- =====================

CREATE TABLE IF NOT EXISTS strategy_market (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT REFERENCES users(id),
    strategy_id BIGINT REFERENCES strategies(id),
    price NUMERIC(10, 2) DEFAULT 0,
    description TEXT,
    performance_metrics JSONB,
    is_public BOOLEAN DEFAULT FALSE,
    subscriber_count INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    UNIQUE(strategy_id)
);

CREATE TABLE IF NOT EXISTS strategy_purchases (
    user_id BIGINT REFERENCES users(id),
    market_item_id BIGINT REFERENCES strategy_market(id),
    purchased_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, market_item_id)
);

CREATE INDEX IF NOT EXISTS idx_market_public ON strategy_market(is_public) WHERE is_public = TRUE;
