package repositories

import (
	"context"

	"github.com/cohesion-org/deepseek-go"
	"github.com/cohesion-org/deepseek-go/constants"
)

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

	BotRepository struct {
		dpsk           *deepseek.Client
		stopWords      []string
		responseFormat *deepseek.ResponseFormat
	}
)

func NewBotRepository(dpsk *deepseek.Client, responseFormat *deepseek.ResponseFormat, stopWords ...string) *BotRepository {
	return &BotRepository{
		dpsk:           dpsk,
		responseFormat: responseFormat,
		stopWords:      stopWords,
	}
}

func chatMessagesToDeepseekMessages(messages ...ChatMessage) []deepseek.ChatCompletionMessage {
	msgs := make([]deepseek.ChatCompletionMessage, 0, len(messages))
	for _, el := range messages {
		if el.Writer == USER {
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

func responseMessagesToChatMessages(messages ...deepseek.ChatCompletionMessage) []ChatMessage {
	msgs := make([]ChatMessage, 0, len(messages))
	for _, el := range messages {
		if el.Role == constants.ChatMessageRoleUser {
			msgs = append(msgs, ChatMessage{
				Writer:  USER,
				Message: el.Content,
			})
		} else {
			msgs = append(msgs, ChatMessage{
				Writer:  BOT,
				Message: el.Content,
			})
		}
	}

	return msgs
}

func (r *BotRepository) Chat(ctx context.Context, messages ...ChatMessage) ([]ChatMessage, error) {
	msgs := chatMessagesToDeepseekMessages(messages...)

	request := deepseek.ChatCompletionRequest{
		Messages:       msgs,
		Temperature:    1.0,
		Stop:           r.stopWords,
		ResponseFormat: r.responseFormat,
	}
	response, err := r.dpsk.CreateChatCompletion(ctx, &request)
	if err != nil {
		return []ChatMessage{}, err
	}

	responseMessage, err := deepseek.MapMessageToChatCompletionMessage(response.Choices[0].Message)
	if err != nil {
		return []ChatMessage{}, err
	}
	msgs = append(msgs, responseMessage)

	return responseMessagesToChatMessages(msgs...), nil
}

func (r *BotRepository) StreamChat(ctx context.Context, chanellForLastMessage chan ChatMessage, messages []ChatMessage) {
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
			messages = append(messages, ChatMessage{
				Message: fullMessage,
				Writer:  BOT,
			})
			break
		}

		for _, choice := range resp.Choices {
			fullMessage += choice.Delta.Content
			chanellForLastMessage <- ChatMessage{
				Writer:  BOT,
				Message: fullMessage,
			}
		}
	}
}
