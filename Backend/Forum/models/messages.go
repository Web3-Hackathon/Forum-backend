package models

import "github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"

type NewMessageRequest struct {
	Content string `json:"content"`
}

type MessageResponseModel struct {
	db_models.Message
	SentByUs bool `json:"sentByUs"`
}

type FetchMessagesResponse struct {
	BaseResponseModel
	Messages []MessageResponseModel `json:"messages"`
}
