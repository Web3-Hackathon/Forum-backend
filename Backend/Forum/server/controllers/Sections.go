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
)

// FetchSections returns all available thread-sections on the forum
func FetchSections(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	var err error

	if middleware.Middleware(ctx) == nil {
		return
	}

	var sections []db_models.ThreadSection

	err = database.MySQLClient.Find(&sections).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not fetch sections from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var result = make(map[string]map[string][]map[string]uint)

	var ok bool
	for _, section := range sections {
		if _, ok = result[section.Category]; !ok {
			result[section.Category] = make(map[string][]map[string]uint)
			result[section.Category][section.Parent] = []map[string]uint{
				{section.Name: section.Id},
			}
		} else {
			result[section.Category][section.Parent] =
				append(result[section.Category][section.Parent], map[string]uint{section.Name: section.Id})
		}
	}

	utils.JSON(ctx, models.SectionsResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched sections",
		},
		Sections: result,
	})
}
