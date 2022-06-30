package db_models

type ThreadSection struct {
	GormModel   `json:"-"`
	Id          uint   `json:"id" gorm:"primaryKey"`
	Category    string `json:"category"`
	Parent      string `json:"parent"`
	Name        string `json:"name"`
	PublicPost  bool   `json:"publicPost"`
	MinimumRank int    `json:"minimumRank"`
}
