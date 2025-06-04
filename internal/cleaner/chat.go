package cleaner

import (
	"context"
	"time"

	chatRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/chat"
	"go.uber.org/zap"
)

func StartChatCleaner(repo *chatRepo.Repository, logger *zap.Logger) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := repo.DeleteOldMessages(context.Background(), 24*time.Hour)
				if err != nil {
					logger.Fatal("Failed to clear messages:", zap.Error(err))
				}
			}
		}
	}()
}
