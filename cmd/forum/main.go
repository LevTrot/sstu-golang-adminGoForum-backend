package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"

	postHandler "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/handler/post"
	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/middleware"
	postRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/post"
	postUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/post"

	topicHandler "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/handler/topic"
	topicRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/topic"
	topicUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/topic"

	commentHandler "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/handler/comment"
	commentRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/comment"
	commentUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/comment"

	chatCleaner "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/cleaner"
	chatHandler "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/handler/chat"
	chatRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/chat"
	chatUC "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/usecase/chat"
)

func NewAuthClient(addr string) (authpb.AuthServiceClient, func(), error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		conn.Close()
	}
	client := authpb.NewAuthServiceClient(conn)
	return client, cleanup, nil
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable&search_path=backend_schema"
	}
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	authClient, closeFunc, err := NewAuthClient("localhost:50051")
	if err != nil {
		log.Fatalf("failed to connect to auth service: %v", err)
	}
	defer closeFunc()

	authMiddleware := middleware.AuthMiddleware(authClient)

	r := gin.Default()

	r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	postRepository := postRepo.New(db)
	postUseCase := postUC.New(postRepository)
	postHandler.NewPostHandler(r, postUseCase, authMiddleware)

	topicRepository := topicRepo.New(db)
	topicUseCase := topicUC.New(topicRepository)
	topicHandler.NewTopicHandler(r.Group("/api"), topicUseCase, authMiddleware)

	commentRepository := commentRepo.New(db)
	commentUseCase := commentUC.New(commentRepository)
	commentHandler.NewCommentHandler(r.Group("/api"), commentUseCase, authClient)

	chatRepository := chatRepo.New(db)
	chatUseCase := chatUC.New(chatRepository)
	chatCleaner.StartChatCleaner(chatRepository)
	chatHandler := chatHandler.New(chatUseCase, authClient)

	r.GET("/chat/messages", chatHandler.GetMessagesHandler)
	r.GET("/chat", chatHandler.ChatWebSocketHandler)
	r.GET("/test-token", func(c *gin.Context) {
		token := c.Query("token")
		res, err := authClient.ValidateToken(context.Background(), &authpb.ValidateTokenRequest{
			Token: token,
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": res.Username})
	})

	log.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
