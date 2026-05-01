package models

// Section represents a high-level category in the architecture (e.g. Core System)
type Section struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Icon        string    `gorm:"default:'📦'" json:"icon"`
	Description string    `json:"description"`
	Color       string    `gorm:"default:'#3b82f6'" json:"color"`
	SortOrder   int       `gorm:"default:0" json:"sort_order"`
	Features    []Feature `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE;" json:"features"` 
}

// Feature represents an individual item within a section
type Feature struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	SectionID    string `gorm:"index:idx_feature_unique,unique;not null" json:"section_id"`
	Title        string `gorm:"index:idx_feature_unique,unique;not null" json:"title"`
	Icon         string `gorm:"default:'📦'" json:"icon"`
	Status       string `gorm:"default:'live'" json:"status"`
	TagExtra     string `json:"tag_extra"`
	Subtitle     string `json:"subtitle"`
	HowItWorks   string `json:"how_it_works"`
	Approach     string `json:"approach"`
	Tech         string `gorm:"default:'[]'" json:"tech"` 
	Capabilities string `gorm:"default:'[]'" json:"capabilities"`
	Flow         string `gorm:"default:'[]'" json:"flow"`
	Impact       string `json:"impact"`
	SortOrder    int    `gorm:"default:0" json:"sort_order"`
}
