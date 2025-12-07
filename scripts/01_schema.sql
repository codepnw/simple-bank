-- Table Users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table Auth
CREATE TABLE IF NOT EXISTS auth (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_token ON auth(token);

-- Table Accounts
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGSERIAL NOT NULL REFERENCES users(id),
    balance BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Index
CREATE INDEX idx_accounts_owner_id ON accounts (owner_id);
-- Rule: 1 owner 1 currency
CREATE UNIQUE INDEX idx_accounts_owner_currency ON accounts (owner_id, currency);

-- Table Entries
CREATE TABLE entries (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL, -- (+) Deposit, (-) Withdraw
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Table Transfers
CREATE TABLE transfers (
    id BIGSERIAL PRIMARY KEY,
    from_account_id BIGINT NOT NULL REFERENCES accounts(id),
    to_account_id BIGINT NOT NULL REFERENCES accounts(id),
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Index
CREATE INDEX idx_transfers_from ON transfers (from_account_id);
CREATE INDEX idx_transfers_to ON transfers (to_account_id);
CREATE INDEX idx_transfers_from_to ON transfers (from_account_id, to_account_id);
