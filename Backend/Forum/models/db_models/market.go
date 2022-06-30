package db_models

import (
	"time"
)

type MarketListing struct {
	GormModel    `json:"-"`
	Id           uint          `json:"id"`
	SellerId     uint          `json:"sellerId"`
	CreationDate time.Time     `json:"creationDate"`
	Price        int64         `json:"price"`
	Title        string        `json:"title"`
	Currency     string        `json:"currency"`
	Description  string        `json:"description"`
	Hidden       bool          `json:"hidden"`
	Type         string        `json:"type"`
	SellerInfo   *UserCardInfo `json:"userInfo" gorm:"-"`
}
