package server

import (
	"fmt"
	"github.com/AubSs/fasthttplogger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/controllers"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
	"log"
	"time"
)

func StartServer(host string, port int) {
	var router = fasthttprouter.New()

	// Authenticate user
	router.POST("/authenticate", controllers.Authenticate)

	// Fetch profile info
	router.GET("/users/:username", controllers.FetchUserInfo)

	// Get all vouches of a user
	router.GET("/users/:username/vouches", controllers.GetUserVouches)
	// Vouch for a user
	router.POST("/users/:username/vouches", controllers.VouchUser)
	// Delete a vouch for a user
	router.DELETE("/users/:username/vouches", controllers.DeleteUserVouch)

	// Modify a users reputation
	router.POST("/users/:username/reputation", controllers.ModifyUsersReputation)
	// Get all feedback of a user
	router.GET("/users/:username/reputation", controllers.GetUsersReputation)

	// Send a message to a user
	router.POST("/users/:username/messages", controllers.SendMessage)
	// Get conversation with user
	router.GET("/users/:username/messages", controllers.FetchMessages)

	// Get all threads made by user
	router.GET("/users/:username/threads/:page", controllers.FetchUserThreads)
	// Get all marketplace listings by user
	router.GET("/users/:username/listings/:page", controllers.FetchUserMarketListings)

	// Edit user profile
	router.PATCH("/me/profile", controllers.UpdateUserProfile)
	// Upload profile picture
	router.POST("/me/picture", controllers.UpdateUserPicture)
	// Select NFT as profile picture
	router.POST("/me/picture/nft/:id", nil)

	// Fetch all categories
	router.GET("/sections", controllers.FetchSections)
	// Post a thread
	router.POST("/sections/:id/new", controllers.PostThread)
	// Fetch all threads of page N in category X
	router.GET("/sections/:id/threads/:page", controllers.FetchThreads)

	// Fetch info of a thread Y in section X
	router.GET("/threads/:id", controllers.FetchThread)
	// Reply to a thread
	router.POST("/threads/:id/new", controllers.PostReply)
	// Fetch all replies of page N in thread X
	router.GET("/threads/:id/pages/:page", controllers.FetchThreadReplies)

	// Fetch all market ads from market category
	router.GET("/market/:category/:page", controllers.FetchMarketListings)
	// Post new market ad
	router.POST("/market/new", controllers.PostMarketListing)
	// Fetch a market ad
	router.GET("/listings/:id", controllers.FetchMarketListing)
	// Delete market ad
	router.DELETE("/listings/:id", controllers.DeleteListing)

	// Create new payment
	router.POST("/payments/new", nil)
	// Payment webhook only accessible by Coinbase Commerce
	router.POST("/payment/coinbase/webhook/update", nil)

	var server = &fasthttp.Server{
		Handler:         fasthttplogger.CombinedColored(router.Handler),
		Name:            "CryptoForum",
		Concurrency:     512 * 2048,
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		ReadTimeout:     0,
		WriteTimeout:    0,
		IdleTimeout:     5 * time.Second,
	}

	logger.Logf(logger.INFO, "Starting HTTP server on %s:%d", host, port)
	var err = server.ListenAndServe(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		log.Fatalf("[SERVER] Could not listen on %s:%d. Error: %s\n", host, port, err.Error())
	}
}
