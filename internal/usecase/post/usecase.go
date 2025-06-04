package post

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/post"
	"go.uber.org/zap"
)

type Repository interface {
	GetAll(ctx context.Context) ([]post.Post, error)
	GetByTopic(ctx context.Context, topicID int) ([]post.Post, error)
	Create(ctx context.Context, p post.Post) error
	Delete(ctx context.Context, postID int) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]post.Post, error) {
	posts, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all posts", zap.Error(err))
		return nil, err
	}
	uc.logger.Info("All posts fetched", zap.Int("count", len(posts)))
	return posts, nil
}

func (uc *UseCase) GetByTopic(ctx context.Context, topicID int) ([]post.Post, error) {
	posts, err := uc.repo.GetByTopic(ctx, topicID)
	if err != nil {
		uc.logger.Error("Failed to get posts by topic", zap.Int("topicID", topicID), zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Posts fetched by topic", zap.Int("topicID", topicID), zap.Int("count", len(posts)))
	return posts, nil
}

func (uc *UseCase) Create(ctx context.Context, p post.Post) error {
	err := uc.repo.Create(ctx, p)
	if err != nil {
		uc.logger.Error("Failed to create post", zap.String("title", p.Title), zap.String("username", p.Username), zap.Error(err))
		return err
	}
	uc.logger.Info("Post created", zap.String("title", p.Title), zap.String("username", p.Username))
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, postID int) error {
	err := uc.repo.Delete(ctx, postID)
	if err != nil {
		uc.logger.Error("Failed to delete post", zap.Int("postID", postID), zap.Error(err))
		return err
	}
	uc.logger.Info("Post deleted", zap.Int("postID", postID))
	return nil
}
