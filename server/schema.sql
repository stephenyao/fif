-- =========================================
-- HOLDINGS TABLE
-- =========================================

CREATE TABLE IF NOT EXISTS holdings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id TEXT NOT NULL,                        -- Firebase UID
    name TEXT NOT NULL,
    symbol VARCHAR(16) NOT NULL,
    quantity NUMERIC(20, 8) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    cost NUMERIC(20, 2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Optional if you want to prevent duplicate tickers per user
    -- UNIQUE (user_id, symbol)
    CHECK (quantity >= 0),
    CHECK (cost >= 0)
);

-- =========================================
-- TRIGGER: update updated_at on row change
-- =========================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_holdings_updated_at
    BEFORE UPDATE ON holdings
    FOR EACH ROW
    EXECUTE PROCEDURE update_updated_at_column();

-- =========================================
-- INDEXES
-- =========================================

-- Fast lookup by user
CREATE INDEX IF NOT EXISTS idx_holdings_user_id
    ON holdings(user_id);

-- Fast "list holdings in order" queries
CREATE INDEX IF NOT EXISTS idx_holdings_user_id_created_at
    ON holdings(user_id, created_at DESC);

-- User UIDs (like Firebase)
-- Replace these with your real Firebase test users if needed
INSERT INTO holdings (user_id, name, symbol, quantity, currency, cost)
VALUES
-- ===== User 1 =====
('v69VFq5fjfhjj4IVckGxL4A1UP92', 'Vanguard Total Stock Market ETF', 'VTI', 12.34567890, 'USD', 2500.00),
('v69VFq5fjfhjj4IVckGxL4A1UP92', 'Apple Inc.', 'AAPL', 20.00000000, 'USD', 3000.00),
('v69VFq5fjfhjj4IVckGxL4A1UP92', 'Tesla Inc.', 'TSLA', 5.00000000, 'USD', 1100.00),