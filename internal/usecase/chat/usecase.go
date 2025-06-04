package chat

import (
	"context"
	"time"

	domain "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/chat"
	"go.uber.org/zap"
)

type Repository interface {
	SaveMessage(ctx context.Context, username, content string) error
	GetRecentMessages(ctx context.Context) ([]domain.ChatMessage, error)
	DeleteOldMessages(ctx context.Context, olderThan time.Duration) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) SendMessage(ctx context.Context, username, content string) error {
	err := u.repo.SaveMessage(ctx, username, content)
	if err != nil {
		u.logger.Error("Failed to save chat message", zap.String("username", username), zap.Error(err))
		return err
	}
	u.logger.Info("Chat message saved", zap.String("username", username))
	return nil
}

func (u *UseCase) GetMessages(ctx context.Context) ([]domain.ChatMessage, error) {
	msgs, err := u.repo.GetRecentMessages(ctx)
	if err != nil {
		u.logger.Error("Failed to get chat messages", zap.Error(err))
		return nil, err
	}
	u.logger.Info("Fetched recent chat messages", zap.Int("count", len(msgs)))
	return msgs, nil
}
