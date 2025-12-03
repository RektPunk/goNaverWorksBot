package db

import (
	"context"
	"database/sql"
	"fmt"
)

const ChatHistoryTableDDL = `
CREATE TABLE IF NOT EXISTS chat_history (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     TEXT NOT NULL,
	role		TEXT NOT NULL,
    text        TEXT NOT NULL
);
`

const UserIDIndexDDL = `
CREATE INDEX IF NOT EXISTS idx_user_id ON chat_history (user_id);
`

func InitializeSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, ChatHistoryTableDDL)
	if err != nil {
		return fmt.Errorf("failed to create chat_history table: %w", err)
	}
	_, err = db.ExecContext(ctx, UserIDIndexDDL)
	if err != nil {
		return fmt.Errorf("failed to create idx_user_id index: %w", err)
	}
	return nil
}
