package dto

type QuestionDto struct {
	SenderEmail string `json:"sender_email" binding:"required,email,max=64,min=6"`
	Message     string `json:"password" binding:"required,min=3"`
}

type AnswerDto struct {
	Message string `json:"message" db:"password" binding:"required,min=3"`
}

type MessageDto struct {
	Role    string `json:"role" binding:"required,oneof=user bot"`
	Message string `json:"message" binding:"required"`
}

type ChatDto struct {
	Messages []MessageDto `json:"messages" binding:"required"`
}
