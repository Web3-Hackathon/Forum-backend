package models

// SectionsResponseModel is used when fetching all the
// SectionInfo columns
type SectionsResponseModel struct {
	BaseResponseModel
	Sections map[string]map[string][]map[string]uint `json:"sections"`
}
