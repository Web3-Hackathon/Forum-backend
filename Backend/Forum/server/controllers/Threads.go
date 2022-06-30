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

// FetchThreads returns all threads on the requested page
// in the requested section
func FetchThreads(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var err error
	var ok bool
	var sectionIDStr string
	var sectionID int
	var pageStr string
	var page int

	sectionIDStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}
	pageStr, ok = utils.ParseParam("page", ctx, p)
	if !ok {
		return
	}

	page, ok = utils.ConvertToInteger(pageStr, "FetchThreads", ctx)
	if !ok {
		return
	}
	sectionID, ok = utils.ConvertToInteger(sectionIDStr, "FetchThreads", ctx)
	if !ok {
		return
	}

	var threads []db_models.Thread

	database.MySQLClient.Where("section_id = ? AND hidden = 0", sectionID).Offset(page*10 - 10).Limit(10).
		Order("creation_date DESC").Find(&threads)
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch threads from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	for k, thread := range threads {
		thread.AuthorInfo = (*db_models.UserCardInfo)(GetUserCard(thread.AuthorId, true))
		if thread.AuthorInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch authors username from database on for threads")
			continue
		}

		threads[k] = thread
	}

	var totalPages int64

	err = database.MySQLClient.Model(&[]db_models.Thread{}).Count(&totalPages).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not count total threads in section. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}
	if totalPages == 0 {
		totalPages = 1
	}

	totalPages = int64(math.Ceil(float64(totalPages) / 10.0))

	utils.JSON(ctx, models.ThreadsResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched threads",
		},
		TotalPages: totalPages,
		Threads:    threads})
}

// FetchThread gets the contents of a thread and other info about it
func FetchThread(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var threadIDStr string
	var threadID int
	var ok bool

	threadIDStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}
	threadID, ok = utils.ConvertToInteger(threadIDStr, "FetchThread", ctx)
	if !ok {
		return
	}

	if database.IsHidden(threadID, &db_models.Thread{}) {
		logger.Logf(logger.WARNING, "User tried accessing hidden thread")
		utils.Error(ctx, 3, "Resource not found")
		return
	}

	var threadContent = database.GetThreadContent(threadID)
	if threadContent == nil {
		logger.Logf(logger.WARNING, "User tried accessing non existent thread")
		utils.Error(ctx, 3, "Resource not found")
		return
	}

	var threadReplies = database.GetThreadReplies(int64(threadID), 1)
	for k, reply := range threadReplies {
		var user db_models.User
		database.MySQLClient.Where("id = ?", reply.AuthorId).First(&user)

		reply.UserInfo = (*db_models.UserCardInfo)(GetUserCard(uint(reply.AuthorId), true))
		if reply.UserInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch users from database for thread replies")
			continue
		}

		threadReplies[k] = reply
	}

	var totalReplies = database.CountThreadReplies(int64(threadID))
	if totalReplies == -1 {
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var thread db_models.Thread
	var err = database.MySQLClient.Where("id=?", threadID).First(&thread).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not parse threadId %d from database", threadID)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	thread.Views++
	err = database.MySQLClient.Save(&thread).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update view counter for threadId %d", threadID)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	var authorInfo = GetUserCard(thread.AuthorId, true)
	if authorInfo == nil {
		logger.Logf(logger.ERROR, "Could not parse author from thread %d", threadID)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.ThreadContentResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched thread",
		},
		TotalPages: int64(math.Ceil(float64(totalReplies) / 10.0)),
		Thread:     *threadContent,
		AuthorInfo: authorInfo,
		Replies:    threadReplies,
	})
}

// FetchThreadReplies is used to fetch all replies on a thread, it is paged
func FetchThreadReplies(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var threadIDStr string
	var threadID int
	var replyPageStr string
	var replyPage int
	var ok bool

	threadIDStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}
	threadID, ok = utils.ConvertToInteger(threadIDStr, "FetchThreadReplies", ctx)
	if !ok {
		return
	}

	replyPageStr, ok = utils.ParseParam("page", ctx, p)
	if !ok {
		return
	}
	replyPage, ok = utils.ConvertToInteger(replyPageStr, "FetchThreadReplies", ctx)
	if !ok {
		return
	}

	if database.IsHidden(threadID, &db_models.Thread{}) {
		logger.Logf(logger.WARNING, "User tried accessing hidden thread")
		utils.Error(ctx, 3, "Resource not found")
		return
	}

	var threadReplies = database.GetThreadReplies(int64(threadID), replyPage)

	for k, reply := range threadReplies {
		var user db_models.User
		database.MySQLClient.Where("id = ?", reply.AuthorId).First(&user)

		reply.UserInfo = (*db_models.UserCardInfo)(GetUserCard(uint(reply.AuthorId), true))
		if reply.UserInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch users from database for thread replies")
			continue
		}

		threadReplies[k] = reply
	}

	utils.JSON(ctx, models.ThreadRepliesResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched replies",
		},
		Replies: threadReplies,
	})
}

