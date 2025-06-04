package comment

import (
	"context"

	models "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/comment"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) *Repository {
	return &Repository{db: db, logger: logger}
}

func (r *Repository) GetByPostID(ctx context.Context, postID int) ([]models.Comment, error) {
	rows, err := r.db.Query(ctx, `SELECT id, post_id, username, content, timestamp FROM backend_schema.comments WHERE post_id=$1 ORDER BY timestamp`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.Username, &c.Content, &c.Timestamp); err != nil {
			r.logger.Fatal("error", zap.Error(err))
			return nil, err
		}
		comments = append(comments, c)
	}
	r.logger.Info("Post by ID return")
	return comments, nil
}

func (r *Repository) Create(ctx context.Context, postID int, username, content string) error {
	_, err := r.db.Exec(ctx, `INSERT INTO backend_schema.comments (post_id, username, content) VALUES ($1, $2, $3)`, postID, username, content)
	return err
}

func (r *Repository) Delete(ctx context.Context, commentID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM backend_schema.comments WHERE id=$1`, commentID)
	return err
}
