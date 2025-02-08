package repositories

import (
	"context"
	"database/sql"

	"example.com/m/internal/api/v1/core/application/dto"
	"example.com/m/internal/api/v1/core/application/services/chat_bot_service"
	"github.com/doug-martin/goqu/v9"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepo {
	return &ChatRepo{
		db: db,
	}
}

func dbMessageToChatMessage(messages ...dto.DbMessageDto) []chat_bot_service.ChatMessage {
	msgs := make([]chat_bot_service.ChatMessage, 0, len(messages))
	for _, el := range messages {
		if el.Role == "user" {
			msgs = append(msgs, chat_bot_service.ChatMessage{
				Message:   el.Text,
				Writer:    chat_bot_service.USER,
				CreatedAt: el.CreatedAt,
			})
		} else {
			msgs = append(msgs, chat_bot_service.ChatMessage{
				Message:   el.Text,
				Writer:    chat_bot_service.BOT,
				CreatedAt: el.CreatedAt,
			})
		}
	}

	return msgs
}

func chatMessageToDto(message *chat_bot_service.ChatMessage, email string) *dto.DbMessageDto {
	if message.Writer == chat_bot_service.USER {
		return &dto.DbMessageDto{
			Email:     email,
			Text:      message.Message,
			Role:      "user",
			CreatedAt: message.CreatedAt,
		}
	} else {
		return &dto.DbMessageDto{
			Email:     email,
			Text:      message.Message,
			Role:      "bot",
			CreatedAt: message.CreatedAt,
		}
	}
}

func (r *ChatRepo) GetChatByUserEmail(ctx context.Context, email string, offset int, limit int) ([]chat_bot_service.ChatMessage, error) {
	query, _, _ := goqu.From("chat_messages").
		Where(goqu.Ex{
			"email": email,
		}).
		Order(goqu.C("created_at").Desc()).
		Offset(uint(offset)).
		Limit(uint(limit)).
		ToSQL()

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []dto.DbMessageDto
	for rows.Next() {
		var msg dto.DbMessageDto
		if err := rows.Scan(&msg.Email, &msg.Text, &msg.Role, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dbMessageToChatMessage(messages...), nil
}

func (r *ChatRepo) AddNewMessageToChatByEmail(ctx context.Context, email string, message *chat_bot_service.ChatMessage) error {
	msgDto := chatMessageToDto(message, email)
	query, _, _ := goqu.Insert("chat_messages").Rows(*msgDto).ToSQL()
	_, err := r.db.Exec(query)

	if err != nil {
		return err
	}

	return nil
}
