package domain

import "time"

type ChatMessage struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
