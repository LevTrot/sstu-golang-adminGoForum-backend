package cleaner

import (
	"context"
	"log"
	"time"

	chatRepo "github.com/LevTrot/sstu-golang-adminGoForum-backend/backend/internal/repository/chat"
)

func StartChatCleaner(repo *chatRepo.Repository) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := repo.DeleteOldMessages(context.Background(), 24*time.Hour)
				if err != nil {
					log.Println("Ошибка очистки сообщений:", err)
				}
			}
		}
	}()
}
