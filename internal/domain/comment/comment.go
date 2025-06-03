package models

import "time"

type Comment struct {
	ID        int       `db:"id" json:"id"`
	PostID    int       `db:"post_id" json:"postId"`
	Username  string    `db:"username" json:"username"`
	Content   string    `db:"content" json:"content"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}
