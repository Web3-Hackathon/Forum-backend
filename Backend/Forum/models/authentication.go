package models

type LoginRequest struct {
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
}

type LoginResponse struct {
	BaseResponseModel
	SessionKey  string `json:"sessionKey"`
	AccessLevel int    `json:"accessLevel"`
}
