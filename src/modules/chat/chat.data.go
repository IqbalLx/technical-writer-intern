package chat

import (
	"context"

	"github.com/IqbalLx/technical-writer-intern/src/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/leporo/sqlf"
)

func insertNewChatHistories(ctx context.Context, db *pgxpool.Pool, userName string, chatHistories []entities.ChatHistory) error {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		InsertInto("chat_histories")

	for _, chat := range chatHistories {
		query.
			NewRow().
			Set("user_name", userName).
			Set("role", chat.Role).
			Set("content", chat.Content)
	}

	sql, args := query.String(), query.Args()
	if _, err := db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

func countUserChatHistories(ctx context.Context, db *pgxpool.Pool, userName string) (int, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		Select("COALESCE(COUNT(*), 0)").
		From("chat_histories").
		Where("user_name = ?", userName).
		Where("timestamp BETWEEN NOW() - INTERVAL '7 day' AND NOW()")

	sql, args := query.String(), query.Args()
	row := db.QueryRow(ctx, sql, args...)

	var count int
	err := row.Scan(&count)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return 0, nil
		}

		return count, err
	}

	return count, nil
}

func getUserChatHistories(ctx context.Context, db *pgxpool.Pool, userName string, offset int, limit int) ([]entities.ChatHistory, error) {
	sqlf.SetDialect(sqlf.PostgreSQL)
	query := sqlf.
		Select("role, content").
		From("chat_histories").
		Where("user_name = ?", userName).
		Where("timestamp BETWEEN NOW() - INTERVAL '7 day' AND NOW()").
		OrderBy("timestamp ASC").
		Offset(offset).
		Limit(limit)

	sql, args := query.String(), query.Args()
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return []entities.ChatHistory{}, err
	}
	chats, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (entities.ChatHistory, error) {
		var chat entities.ChatHistory
		err = row.Scan(
			&chat.Role,
			&chat.Content,
		)
		if err != nil {
			return chat, err
		}

		return chat, err
	})

	if err != nil {
		return []entities.ChatHistory{}, err
	}

	return chats, nil
}
