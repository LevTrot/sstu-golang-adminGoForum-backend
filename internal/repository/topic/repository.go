package repository

import (
	"context"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/domain/topic"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TopicRepository struct {
	DB *pgxpool.Pool
}

func New(db *pgxpool.Pool) *TopicRepository {
	return &TopicRepository{DB: db}
}

func (r *TopicRepository) GetAll(ctx context.Context) ([]topic.Topic, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, title, description FROM backend_schema.topics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []topic.Topic
	for rows.Next() {
		var t topic.Topic
		if err := rows.Scan(&t.ID, &t.Title, &t.Description); err != nil {
			return nil, err
		}
		topics = append(topics, t)
	}
	return topics, nil
}

func (r *TopicRepository) Create(ctx context.Context, title, description string) error {
	_, err := r.DB.Exec(ctx,
		"INSERT INTO backend_schema.topics (title, description) VALUES ($1, $2)",
		title, description)
	return err
}

func (r *TopicRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.Exec(ctx,
		"DELETE FROM backend_schema.topics WHERE id = $1", id)
	return err
}
