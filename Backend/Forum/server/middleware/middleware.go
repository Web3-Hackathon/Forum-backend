package middleware

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/server/sessions"
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils"
	"github.com/valyala/fasthttp"
	"strings"
)

// Middleware is used to authenticate a users request on an endpoint
// this runs all checks related to the session
// - checking if it is still active
// - checking the metadata
// TODO: Add middleware with example data to start running live tests!
func Middleware(ctx *fasthttp.RequestCtx) *sessions.Session {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PATCH")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "false")

	//return &sessions.Session{
	//	Username:    "0xD",
	//	AccessLevel: 1,
	//	Rank:        2,
	//	ActiveSince: 0,
	//	SessionId:   "asfoasforo23",
	//	LastRequest: time.Now().Unix(),
	//	Metadata:    sessions.Metadata{},
	//}
	//

	var sessionId = string(ctx.Request.Header.Peek("Authorization"))
	sessionId = strings.Replace(strings.ToLower(sessionId), "bearer ", "", -1)
	var sessionInterface, exists = sessions.Active.Get(sessionId)

	if !exists {
		utils.Error(ctx, 1, "User is not authenticated")
		return nil
	}

	var session = sessionInterface.(*sessions.Session)

	// getting data to verify the metadata
	//var currentMetadata = sessions.Metadata{
	//	UserAgent:          string(ctx.Request.Header.UserAgent()),
	//	IpAddress:          realip.FromRequest(ctx),
	//	BrowserFingerprint: string(ctx.Request.Header.Peek("fingerprint")),
	//}
	//
	//if !cmp.Equal(session.Metadata, currentMetadata) {
	//	logger.Logf(logger.INFO, "The metadata of user %s has changed mid session. User has been logged out.",
	//		session.Username)
	//
	//	return nil
	//}

	return session
}
