package queries

//import (
//	"database/sql"
//	"github.com/alxalx14/CryptoForum/Backend/Forum/database"
//	"github.com/alxalx14/CryptoForum/Backend/Forum/logger"
//	"os"
//)
//
//// Threads related SQL queries
//
//// FetchSections is used to get all available sections
//// for the forum
//var FetchSections *sql.Stmt
//
//// FetchThreads is used to get the metadata of all threads on
//// page N inside section X. It takes inputs session_id and the
//// offset (used for pagination)
//var FetchThreads *sql.Stmt
//
//// CountThreads counts all threads of a section, used to
//// count total pages
//var CountThreads *sql.Stmt
//
//// FetchUserPublic is used to return publicly available
//// on a user, needed to display a users profile
//var FetchUserPublic *sql.Stmt
//
//// GetUsersReputation is sued to get all feedback given to a user
//var GetUsersReputation *sql.Stmt
//
//// ModifyUserReputation is used to modify a user's
//// reputation, used by the GiveUserFeedback endpoint
//var ModifyUserReputation *sql.Stmt
//
//// VouchUser is used to leave a vouch for a service provided
//// by a user
//var VouchUser *sql.Stmt
//
//// CanVouchUser checks if the user has ever left a vouch
//// for this user and if so the new vouch is not allowed
//var CanVouchUser *sql.Stmt
//
//// GetUsersVouches is used to fetch all vouches of a user
//var GetUsersVouches *sql.Stmt
//
//// RemoveUserVouch is used to remove the vouch of a user
//var RemoveUserVouch *sql.Stmt
//
//// IsThreadHidden is used to check if a thread has hidden = 1
//var IsThreadHidden *sql.Stmt
//
//// ModifyUserReputationCounter is used to increase/decrease the rep counter
//var ModifyUserReputationCounter *sql.Stmt
//
//// ModifyUserVouchCounter is used to increase/decrease the vouch counter
//var ModifyUserVouchCounter *sql.Stmt
//
//// AuthenticateUser is used on login to check if a user exists
//var AuthenticateUser *sql.Stmt
//
//// RegisterUser is used on login to create a new account in case
//// the user has no account
//var RegisterUser *sql.Stmt
//
//// UpdateUserLastLogin is used to update the last login info of a user
//var UpdateUserLastLogin *sql.Stmt
//
//// PrepareMySQLStatements is used to prepare all statements
//// used by the app, every executed query is listed here
//func PrepareMySQLStatements() {
//	var err error
//
//	FetchSections, err = database.MySQLClient.Prepare("SELECT id, parent_category, parent_section, section_name FROM thread_sections")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare FetchCategories statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	CountThreads, err = database.MySQLClient.Prepare("SELECT COUNT(*) FROM threads WHERE section_id=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare FetchCategories statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	IsThreadHidden, err = database.MySQLClient.Prepare("SELECT hidden FROM threads WHERE id=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare FetchCategories statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	FetchThreads, err = database.MySQLClient.Prepare("SELECT id, title, author, creation_date, last_post, post_count, hidden FROM threads WHERE section_id=? ORDER BY creation_date DESC LIMIT ?,10")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare FetchCategories statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	FetchUserPublic, err = database.MySQLClient.Prepare("SELECT username, user_groups, join_date, wallet_address, posts, threads, likes, reputation, vouches, banned, ban_reason, ban_expiry, muted, mute_reason, mute_expiry, discord_id, discord_tag, telegram FROM users WHERE username=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare FetchUserPublic statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	GetUsersReputation, err = database.MySQLClient.Prepare("SELECT sender, recipient, message, modifier, creation_date FROM feedbacks WHERE recipient=? ORDER BY creation_date DESC")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare GetUsersVouches statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	// TODO: Implement a check to know if the reputation decreased and calculate by how much it decreased or increased!!
//	ModifyUserReputation, err = database.MySQLClient.Prepare("INSERT INTO feedbacks (sender, recipient, message, modifier, creation_date) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP())")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare ModifyUserReputation statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	GetUsersVouches, err = database.MySQLClient.Prepare("SELECT sender, recipient, message, deal_amount, show_amount, creation_date FROM vouches WHERE recipient=? ORDER BY creation_date DESC")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare GetUsersVouches statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	CanVouchUser, err = database.MySQLClient.Prepare("SELECT id FROM vouches WHERE recipient=? AND sender=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare GetUsersVouches statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	VouchUser, err = database.MySQLClient.Prepare("INSERT INTO vouches (sender, recipient, message, deal_amount, show_amount, creation_date) VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP())")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare VouchUser statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	RemoveUserVouch, err = database.MySQLClient.Prepare("DELETE FROM vouches WHERE sender=? AND recipient=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare RemoveUserVouch statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	ModifyUserVouchCounter, err = database.MySQLClient.Prepare("UPDATE users SET vouches = vouches + ? WHERE username=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare ModifyUserVouchCounter statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	ModifyUserReputationCounter, err = database.MySQLClient.Prepare("UPDATE users SET reputation = reputation + ? WHERE username=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare ModifyUserReputationCounter statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	AuthenticateUser, err = database.MySQLClient.Prepare("SELECT username, user_groups, banned, ban_reason, ban_expiry FROM users WHERE wallet_address=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare AuthenticateUser statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	RegisterUser, err = database.MySQLClient.Prepare("INSERT INTO users (username, email, last_ip, user_groups, last_login, join_date, wallet_address, posts, threads, likes, reputation, vouches, banned, ban_reason, ban_expiry, muted, mute_reason, mute_expiry, discord_id, discord_tag, telegram) VALUES (?, '', ?, '[]', CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP(), ?, 0, 0, 0, 0, 0, 0, '', NULL, 0, '', NULL, '', '', '')")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare RegisterUser statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//
//	UpdateUserLastLogin, err = database.MySQLClient.Prepare("UPDATE users SET last_login=CURRENT_TIMESTAMP(), last_ip=? WHERE username=?")
//	if err != nil {
//		logger.Logf(logger.ERROR, "Could not prepare UpdateUserLastLogin statement. Error: %s", err.Error())
//		os.Exit(0)
//	}
//}
