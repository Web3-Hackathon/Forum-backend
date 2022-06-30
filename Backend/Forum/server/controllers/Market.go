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
	"math"
	"time"
)

// FetchMarketListings returns all market listings on the requested page
func FetchMarketListings(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var err error
	var ok bool
	var pageStr string
	var category string
	var page int

	pageStr, ok = utils.ParseParam("page", ctx, p)
	if !ok {
		return
	}

	category, ok = utils.ParseParam("category", ctx, p)
	if !ok {
		return
	}

	page, ok = utils.ConvertToInteger(pageStr, "FetchMarketListings", ctx)
	if !ok {
		return
	}

	if category != "purchase" && category != "sale" {
		utils.Error(ctx, 5, "Listing category must be either purchase or sale")
		return
	}

	var listings []db_models.MarketListing

	database.MySQLClient.Where("hidden = 0 AND type = ?", category).Offset(page*10 - 10).Limit(10).
		Order("creation_date DESC").Find(&listings)
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch listings from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	for k, listing := range listings {
		listing.SellerInfo = (*db_models.UserCardInfo)(GetUserCard(listing.SellerId, true))
		if listing.SellerInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch seller username from database for listings")
			continue
		}

		listings[k] = listing
	}

	var totalPages int64

	err = database.MySQLClient.Model(&[]db_models.MarketListing{}).Count(&totalPages).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not count total listings. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}
	if totalPages == 0 {
		totalPages = 1
	}

	totalPages = int64(math.Ceil(float64(totalPages) / 10.0))

	utils.JSON(ctx, models.MarketplaceResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched listings",
		},
		TotalPages: totalPages,
		Listings:   listings})
}

// PostMarketListing is used to create a new market listing
func PostMarketListing(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	// TODO: Make me match the most expensive role
	// if 3 > session.Rank {
	// 	utils.Error(ctx, 2, "You do not have the permissions to post in this section")
	// 	return
	// }

	var listingData = new(models.NewMarketListingRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &listingData) == nil {
		return
	}

	if len(listingData.Title) < 3 || len(listingData.Title) > 128 || listingData.Title == "" {
		utils.Error(ctx, 5, "Listings titles must be longer than 3 characters and shorter than 128 characters.")
		return
	}

	if len(listingData.Description) < 3 || len(listingData.Description) > 1024 || listingData.Description == "" {
		utils.Error(ctx, 5, "Listing content must be longer than 3 characters. And shorter than 1024 characters")
		return
	}

	if listingData.Price < 1 {
		utils.Error(ctx, 5, "Listing price must be more than 1$")
		return
	}

	if listingData.Type != "purchase" && listingData.Type != "sale" {
		utils.Error(ctx, 5, "Listing type must be either purchase or sale")
		return
	}

	if listingData.Currency != "ETH" && listingData.Currency != "SOL" {
		utils.Error(ctx, 5, "Listings only accept ETH & SOL transactions")
		return
	}

	var newListing = db_models.MarketListing{
		Title:        listingData.Title,
		SellerId:     session.UserId,
		CreationDate: time.Now(),
		Price:        listingData.Price,
		Currency:     listingData.Currency,
		Description:  listingData.Description,
		Type:         listingData.Type,
	}

	var response = database.MySQLClient.Create(&newListing)
	if response.Error != nil {
		logger.Logf(logger.ERROR, "Could not create market listing. Error: %s", response.Error)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var err = database.MySQLClient.Model(&db_models.User{}).
		Where("id=?", session.UserId).
		UpdateColumn("market_listings", gorm.Expr("market_listings + 1")).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update user profile on new market listing. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.NewThreadResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully posted listing",
		},
		Path: fmt.Sprintf("/market/listings/%s/%d",
			newListing.Type, newListing.Id),
	})
}

// FetchMarketListing is used to fetch a single market listing
func FetchMarketListing(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var listingIDStr string
	var listingID int
	var ok bool

	listingIDStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}
	listingID, ok = utils.ConvertToInteger(listingIDStr, "FetchMarketListing", ctx)
	if !ok {
		return
	}

	var listing db_models.MarketListing
	if database.IsHidden(listingID, &listing) {
		logger.Logf(logger.WARNING, "User tried accessing hidden listing")
		utils.Error(ctx, 3, "Resource not found")
		return
	}

	var sellerInfo = GetUserCard(listing.SellerId, false)
	if sellerInfo == nil {
		logger.Logf(logger.ERROR, "Could not parse seller from listing %d", listingID)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.MarketListingResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched listings",
		},
		Listing:    listing,
		SellerInfo: sellerInfo,
	})
}

func FetchUserMarketListings(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var err error
	var ok bool
	var username string
	var page int

	username, ok = utils.ParseParam("username", ctx, p)
	if !ok {
		return
	}

	var userId uint
	userId, ok = database.GetUserId(username)
	if !ok {
		utils.Error(ctx, 3, "User does not exist")
		return
	}

	var listings []db_models.MarketListing

	database.MySQLClient.Where("seller_id = ? AND hidden = 0", userId).Offset(page*10 - 10).Limit(10).
		Order("creation_date DESC").Find(&listings)
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch listings from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	for k, listing := range listings {
		listing.SellerInfo = (*db_models.UserCardInfo)(GetUserCard(listing.SellerId, true))
		if listing.SellerInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch seller card from database for listings.")
			continue
		}

		listings[k] = listing
	}

	var totalPages int64

	err = database.MySQLClient.Model(&[]db_models.MarketListing{}).Count(&totalPages).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not count total listings. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}
	if totalPages == 0 {
		totalPages = 1
	}

	totalPages = int64(math.Ceil(float64(totalPages) / 10.0))

	utils.JSON(ctx, models.MarketplaceResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched listings",
		},
		TotalPages: totalPages,
		Listings:   listings})
}

func DeleteListing(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var listingIDStr string
	var listingID int
	var ok bool

	listingIDStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}
	listingID, ok = utils.ConvertToInteger(listingIDStr, "FetchMarketListing", ctx)
	if !ok {
		return
	}

	var listing db_models.MarketListing
	if database.IsHidden(listingID, &listing) {
		logger.Logf(logger.WARNING, "User tried accessing hidden listing")
		utils.Error(ctx, 3, "Resource not found")
		return
	}

	var seller db_models.User
	var err = database.MySQLClient.Where("id = ?", listing.SellerId).First(&seller).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not parse seller from listing %d", listingID)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	if seller.Id != session.UserId {
		logger.Logf(logger.INFO, "User %d has attempted to delete foreign thread", session.UserId)
		utils.Error(ctx, 2, "You have not created this thread. Action logged.")
		return
	}

	listing.Hidden = true
	err = database.MySQLClient.Save(&listing).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not hide listing %d. Error: %s", listingID, err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully deleted listing",
	})
}
