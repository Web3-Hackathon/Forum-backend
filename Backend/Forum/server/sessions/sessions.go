package sessions

import (
	"github.com/alxalx14/CryptoForum/Backend/Forum/utils"
	cmap "github.com/orcaman/concurrent-map"
	"time"
)

// Metadata contains metadata to make sure a session
// is not shared across devices, this data can not be fetched
// by any endpoint
type Metadata struct {
	UserAgent          string
	IpAddress          string
	BrowserFingerprint string
}

type Session struct {
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	// AccessLevel defines 3 different access levels
	// 0 - normal user (create threads, like threads, reply to thread)
	// 1 - moderator (edit threads, hide threads, warn users, ban users)
	// 2 - admin (all of the above, create thread sections, user management)
	AccessLevel int `json:"accessLevel"`
	// Rank is like a role that can be purchased by the user to make
	// their threads look cooler and make the user look fancier
	Rank int `json:"rank"`
	// ActiveSince is the timestamp that is created when the user logs in
	ActiveSince int64 `json:"activeSince"`
	// SessionId is the identifier that is used by the user in requests to
	// prove he is logged in
	SessionId string `json:"-"`
	// LastRequest is the timestamp of the last request sent by the user
	// it's used to remove inactive sessions if the last request was over
	// 3 days ago
	LastRequest int64 `json:"-"`
	// Metadata is a bot prevention technique that protects the server
	// against basic & noob wanna-be hackers. At login, we store info
	// about the user, if on any request that info does not match
	// the user is logged out on all sessions. This data is not accessible
	// to anyone
	Metadata Metadata `json:"-"`
}

// Active is a concurrent Map that keeps track of all active sessions
var Active = cmap.New()

// Create returns a new session id that is linked to an entry
// inside Sessions
func Create(userId uint, username string, accessLevel int, metadata Metadata) string {
	var sessionId = utils.RandomString(128)

	Active.Set(sessionId, &Session{
		UserId:      userId,
		Username:    username,
		AccessLevel: accessLevel,
		ActiveSince: time.Now().Unix(),
		SessionId:   sessionId,
		LastRequest: time.Now().Unix(),
		Metadata:    metadata,
	})

	return sessionId
}

// Delete is used to log a user out, it also deletes the session
// inside Sessions
func Delete(sessionId string) {
	Active.Remove(sessionId)
}

// WatchDog runs in the background and checks all sessions every
// 5 minutes to make sure that they are still active, otherwise it removes
// them
func WatchDog() {
	for {
		var sessions = Active.Items()
		for sessionId, sessionInterface := range sessions {
			var session = sessionInterface.(*Session)

			// Deleting if 3 days passed
			if session.LastRequest+259200 < time.Now().Unix() {
				Delete(sessionId)
				continue
			}

			// TODO: Check if the rank of the user has expired or not
		}

		time.Sleep(30 * time.Second)
	}
}
