package repositories

import (
	"context"
	"time"

	"example.com/m/internal/api/v1/core/application/services/chat_bot_service"
	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

type BotRepository struct {
	dpsk           *deepseek.Client
	stopWords      []string
	responseFormat *deepseek.ResponseFormat
}

func NewBotRepository(dpsk *deepseek.Client, responseFormat *deepseek.ResponseFormat, stopWords ...string) *BotRepository {
	return &BotRepository{
		dpsk:           dpsk,
		responseFormat: responseFormat,
		stopWords:      stopWords,
	}
}

func chatMessagesToDeepseekMessages(messages ...chat_bot_service.ChatMessage) []deepseek.ChatCompletionMessage {
	msgs := make([]deepseek.ChatCompletionMessage, 0, len(messages))
	for _, el := range messages {
		if el.Writer == chat_bot_service.USER {
			msgs = append(msgs, deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleUser,
				Content: el.Message,
			})
		} else {
			msgs = append(msgs, deepseek.ChatCompletionMessage{
				Role:    constants.ChatMessageRoleSystem,
				Content: el.Message,
			})
		}
	}

	return msgs
}

func responseMessagesToChatMessages(messages ...deepseek.ChatCompletionMessage) []chat_bot_service.ChatMessage {
	msgs := make([]chat_bot_service.ChatMessage, 0, len(messages))
	for _, el := range messages {
		if el.Role == constants.ChatMessageRoleUser {
			msgs = append(msgs, chat_bot_service.ChatMessage{
				Writer:    chat_bot_service.USER,
				Message:   el.Content,
				CreatedAt: time.Now().Format("02.01.2006"),
			})
		} else {
			msgs = append(msgs, chat_bot_service.ChatMessage{
				Writer:    chat_bot_service.BOT,
				Message:   el.Content,
				CreatedAt: time.Now().Format("02.01.2006"),
			})
		}
	}

	return msgs
}

func (r *BotRepository) Chat(ctx context.Context, messages ...chat_bot_service.ChatMessage) ([]chat_bot_service.ChatMessage, error) {
	msgs := chatMessagesToDeepseekMessages(messages...)

	request := deepseek.ChatCompletionRequest{
		Messages:       msgs,
		Temperature:    1.0,
		Stop:           r.stopWords,
		ResponseFormat: r.responseFormat,
	}
	response, err := r.dpsk.CreateChatCompletion(ctx, &request)
	if err != nil {
		return []chat_bot_service.ChatMessage{}, err
	}

	responseMessage, err := deepseek.MapMessageToChatCompletionMessage(response.Choices[0].Message)
	if err != nil {
		return []chat_bot_service.ChatMessage{}, err
	}
	msgs = append(msgs, responseMessage)

	return responseMessagesToChatMessages(msgs...), nil
}

func (r *BotRepository) StreamChat(ctx context.Context, chanellForLastMessage chan chat_bot_service.ChatMessage, messages []chat_bot_service.ChatMessage) {
	msgs := chatMessagesToDeepseekMessages(messages...)

	request := deepseek.StreamChatCompletionRequest{
		Messages:       msgs,
		Temperature:    1.0,
		Stop:           r.stopWords,
		ResponseFormat: r.responseFormat,
		Stream:         true,
	}

	stream, err := r.dpsk.CreateChatCompletionStream(ctx, &request)
	if err != nil {
		close(chanellForLastMessage)
		return
	}
	defer stream.Close()

	var fullMessage string

	for {
		resp, err := stream.Recv()
		if err != nil {
			messages = append(messages, chat_bot_service.ChatMessage{
				Message: fullMessage,
				Writer:  chat_bot_service.BOT,
			})
			break
		}

		for _, choice := range resp.Choices {
			fullMessage += choice.Delta.Content
			chanellForLastMessage <- chat_bot_service.ChatMessage{
				Writer:  chat_bot_service.BOT,
				Message: fullMessage,
			}
		}
	}
	close(chanellForLastMessage)
}
