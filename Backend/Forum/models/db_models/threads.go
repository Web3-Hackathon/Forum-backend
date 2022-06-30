package db_models

import (
	"time"
)

type ThreadContent struct {
	ThreadId     int64  `bson:"thread_id" json:"threadId"`
	CreationDate int64  `bson:"creation_date" json:"creationDate"`
	Content      string `bson:"content" json:"content"`
}

type Thread struct {
	GormModel    `json:"-"`
	Id           uint          `json:"id" gorm:"primaryKey"`
	Title        string        `json:"title"`
	SectionId    uint          `json:"sectionId"`
	AuthorId     uint          `json:"author_id"`
	CreationDate time.Time     `json:"creationDate"`
	LastPost     time.Time     `json:"lastPost"`
	PostCount    uint64        `json:"postCount"`
	Views        uint64        `json:"views"`
	Hidden       bool          `json:"-"`
	AuthorInfo   *UserCardInfo `json:"userInfo" gorm:"-"`
}

type ReplyContent struct {
	AuthorId     int64  `bson:"author" json:"author_id"`
	Content      string `bson:"content" json:"content"`
	ThreadId     int64  `bson:"thread_id" json:"threadId"`
	ReplyIndex   int64  `bson:"reply_index" json:"replyIndex"`
	CreationDate int64  `bson:"creation_date" json:"creationDate"`
	//Hidden       bool   `bson:"hidden" json:"-"`
	UserInfo *UserCardInfo `bson:"-" json:"userInfo"`
}
