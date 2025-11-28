DROP INDEX IF EXISTS idx_transfers_from_to;
DROP INDEX IF EXISTS idx_transfers_to;
DROP INDEX IF EXISTS idx_transfers_from;
DROP INDEX IF EXISTS idx_accounts_owner_currency;
DROP INDEX IF EXISTS idx_accounts_owner_id;

DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;
