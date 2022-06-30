package database

import (
	"fmt"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	// _ "github.com/go-sql-driver/mysql"
)

var MySQLClient *gorm.DB // *sql.DB

// ConnectMySQL is used to connect to the MySQL database
// it holds meta info about threads and listings but no actual
// content of either one
//func ConnectMySQL(address, username, password string) {
//	var err error
//
//	MySQLClient, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/crypto_forum",
//		username, password, address))
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not connect to MySQL. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	// See "Important settings" section.
//	MySQLClient.SetConnMaxLifetime(time.Minute * 5)
//	MySQLClient.SetMaxOpenConns(1000 ^ 5)
//	MySQLClient.SetMaxIdleConns(1000 ^ 5)
//}
func ConnectMySQL(address, username, password string) {
	var err error

	MySQLClient, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/crypto_forum?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, address)))
	if err != nil {
		logger.Logf(logger.ERROR, "Could not connect to MySQL. Error: %s", err.Error())
		os.Exit(0)
	}

	err = MySQLClient.AutoMigrate(&db_models.User{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate User model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.Vouches{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Vouches model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.Reputation{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Reputation model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.Thread{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Thread model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.ThreadSection{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate ThreadSection model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.IPBan{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate IPBan model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.Rank{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Rank model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.MarketListing{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Rank model. Error: %s", err.Error())
		os.Exit(1)
	}

	err = MySQLClient.AutoMigrate(&db_models.Message{})
	if err != nil {
		logger.Logf(logger.ERROR, "could not autoMigrate Message model. Error: %s", err.Error())
		os.Exit(1)
	}
}

// UsernameExists is used to check if a username currently exists in the
// database
func UsernameExists(username string) (bool, error) {
	var user db_models.User
	var exists bool

	var err = MySQLClient.Model(&user).Select("count(*) > 0").Where("username = ?", username).First(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetUserId is used to get a users ID by their username
func GetUserId(username string) (uint, bool) {
	var user db_models.User
	var userId uint

	var err = MySQLClient.Model(&user).Select("id").Where("username = ?", username).First(&userId).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false
	} else if err != nil {
		return 0, false
	}

	return userId, true
}

// GetUsernameById is used to get a users username by their ID
func GetUsernameById(id uint) (string, bool) {
	var user db_models.User
	var username string

	var err = MySQLClient.Model(&user).Select("username").Where("id = ?", id).First(&username).Error
	if err == gorm.ErrRecordNotFound {
		return "", false
	} else if err != nil {
		return "", false
	}

	return username, true
}

// IsHidden is used to find out if the object is
// hidden or not (deleted)
func IsHidden(objectId int, object interface{}) bool {
	var err error

	err = MySQLClient.Where("hidden = 0 AND id = ?", objectId).First(object).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Logf(logger.ERROR, "Could not check if an object is hidden. Error: %s", err.Error())
		return true
	} else if err == gorm.ErrRecordNotFound {
		return true
	}

	return false
}
