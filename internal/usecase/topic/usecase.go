package topic

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/topic"
)

type Repository interface {
	GetAll(ctx context.Context) ([]topic.Topic, error)
	Create(ctx context.Context, title, description string) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]topic.Topic, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) Create(ctx context.Context, title, description string) error {
	return uc.repo.Create(ctx, title, description)
}

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
