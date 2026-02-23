package chat

import "time"

type Message struct {
	ID             int64     `json:"id"`
	SenderUsername string    `json:"sender_username"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateMessageInput struct {
	SenderUsername string `json:"sender_username" validate:"required,min=3,max=50,alphanum"`
	Message        string `json:"message" validate:"required,min=1,max=2000"`
}
