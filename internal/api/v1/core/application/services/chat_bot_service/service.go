package chat_bot_service

import "context"

const (
	USER writerType = "user"
	BOT  writerType = "bot"
)

type (
	writerType string

	ChatMessage struct {
		Message string
		Writer  writerType
	}

	IBotRepository interface {
		Chat(ctx context.Context, messages ...ChatMessage) ([]ChatMessage, error)
	}

	IStorageRepo interface {
		GetChatByUserEmail(ctx context.Context, email string, offset int, limit int) ([]ChatMessage, error)
		AddNewMessageToChatByEmail(ctx context.Context, email string, message ChatMessage) error
	}

	ChatBotService struct {
		botRepo     IBotRepository
		storageRepo IStorageRepo
	}
)

func NewChatBotService(botRepo IBotRepository, storageRepo IStorageRepo) *ChatBotService {
	return &ChatBotService{
		botRepo:     botRepo,
		storageRepo: storageRepo,
	}
}
