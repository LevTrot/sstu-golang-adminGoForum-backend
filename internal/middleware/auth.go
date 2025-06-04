package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/AdminGo/proto/authpb"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(authClient authpb.AuthServiceClient, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		resp, err := authClient.ValidateToken(ctx, &authpb.ValidateTokenRequest{Token: token})
		if err != nil || !resp.Valid {
			logger.Fatal("error:", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": resp.GetError()})
			return
		}

		logger.Info("Middleware success")

		c.Set("user_id", resp.UserId)
		c.Set("username", resp.Username)
		c.Set("role", resp.Role)

		c.Next()
	}
}
