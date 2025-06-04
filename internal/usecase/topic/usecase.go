package topic

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/topic"
	"go.uber.org/zap"
)

type Repository interface {
	GetAll(ctx context.Context) ([]topic.Topic, error)
	Create(ctx context.Context, title, description string) error
	Delete(ctx context.Context, id int64) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (uc *UseCase) GetAll(ctx context.Context) ([]topic.Topic, error) {
	topics, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to fetch topics", zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Topics fetched", zap.Int("count", len(topics)))
	return topics, nil
}

func (uc *UseCase) Create(ctx context.Context, title, description string) error {
	err := uc.repo.Create(ctx, title, description)
	if err != nil {
		uc.logger.Error("Failed to create topic", zap.String("title", title), zap.Error(err))
		return err
	}
	uc.logger.Info("Topic created", zap.String("title", title))
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to delete topic", zap.Int64("topicID", id), zap.Error(err))
		return err
	}
	uc.logger.Info("Topic deleted", zap.Int64("topicID", id))
	return nil
}
