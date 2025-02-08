package chat_bot_service

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/core/application/exceptions"
)

const (
	USER writerType = "user"
	BOT  writerType = "bot"
)

type (
	writerType string

	ChatMessage struct {
		Message   string
		Writer    writerType
		CreatedAt string
	}

	IBotRepository interface {
		Chat(ctx context.Context, messages ...ChatMessage) ([]ChatMessage, error)
	}

	IStorageRepo interface {
		GetChatByUserEmail(ctx context.Context, email string, offset int, limit int) ([]ChatMessage, error)
		AddNewMessageToChatByEmail(ctx context.Context, email string, message *ChatMessage) error
	}

	ChatBotService struct {
		botRepo     IBotRepository
		storageRepo IStorageRepo
	}
)

func (s *ChatBotService) GetAnswer(ctx context.Context, email string, message *ChatMessage) (*ChatMessage, *exceptions.Error_) {
	if message == nil || message.Writer != USER {
		return nil, &exceptions.ErrInvalidQuestion
	}
	message.CreatedAt = time.Now().Format("02.01.2006")

	lastMessages, err := s.storageRepo.GetChatByUserEmail(ctx, email, 0, 20)
	if err != nil {
		return nil, &exceptions.InternalServerError
	}
	lastMessages = append(lastMessages, *message)

	answer, err := s.botRepo.Chat(ctx, lastMessages...)
	if err != nil {
		return nil, &exceptions.InternalServerError
	}

	if err := s.storageRepo.AddNewMessageToChatByEmail(ctx, email, answer[len(answer)-1]); err != nil {
		return nil, &exceptions.InternalServerError
	}

	return &answer[len(answer)-1], nil
}

func (s *ChatBotService) GetChatByEmail(ctx context.Context, email string, offset int, limit int) ([]ChatMessage, *exceptions.Error_) {
	res, err := s.storageRepo.GetChatByUserEmail(ctx, email, limit, offset)
	if err != nil {
		return []ChatMessage{}, &exceptions.InternalServerError
	}

	return res, nil
}

func NewChatBotService(botRepo IBotRepository, storageRepo IStorageRepo) *ChatBotService {
	return &ChatBotService{
		botRepo:     botRepo,
		storageRepo: storageRepo,
	}
}
