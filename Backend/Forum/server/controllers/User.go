package controllers

import (
	"fmt"
	"github.com/alxalx14/CryptoForum/Backend/Forum/database"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/middleware"
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gorm.io/gorm"
	"net/url"
)

func FetchUserInfo(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var err error

	var username, ok = utils.ParseParam("username", ctx, p)
	if !ok {
		return
	}

	var user = new(db_models.User)

	err = database.MySQLClient.Where("username = ?", username).First(user).Error
	if err == gorm.ErrRecordNotFound {
		utils.Error(ctx, 3, "User does not exist")
		return
	} else if err != nil {
		logger.Logf(logger.ERROR, "Could not fetch a users profile info. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var jsonData []byte
	jsonData, err = user.UserGroups.MarshalJSON()
	if err != nil {
		logger.Logf(logger.ERROR, "Could parse the user groups of user %s. Error: %s",
			user.Username, err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var userGroups = utils.ParseMySQLJson(jsonData).([]interface{})
	if userGroups == nil {
		err = fmt.Errorf("got back invalid JSON response from MySQL")
	}
	if err != nil {
		logger.Logf(logger.ERROR, "Could not parse a users user groups. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, &models.UserProfileResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched user",
		},
		Info: struct {
			Username      string      `json:"username"`
			UserGroups    interface{} `json:"userGroups"`
			JoinDate      string      `json:"joinDate"`
			Signature     string      `json:"signature"`
			PictureLink   string      `json:"pictureLink"`
			VerifiedMark  bool        `json:"verifiedMark"`
			WalletAddress string      `json:"walletAddress"`
		}{
			user.Username,
			userGroups,
			utils.FormatTime(user.JoinDate),
			user.Signature,
			user.PictureLink,
			user.VerifiedMark,
			user.WalletAddress,
		},
		Stats: struct {
			Posts          int64 `json:"posts"`
			Threads        int64 `json:"threads"`
			Reputation     int64 `json:"reputation"`
			MarketListings int64 `json:"marketListings"`
			Vouches        int64 `json:"vouches"`
		}{
			user.Posts,
			user.Threads,
			user.Reputation,
			user.MarketListings,
			user.Vouches,
		},
		Restrictions: struct {
			Banned     bool   `json:"banned"`
			BanReason  string `json:"banReason"`
			BanExpiry  string `json:"banExpiry"`
			Muted      bool   `json:"muted"`
			MuteReason string `json:"muteReason"`
			MuteExpiry string `json:"muteExpiry"`
		}{
			user.Banned,
			user.BanReason,
			utils.FormatTime(user.BanExpiry),
			user.Muted,
			user.MuteReason,
			utils.FormatTime(user.MuteExpiry),
		},
		ContactInfo: struct {
			DiscordID  string `json:"discordID"`
			DiscordTag string `json:"discordTag"`
			Telegram   string `json:"telegram"`
		}{
			user.DiscordID,
			user.DiscordTag,
			user.Telegram,
		},
	})
}

// TODO: On username change, log all instances of the user out

// UpdateUserProfile is used to update the profile info of the own user
func UpdateUserProfile(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var user db_models.User
	var err = database.MySQLClient.Where("id = ?", session.UserId).First(&user).Error
	if err != nil {
		logger.Logf(logger.ERROR, "")
		utils.Error(ctx, 3, "Your user could not be found. Contact support.")
		return
	}

	var newProfileData = new(models.UpdateUserProfile)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &newProfileData) == nil {
		return
	}

	if newProfileData.Username != user.Username {
		if user.ChangeUsername {
			user.Username = newProfileData.Username

			// TODO: Log user out
			user.ChangeUsername = false
		}

	}

	user.Signature = newProfileData.Signature
	user.DiscordID = newProfileData.DiscordId
	user.DiscordTag = newProfileData.DiscordTag
	user.Telegram = newProfileData.Telegram

	err = database.MySQLClient.Save(&user).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update user profile while updating user profile. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully updated profile0",
	})
}

func UpdateUserPicture(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var pictureData = new(models.NewUserProfilePicture)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &pictureData) == nil {
		return
	}

	var pictureLink = pictureData.PictureLink
	var err error

	// TODO: check if user owns the NFT and if he does use the NFTs metadata URL
	if pictureData.UseNFT {

	} else {
		// checking url
		var parsedUrl *url.URL
		parsedUrl, err = url.Parse(pictureData.PictureLink)
		if err != nil {
			utils.Error(ctx, 5, "Invalid picture link. Example link: https://imgur.com/a/<picture_link>")
			return
		}

		if parsedUrl.Host != "imgur.com" && parsedUrl.Host != "imgbb.com" {
			utils.Error(ctx, 5, "Only allowed domains: imgurl & imgbb. Do not attempt to use other picture hosts.")
			return
		}
	}

	err = database.MySQLClient.Model(&db_models.User{}).Where("id = ?", session.UserId).Update("picture_link", pictureLink).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update picture link of user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully updated profile picture link",
	})
}

// GetUserCard is used to get a brief description of the user
// to display as a card or next to a thread/listing
// NOT AN ENDPOINT
func GetUserCard(userId uint, signature bool) *models.UserCardInfo {
	var user db_models.User

	var err = database.MySQLClient.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil
	}

	if !signature {
		user.Signature = ""
	}

	return &models.UserCardInfo{
		Username:     user.Username,
		Signature:    user.Signature,
		VerifiedMark: user.VerifiedMark,
		PictureLink:  user.PictureLink,
		Reputation:   user.Reputation,
		Vouches:      user.Vouches,
		JoinDate:     utils.FormatTime(user.JoinDate),
	}
}
