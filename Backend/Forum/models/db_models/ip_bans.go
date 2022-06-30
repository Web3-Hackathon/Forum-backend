package db_models

import (
	"time"
)

type IPBan struct {
	GormModel `json:"-"`
	Id        uint      `json:"id" gorm:"primaryKey"`
	IpAddress string    `json:"ipAddress"`
	BanExpiry time.Time `json:"banExpiry"`
	BannedBy  string    `json:"bannedBy"`
	BanReason string    `json:"banReason"`
}
