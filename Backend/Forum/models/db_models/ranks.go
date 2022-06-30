package db_models

type Rank struct {
	GormModel          `json:"-"`
	Id                 uint   `json:"id"`
	NeededNFT          string `json:"neededNFT"`
	ReputationModifier int    `json:"reputationModifier"`
	MarketAccess       bool   `json:"marketAccess"`
}
