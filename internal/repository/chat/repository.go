package chat

import (
	"context"
	"time"

	domain "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/chat"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) SaveMessage(ctx context.Context, username, content string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO backend_schema.chat_messages (username, content, timestamp) 
		 VALUES ($1, $2, NOW())`,
		username, content,
	)
	r.logger.Info("Saved message")
	return err
}

func (r *Repository) GetRecentMessages(ctx context.Context) ([]domain.ChatMessage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT username, content, timestamp 
		 FROM backend_schema.chat_messages 
		 ORDER BY timestamp ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var msg domain.ChatMessage
		if err := rows.Scan(&msg.Username, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	r.logger.Info("Message return")
	return messages, nil
}

func (r *Repository) DeleteOldMessages(ctx context.Context, olderThan time.Duration) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM backend_schema.chat_messages 
		 WHERE timestamp < NOW() - ($1)::interval`,
		olderThan.String(),
	)
	r.logger.Info("Old Message deleted")
	return err
}
