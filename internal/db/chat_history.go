package db

import (
	"context"
	"database/sql"
	"fmt"
)

const HistoryLimit int = 100

type HistoryRepository struct {
	DB *sql.DB
}

func NewHistoryRepository(db *sql.DB) *HistoryRepository {
	return &HistoryRepository{DB: db}
}

type ChatTurn struct {
	Role string // "user" or "assistant"
	Text string
}

func (r *HistoryRepository) SaveAndLimitChatHistory(ctx context.Context, userID string, role string, text string) (err error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rberr := tx.Rollback(); rberr != nil {
				fmt.Printf("failed to rollback transaction for user %s: %v", userID, rberr)
			}
		}
	}()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO chat_history (user_id, role, text) VALUES (?, ?, ?)",
		userID, role, text,
	)
	if err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	var count int
	err = tx.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM chat_history WHERE user_id = ?",
		userID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count history: %w", err)
	}

	if count > HistoryLimit {
		N := count - HistoryLimit
		_, err = tx.ExecContext(ctx,
			"DELETE FROM chat_history WHERE user_id = ? ORDER BY id ASC LIMIT ?",
			userID, N,
		)
		if err != nil {
			return fmt.Errorf("failed to delete old history: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *HistoryRepository) GetRecentChatHistory(ctx context.Context, userID string) ([]ChatTurn, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT role, text FROM chat_history WHERE user_id = ? ORDER BY id DESC LIMIT ?",
		userID, HistoryLimit-1,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	var history []ChatTurn
	for rows.Next() {
		var turn ChatTurn
		if err := rows.Scan(&turn.Role, &turn.Text); err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		history = append(history, turn)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return history, nil
}
