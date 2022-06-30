package database

import (
	"context"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var ThreadMongoCollection *mongo.Collection
var ReplyMongoCollection *mongo.Collection
var MongoContext = context.TODO()
var MongoClient *mongo.Client

// ConnectMongo is used to connect to our mongo DB instance
// only needed to fetch contents of threads and their replies
func ConnectMongo(url, username, password string) {
	var err error

	var clientOptions = options.Client().SetHosts(
		[]string{url},
	).SetAuth(
		options.Credential{
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-256",
			Username:      username,
			Password:      password,
		},
	)

	MongoClient, err = mongo.Connect(MongoContext, clientOptions)
	if err != nil {
		logger.Logf(logger.ERROR, "Could not connect to MongoDB. Error: %s", err.Error())
		os.Exit(0)
	}

	ThreadMongoCollection = MongoClient.Database("crypto_forum").Collection("threads")
	if ThreadMongoCollection == nil {
		logger.Logf(logger.ERROR, "Could not find collection 'threads' in mongodb.")
		os.Exit(0)
	}

	ReplyMongoCollection = MongoClient.Database("crypto_forum").Collection("replies")
	if ReplyMongoCollection == nil {
		logger.Logf(logger.ERROR, "Could not find collection 'replies' in mongodb.")
		os.Exit(0)
	}
}

// GetThreadContent is used to fetch the content of a thread
// from the mongo database
func GetThreadContent(threadId int) *db_models.ThreadContent {
	var err error
	var thread = new(db_models.ThreadContent)

	err = ThreadMongoCollection.FindOne(
		MongoContext,
		bson.M{
			"thread_id": threadId,
		}).Decode(thread)
	if err == mongo.ErrNoDocuments {
		return nil
	}
	if err != nil {
		logger.Logf(logger.ERROR, "Could not fetch thread content. Error: %s", err.Error())
		return nil
	}

	return thread
}

// GetThreadReplies is used to fetch the replies
// of a given thread on the given page
func GetThreadReplies(threadId int64, page int) []db_models.ReplyContent {
	var err error
	var replies []db_models.ReplyContent
	var opts *options.FindOptions
	var cursor *mongo.Cursor

	opts = options.Find().SetSort(
		bson.D{
			{"reply_index", 1},
		},
	).SetLimit(10)

	cursor, err = ReplyMongoCollection.Find(
		MongoContext,
		bson.D{
			{"reply_index", bson.M{"$gt": page*10 - 10}},
			{"thread_id", threadId},
		},
		opts)
	if err == mongo.ErrNoDocuments {
		return []db_models.ReplyContent{}
	}
	if err != nil {
		logger.Logf(logger.ERROR, "Could not fetch thread replies. Error: %s", err.Error())
		return []db_models.ReplyContent{}
	}

	err = cursor.All(MongoContext, &replies)
	if err != nil {
		logger.Logf(logger.ERROR, "Could not decode thread replies. Error: %s", err.Error())
		return []db_models.ReplyContent{}
	}

	return replies
}

// CountThreadReplies is used to count the amount of
// replies a thread has, used to count available pages
func CountThreadReplies(threadId int64) int64 {
	var count int64
	var err error

	count, err = ReplyMongoCollection.CountDocuments(MongoContext, bson.M{"thread_id": threadId})
	if err != nil {
		logger.Logf(logger.ERROR, "Could not count thread replies. Error: %s", err.Error())
		return -1
	}

	return count
}
