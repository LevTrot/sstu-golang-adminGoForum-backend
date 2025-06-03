package post

import "time"

type Post struct {
	ID        int       `json:"id"`
	TopicID   int       `json:"topic_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}
