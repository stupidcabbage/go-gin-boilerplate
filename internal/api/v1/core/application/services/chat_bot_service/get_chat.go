package chat_bot_service

import (
	"context"

	"example.com/m/internal/api/v1/core/application/exceptions"
)

func (s *ChatBotService) GetChatByEmail(ctx context.Context, email string, offset int, limit int) ([]ChatMessage, *exceptions.Error_) {
	res, err := s.storageRepo.GetChatByUserEmail(ctx, email, limit, offset)
	if err != nil {
		return []ChatMessage{}, &exceptions.InternalServerError
	}

	return res, nil
}
