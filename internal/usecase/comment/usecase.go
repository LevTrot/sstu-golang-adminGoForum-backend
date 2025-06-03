package comment

import (
	"context"

	models "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/comment"
	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/comment"
)

type Usecase struct {
	repo *comment.Repository
}

func New(repo *comment.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetCommentsByPost(ctx context.Context, postID int) ([]models.Comment, error) {
	return u.repo.GetByPostID(ctx, postID)
}

func (u *Usecase) CreateComment(ctx context.Context, postID int, username, content string) error {
	return u.repo.Create(ctx, postID, username, content)
}

func (u *Usecase) DeleteComment(ctx context.Context, commentID int) error {
	return u.repo.Delete(ctx, commentID)
}
