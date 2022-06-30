package models

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
)

type NewMarketListingRequest struct {
	Price       int64  `json:"price"`
	Title       string `json:"title"`
	Currency    string `json:"currency"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// MarketplaceResponseModel is used when fetching a list of
// ThreadInfo from a given thread-section
type MarketplaceResponseModel struct {
	BaseResponseModel
	TotalPages int64                     `json:"totalPages"`
	Listings   []db_models.MarketListing `json:"listings"`
}

// MarketListingResponseModel is used to fetch a single marketplace
// listing
type MarketListingResponseModel struct {
	BaseResponseModel
	Listing    db_models.MarketListing `json:"listing"`
	SellerInfo *UserCardInfo           `json:"userInfo"`
}
