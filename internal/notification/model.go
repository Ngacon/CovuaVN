package notification

import "time"

type Notification struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateNotificationInput struct {
	Title string `json:"title" validate:"required,min=1,max=200"`
	Body  string `json:"body" validate:"required,min=1,max=2000"`
}
