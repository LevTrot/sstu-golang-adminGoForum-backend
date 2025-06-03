package chat

import (
	"context"
	"time"

	domain "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/chat"
)

type Repository interface {
	SaveMessage(ctx context.Context, username, content string) error
	GetRecentMessages(ctx context.Context) ([]domain.ChatMessage, error)
	DeleteOldMessages(ctx context.Context, olderThan time.Duration) error
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (u *UseCase) SendMessage(ctx context.Context, username, content string) error {
	return u.repo.SaveMessage(ctx, username, content)
}

func (u *UseCase) GetMessages(ctx context.Context) ([]domain.ChatMessage, error) {
	return u.repo.GetRecentMessages(ctx)
}
