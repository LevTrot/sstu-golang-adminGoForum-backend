package response

type ErrorResponse struct {
	Error string `json:"error"` // Ошибка
}

type MessageResponse struct {
	Message string `json:"message"` // Сообщение
}

type DataCommentsResponse struct {
	Data []Comment `json:"data"` // Комментарии
}

type DataPostsResponse struct {
	Data []Post `json:"data"` // Посты
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	Content   string `json:"content"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
}

type Post struct {
	ID        int    `json:"id"`
	TopicID   int    `json:"topic_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
}