// FetchUserThreads returns all threads created by a specific user
func FetchUserThreads(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	if middleware.Middleware(ctx) == nil {
		return
	}

	var err error
	var ok bool
	var pageStr string
	var page int
	var username string

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

	pageStr, ok = utils.ParseParam("page", ctx, p)
	if !ok {
		return
	}

	page, ok = utils.ConvertToInteger(pageStr, "FetchUserThreads", ctx)
	if !ok {
		return
	}

	var threads []db_models.Thread

	database.MySQLClient.Where("author_id = ? AND hidden = 0", userId).Offset(page*10 - 10).Limit(10).
		Order("creation_date DESC").Find(&threads)
	if err != gorm.ErrRecordNotFound && err != nil {
		logger.Logf(logger.ERROR, "Could not fetch threads from database. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	for k, thread := range threads {
		thread.AuthorInfo = (*db_models.UserCardInfo)(GetUserCard(thread.AuthorId, true))
		if thread.AuthorInfo == nil {
			logger.Logf(logger.WARNING, "Could not fetch authors username from database on for threads")
			continue
		}

		threads[k] = thread
	}

	var totalPages int64

	err = database.MySQLClient.Model(&[]db_models.Thread{}).Count(&totalPages).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not count total threads in section. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}
	if totalPages == 0 {
		totalPages = 1
	}

	totalPages = int64(math.Ceil(float64(totalPages) / 10.0))

	utils.JSON(ctx, models.ThreadsResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully fetched threads",
		},
		TotalPages: totalPages,
		Threads:    threads,
	})
}

// PostThread is used to create a new thread in the specified section
func PostThread(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var ok bool
	var sectionIdStr string
	var sectionId int

	sectionIdStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}

	sectionId, ok = utils.ConvertToInteger(sectionIdStr, "PostThread", ctx)
	if !ok {
		return
	}

	var section db_models.ThreadSection

	var err = database.MySQLClient.Where("id=?", sectionId).First(&section).Error
	if err != nil {
		utils.Error(ctx, 3, "Thread section does not exist")
		return
	}

	if section.MinimumRank > session.Rank {
		utils.Error(ctx, 2, "You do not have the permissions to post in this section")
		return
	}

	var threadData = new(models.NewThreadRequest)
	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &threadData) == nil {
		return
	}

	if len(threadData.Title) < 3 || len(threadData.Title) > 128 || threadData.Title == "" {
		utils.Error(ctx, 5, "Thread titles must be longer than 3 characters and shorter than 128 characters.")
		return
	}

	if len(threadData.Content) < 3 || threadData.Content == "" {
		utils.Error(ctx, 5, "Thread content must be longer than 3 characters.")
		return
	}

	var newThread = db_models.Thread{
		Title:        threadData.Title,
		SectionId:    uint(sectionId),
		AuthorId:     session.UserId,
		CreationDate: time.Now(),
		LastPost:     time.Now(),
		PostCount:    1,
		Views:        1,
		Hidden:       false,
	}

	var response = database.MySQLClient.Create(&newThread)
	if response.Error != nil {
		logger.Logf(logger.ERROR, "Could not create thread. Error: %s", response.Error)
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	_, err = database.ThreadMongoCollection.InsertOne(
		database.MongoContext,
		db_models.ThreadContent{
			ThreadId:     int64(newThread.Id),
			CreationDate: time.Now().Unix(),
			Content:      threadData.Content,
		},
	)
	if err != nil {
		logger.Logf(logger.ERROR, "Could not insert thread data. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Model(&db_models.User{}).
		Where("id=?", session.UserId).
		UpdateColumn("threads", gorm.Expr("threads + 1")).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update user profile on new thread. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.NewThreadResponseModel{
		BaseResponseModel: models.BaseResponseModel{
			Status:  true,
			Message: "Successfully posted thread",
		},
		Path: fmt.Sprintf("/threads/%d", newThread.Id),
	})
}

// PostReply is used to reply to a thread
func PostReply(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
	var session = middleware.Middleware(ctx)
	if session == nil {
		return
	}

	var ok bool
	var threadIdStr string
	var threadId int
	var err error

	threadIdStr, ok = utils.ParseParam("id", ctx, p)
	if !ok {
		return
	}

	threadId, ok = utils.ConvertToInteger(threadIdStr, "PostReply", ctx)
	if !ok {
		return
	}

	var replyData struct {
		Content string `json:"content"`
	}

	if utils.ParseBody(&utils.UserRequest{
		Username: session.Username,
		Request:  ctx,
	}, &replyData) == nil {
		return
	}

	if len(replyData.Content) < 3 {
		utils.Error(ctx, 5, "Reply content must be longer than 3 characters.")
		return
	}

	var thread db_models.Thread
	err = database.MySQLClient.
		Where("id = ?", threadId).
		First(&thread).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not get post count of thread. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	_, err = database.ReplyMongoCollection.InsertOne(
		database.MongoContext,
		db_models.ReplyContent{
			ReplyIndex:   int64(thread.PostCount),
			AuthorId:     int64(session.UserId),
			ThreadId:     int64(threadId),
			CreationDate: time.Now().Unix(),
			Content:      replyData.Content,
		},
	)
	if err != nil {
		logger.Logf(logger.ERROR, "Could not insert thread data. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	thread.PostCount++
	thread.LastPost = time.Now()

	err = database.MySQLClient.Save(&thread).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update user profile on new thread. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	err = database.MySQLClient.Model(&db_models.User{}).
		Where("id=?", session.UserId).
		UpdateColumn("posts", gorm.Expr("posts + 1")).Error
	if err != nil {
		logger.Logf(logger.ERROR, "Could not update user profile on new thread. Error: %s", err.Error())
		utils.Error(ctx, 7, "Error occurred on the server. Contact staff")
		return
	}

	utils.JSON(ctx, models.BaseResponseModel{
		Status:  true,
		Message: "Successfully posted reply",
	})

}
