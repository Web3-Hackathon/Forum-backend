package models

import "github.com/alxalx14/CryptoForum/Backend/Forum/models/db_models"

// VouchUserRequest is used when a user sends a vouch
type VouchUserRequest struct {
	Message    string `json:"message"`
	DealAmount int    `json:"dealAmount"`
	ShowAmount bool   `json:"showAmount"`
}

// UserReputationRequest is used when a user modifies
// the reputation of another user
type UserReputationRequest struct {
	Message  string `json:"message"`
	Modifier int    `json:"modifier"`
}

// VouchesResponseModel is used when fetching a users
// vouches
type VouchesResponseModel struct {
	BaseResponseModel
	Vouches []db_models.Vouches `json:"vouches"`
}

// ReputationResponseModel is used when fetching a users
// reputation
type ReputationResponseModel struct {
	BaseResponseModel
	Reputation []db_models.Reputation `json:"reputation"`
}
