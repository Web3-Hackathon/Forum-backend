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

func GetUsersReputation(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var err error
	var username string
	var ok bool

	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	// TODO: Add check to see if username exists
	username, ok = utils.ParseParam("username", ctx, p)
	if !ok {
		return
	}

	var recipientId uint
	recipientId, ok = database.GetUserId(username)
	if !ok {
		utils.Error(ctx, 3, "User does not exist")
		return
	}

	var reputation []db_models.Reputation

	err = database.MySQLClient.Order("creation_date DESC").Where("recipient_id = ?", recipientId).First(&reputation).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch reputation from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.ReputationResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched reputation",
		},
		Reputation: reputation,
	})
}

func ModifyUsersReputation(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var username string
	var ok bool
	var err error

	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	// TODO: Make a check to make sure that the user exists
	username, ok = utils.ParseParam("username", ctx, p)
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
		utils.Error(ctx, 9, "Cannot vouch/modify your own reputation")
		return
	}

	var repData = new(models.UserReputationRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, repData) == nil {
		return
	}

	if repData.Message == "" {
		utils.Error(ctx, 5, "Invalid request body")
		return
	}

	// TODO: Check if the modifier sent by the client is valid, otherwise ban the user instantly for abusing API
	var reputation db_models.Reputation
	var exists = true

	err = database.MySQLClient.Where("sender_id = ? AND recipient_id = ?",
		session.UserId, recipientId).First(&reputation).Error
	if err == gorm.ErrRecordNotFound {
		exists = false
	}

	// TODO: Fix this. Issue: When user resubmits the reputation can double. Make check for the current rep and calculate new one

	if exists {
		err = database.MySQLClient.Model(&reputation).Where("sender_id = ? AND recipient_id = ?").
			Updates(db_models.Reputation{
				Message:      repData.Message,
				Modifier:     repData.Modifier,
				CreationDate: time.Now(),
			}).Error

	} else {
		err = database.MySQLClient.Create(&db_models.Reputation{
			SenderId:     session.UserId,
			RecipientId:  recipientId,
			Message:      repData.Message,
			Modifier:     repData.Modifier,
			CreationDate: time.Now(),
		}).Error
	}
	if err != nil {
		logger.Logf(logger.ERROR, "Could not modify user reputation. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var repCounter = 0
	err = database.MySQLClient.Model(&db_models.User{}).Select("reputation").Where("id = ?", recipientId).First(&repCounter).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not increase reputation counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Model(&db_models.User{}).Where("id = ?", recipientId).Update("reputation", repCounter+repData.Modifier).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not increase reputation counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully modified reputation of user",
	})
}
