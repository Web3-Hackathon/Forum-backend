package controllers

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/database"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/middleware"
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gorm.io/gorm"
	"time"
)

func SendMessage(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var username, ok = utils.ParseParam("username", ctx, p)
	if !ok {
		return
	}

	var recipientId uint
	recipientId, ok = database.GetUserId(username)
	if !ok {
		utils.Error(ctx, 3, "User does not exist")
		return
	}

	if session.UserId == recipientId {
		utils.Error(ctx, 9, "Cannot send message to yourself")
		return
	}

	var messageData = new(models.NewMessageRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &messageData) == nil {
		return
	}

	if len(messageData.Content) < 5 {
		utils.Error(ctx, 5, "Message content should be at least 5 characters long.")
		return
	}

	var newMessage = db_models.Message{
		SenderId:     session.UserId,
		RecipientId:  recipientId,
		Content:      messageData.Content,
		CreationDate: time.Now(),
	}

	var response = database.MySQLClient.Create(&newMessage)
	if response.Error != nil {
		logger.Logf(logger.ERROR, "Could not create message from %d to %d. Error: %s",
			session.UserId, recipientId, response.Error)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully sent message",
	})
}

func FetchMessages(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var username, ok = utils.ParseParam("username", ctx, p)
	if !ok {
		return
	}

	var userId uint
	userId, ok = database.GetUserId(username)
	if !ok {
		utils.Error(ctx, 3, "User does not exist")
		return
	}

	if session.UserId == userId {
		utils.Error(ctx, 9, "Cannot fetch chat with yourself")
		return
	}

	var messageModels []db_models.Message
	var err = database.MySQLClient.Where("").Order("creation_date ASC").Find(&messageModels).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch chats from database between %d and %d. Error: %s",
			session.UserId, userId, err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var messages []models.MessageResponseModel

	for _, msg := range messageModels {
		messages = append(messages, models.MessageResponseModel{
			Message:  msg,
			SentByUs: msg.SenderId == session.UserId,
		})
	}

	utils.JSON(ctx, models.FetchMessagesResponse{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched messages",
		},
		Messages: messages,
	})
}
