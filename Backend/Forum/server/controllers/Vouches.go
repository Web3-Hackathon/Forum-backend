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

// GetUserVouches is used to get all vouches associated
// with a specific user
func GetUserVouches(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
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

	var vouches []db_models.Vouches

	err = database.MySQLClient.Where("recipient_id = ?", recipientId).Order("creation_date DESC").Find(&vouches).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch vouches from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.VouchesResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched vouches",
		},
		Vouches: vouches,
	})
}

// VouchUser is used to leave a vouch for a user after a deal
func VouchUser(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var username string
	var ok bool
	var err error

	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

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
		utils.Error(ctx, 9, "Cannot vouch/modify user reputation of yourself")
		return
	}

	var vouchData = new(models.VouchUserRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, vouchData) == nil {
		return
	}

	if vouchData.Message == "" || vouchData.DealAmount < 5 {
		utils.Error(ctx, 5, "Invalid request body")
		return
	}

	if !vouchData.ShowAmount {
		vouchData.DealAmount = -1
	}

	var vouch db_models.Vouches
	var exists = true

	err = database.MySQLClient.Where("sender_id = ? AND recipient_id = ?",
		session.UserId, recipientId).First(&vouch).Error
	if err == gorm.ErrRecordNotFound {
		exists = false
	} else if err != nil {
		logger.Logf(logger.ERROR, "Could not check if a user can rep another user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	if exists {
		utils.Error(ctx, 8, "You can vouch user just 1 time")
		return
	}

	err = database.MySQLClient.Create(&db_models.Vouches{
		SenderId:     session.UserId,
		RecipientId:  recipientId,
		Message:      vouchData.Message,
		DealAmount:   vouchData.DealAmount,
		ShowAmount:   vouchData.ShowAmount,
		CreationDate: time.Now(),
	}).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not vouch a user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var vouchCounter = 0
	err = database.MySQLClient.Model(&db_models.User{}).Select("vouches").Where("id = ?", recipientId).First(&vouchCounter).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not increase vouch counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Model(&db_models.User{}).Where("id = ?", recipientId).Update("vouches", vouchCounter+1).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not increase vouch counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully vouched for user",
	})
}

// DeleteUserVouch is used to remove a vouch from a user
func DeleteUserVouch(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var username string
	var ok bool
	var err error

	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

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

	var vouch db_models.Vouches

	err = database.MySQLClient.Where("sender_id = ? AND recipient_id = ?",
		session.UserId, recipientId).First(&vouch).Error
	if err == gorm.ErrRecordNotFound {
		utils.Error(ctx, 10, "You have not vouched for user")
		return
	} else if err != nil {
		logger.Logf(logger.ERROR, "Could not check if a user can vouch another user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Where("sender_id = ? AND recipient_id = ?",
		session.UserId, recipientId).Delete(&vouch).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not remove the vouch of a user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var vouchCounter = 0
	err = database.MySQLClient.Model(&db_models.User{}).Select("vouches").Where("id = ?", recipientId).First(&vouchCounter).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not decrease vouches counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Model(&db_models.User{}).Where("id = ?", recipientId).Update("vouches", vouchCounter-1).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not decrease vouches counter of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully removed vouch",
	})

}
