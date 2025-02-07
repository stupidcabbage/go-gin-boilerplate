package chat_bot_service

import (
	"context"

	"example.com/m/internal/api/v1/core/application/exceptions"
)

func (s *ChatBotService) GetAnswer(ctx context.Context, email string, message *ChatMessage) (*ChatMessage, *exceptions.Error_) {
	if message == nil || message.Writer != USER {
		return nil, &exceptions.ErrInvalidQuestion
	}

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
