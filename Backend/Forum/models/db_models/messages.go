package db_models

import "time"

type Message struct {
	GormModel    `json:"-"`
	Id           uint      `json:"id"`
	SenderId     uint      `json:"senderId"`
	RecipientId  uint      `json:"recipientId"`
	Content      string    `json:"content"`
	CreationDate time.Time `json:"creationDate"`
}
