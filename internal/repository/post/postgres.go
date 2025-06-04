package post

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/post"
	"go.uber.org/zap"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func New(db *pgxpool.Pool, logger *zap.Logger) *PostgresRepo {
	return &PostgresRepo{db: db, logger: logger}
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]post.Post, error) {
	rows, err := r.db.Query(ctx, `SELECT id, topic_id, title, content, username, timestamp FROM backend_schema.posts ORDER BY timestamp DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		if err := rows.Scan(&p.ID, &p.TopicID, &p.Title, &p.Content, &p.Username, &p.Timestamp); err != nil {
			return nil, err
		}
		r.logger.Info("Posts: %+v\n")
		posts = append(posts, p)
	}
	r.logger.Info("All posts return")
	return posts, nil
}

func (r *PostgresRepo) GetByTopic(ctx context.Context, topicID int) ([]post.Post, error) {
	rows, err := r.db.Query(ctx, `SELECT id, topic_id, title, content, username, timestamp FROM backend_schema.posts WHERE topic_id=$1 ORDER BY timestamp DESC`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []post.Post
	for rows.Next() {
		var p post.Post
		if err := rows.Scan(&p.ID, &p.TopicID, &p.Title, &p.Content, &p.Username, &p.Timestamp); err != nil {
			r.logger.Fatal("error", zap.Error(err))
			return nil, err
		}
		posts = append(posts, p)
	}
	r.logger.Info("Posts by Topic return")
	return posts, nil
}

func (r *PostgresRepo) Create(ctx context.Context, p post.Post) error {
	_, err := r.db.Exec(ctx, `INSERT INTO backend_schema.posts (topic_id, title, content, username) VALUES ($1, $2, $3, $4)`,
		p.TopicID, p.Title, p.Content, p.Username)
	return err
}

func (r *PostgresRepo) Delete(ctx context.Context, postID int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM backend_schema.posts WHERE id = $1`, postID)
	return err
}
