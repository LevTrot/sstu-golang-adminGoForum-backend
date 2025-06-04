package comment

import (
	"context"

	models "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/comment"
	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/comment"
	"go.uber.org/zap"
)

type Usecase struct {
	repo   *comment.Repository
	logger *zap.Logger
}

func New(repo *comment.Repository, logger *zap.Logger) *Usecase {
	return &Usecase{repo: repo, logger: logger}
}

func (u *Usecase) GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	comments, err := u.repo.GetByPostID(ctx, postID)
	if err != nil {
		u.logger.Error("Failed to fetch comments", zap.Int("postID", postID), zap.Error(err))
		return nil, err
	}
	u.logger.Info("Comments fetched", zap.Int("postID", postID), zap.Int("count", len(comments)))
	return comments, nil
}

func (u *Usecase) CreateComment(ctx context.Context, postID int, username, content string) error {
	err := u.repo.Create(ctx, postID, username, content)
	if err != nil {
		u.logger.Error("Failed to create comment", zap.Int("postID", postID), zap.String("username", username), zap.Error(err))
		return err
	}
	u.logger.Info("Comment created", zap.Int("postID", postID), zap.String("username", username))
	return nil
}

func (u *Usecase) DeleteComment(ctx context.Context, commentID int) error {
	err := u.repo.Delete(ctx, commentID)
	if err != nil {
		u.logger.Error("Failed to delete comment", zap.Int("commentID", commentID), zap.Error(err))
		return err
	}
	u.logger.Info("Comment deleted", zap.Int("commentID", commentID))
	return nil
}
