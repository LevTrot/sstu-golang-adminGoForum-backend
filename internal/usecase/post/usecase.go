package post

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/post"
)

type Repository interface {
	GetAll(ctx context.Context) ([]post.Post, error)
	GetByTopic(ctx context.Context, topicID int) ([]post.Post, error)
	Create(ctx context.Context, p post.Post) error
	Delete(ctx context.Context, postID int) error
}

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]post.Post, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) GetByTopic(ctx context.Context, topicID int) ([]post.Post, error) {
	return uc.repo.GetByTopic(ctx, topicID)
}

func (uc *UseCase) Create(ctx context.Context, p post.Post) error {
	return uc.repo.Create(ctx, p)
}

func (uc *UseCase) Delete(ctx context.Context, postID int) error {
	return uc.repo.Delete(ctx, postID)
}
