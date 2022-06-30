package models

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
)

type NewThreadRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ThreadContentResponseModel is used when fetching a
// specific thread from the Mongo DB. It contains the
// initial thread and the first page of the replies
type ThreadContentResponseModel struct {
	BaseResponseModel
	Thread     db_models.ThreadContent  `json:"thread"`
	TotalPages int64                    `json:"totalPages"`
	Replies    []db_models.ReplyContent `json:"replies"`
	AuthorInfo *UserCardInfo            `bson:"-" json:"userInfo"`
}

type NewThreadResponseModel struct {
	BaseResponseModel
	Path string `json:"threadPath"`
}

// ThreadRepliesResponseModel is used when fetching the
// replies from a specific thread from the Mongo DB. It
// contains 10 replies from the requested page
type ThreadRepliesResponseModel struct {
	BaseResponseModel
	Replies []db_models.ReplyContent `json:"replies"`
}

// ThreadsResponseModel is used when fetching a list of
// ThreadInfo from a given thread-section
type ThreadsResponseModel struct {
	BaseResponseModel
	TotalPages int64              `json:"totalPages"`
	Threads    []db_models.Thread `json:"threads"`
}
