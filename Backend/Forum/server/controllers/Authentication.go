package controllers

import (
	"encoding/hex"
	"fmt"
	"github.com/alxalx14/CryptoForum/Backend/Forum/database"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/sessions"
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils"
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils/realip"
	"github.com/gagliardetto/solana-go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

var LoginNonce = []byte("Welcome to CryptoForum, sign this message to login!")

func Authenticate(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	var loginData = new(models.LoginRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: "",
		Request:  ctx,
	}, &loginData) == nil {
		return
	}

	var pubKey, err = solana.PublicKeyFromBase58(loginData.PublicKey)
	if err != nil {
		logger.Logf(logger.ERROR, "Could verify login signature. Error: %s", err.Error())
		utils.Error(ctx, 5, "Error occurred on the server. Contact staff")
		return
	}

	var sig []byte
	sig, err = hex.DecodeString(loginData.Signature)
	if err != nil {
		utils.Error(ctx, 5, "Invalid request body")
		return
	}

	var signature = solana.SignatureFromBytes(sig)
	var ok = signature.Verify(pubKey, LoginNonce)
	if !ok {
		database.MySQLClient.Create(&db_models.IPBan{
			IpAddress: realip.FromRequest(ctx),
			BanExpiry: time.Now().Add((7 * 24) * time.Hour),
			BannedBy:  "SYSTEM",
			BanReason: "Attempted to use fake signature on login endpoint.",
		})
		logger.Logf(logger.ERROR, "Someone attempted to verify fake signature, IP has been banned.")

		utils.Error(ctx, 11, "Could not verify signature. IP has been banned")
		return
	}

	var user db_models.User

	err = database.MySQLClient.Where("wallet_address = ?", loginData.PublicKey).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user.Username = "user" + utils.RandomString(28)
		var exists = true

		for exists {
			exists, err = database.UsernameExists(user.Username)
			if err != nil {
				logger.Logf(logger.ERROR, "Could not check user existence. Error: %s", err.Error())
				utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
				return
			}
		}

		var emptyArray = datatypes.JSON{}
		err = emptyArray.UnmarshalJSON([]byte(`[]`))
		if err != nil {
			logger.Logf(logger.ERROR, "Could create user groups for user. Error: %s", err.Error())
			utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
			return
		}

		user.JoinDate = time.Now()
		user.WalletAddress = loginData.PublicKey
		user.LastIP = realip.FromRequest(ctx)
		user.LastLogin = time.Now()
		user.OldUsernames = emptyArray
		user.UserGroups = emptyArray
		user.BanExpiry = time.Now()
		user.MuteExpiry = time.Now()

		err = database.MySQLClient.Create(&user).Error
		if err != nil {
			logger.Logf(logger.ERROR, "Could register user. Error: %s", err.Error())
			utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
			return
		}

		err = database.MySQLClient.Where("wallet_address = ?", loginData.PublicKey).First(&user).Error
		if err != nil {
			logger.Logf(logger.ERROR, "Could not login user. Error: %s", err.Error())
			utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
			return
		}

	} else if err != nil {
		logger.Logf(logger.ERROR, "Could not login user. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	if user.Banned {
		utils.Error(ctx, 12, fmt.Sprintf("You are banned from the forum until %s. Reason: %s",
			utils.FormatTime(user.BanExpiry), user.BanReason))
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

	var accessLevel = 0

	var userGroups []interface{}

	userGroups, ok = utils.ParseMySQLJson(jsonData).([]interface{})
	if !ok {
		return
	} else {
		if utils.IntfSliceContains(userGroups, "admin") {
			accessLevel = 2
		} else if utils.IntfSliceContains(userGroups, "moderator") {
			accessLevel = 1
		} else {
			accessLevel = 0
		}
	}

	err = database.MySQLClient.Model(&user).Where("wallet_address = ?", loginData.PublicKey).
		Update("last_login", time.Now()).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update users last login info. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var sessionKey = sessions.Create(user.Id, user.Username, accessLevel, sessions.Metadata{
		UserAgent:          string(ctx.Request.Header.UserAgent()),
		IpAddress:          realip.FromRequest(ctx),
		BrowserFingerprint: string(ctx.Request.Header.Peek("fingerprint")),
	})

	utils.JSON(ctx, models.LoginResponse{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully logged in, welcome",
		},
		SessionKey:  sessionKey,
		AccessLevel: accessLevel,
	})
}
