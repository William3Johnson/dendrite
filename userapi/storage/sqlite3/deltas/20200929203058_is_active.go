package deltas

import (
	"context"
	"database/sql"
	"fmt"
)

func UpIsActive(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	ALTER TABLE account_accounts RENAME TO account_accounts_tmp;
CREATE TABLE account_accounts (
    localpart TEXT NOT NULL PRIMARY KEY,
    created_ts BIGINT NOT NULL,
    password_hash TEXT,
    appservice_id TEXT,
    is_deactivated BOOLEAN DEFAULT 0
);
INSERT
    INTO account_accounts (
      localpart, created_ts, password_hash, appservice_id
    ) SELECT
        localpart, created_ts, password_hash, appservice_id
    FROM account_accounts_tmp
;
DROP TABLE account_accounts_tmp;`)
	if err != nil {
		return fmt.Errorf("failed to execute upgrade: %w", err)
	}
	return nil
}

func DownIsActive(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	ALTER TABLE account_accounts RENAME TO account_accounts_tmp;
CREATE TABLE account_accounts (
    localpart TEXT NOT NULL PRIMARY KEY,
    created_ts BIGINT NOT NULL,
    password_hash TEXT,
    appservice_id TEXT
);
INSERT
    INTO account_accounts (
      localpart, created_ts, password_hash, appservice_id
    ) SELECT
        localpart, created_ts, password_hash, appservice_id
    FROM account_accounts_tmp
;
DROP TABLE account_accounts_tmp;`)
	if err != nil {
		return fmt.Errorf("failed to execute downgrade: %w", err)
	}
	return nil
}
