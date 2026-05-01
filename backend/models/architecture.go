package models

// Section represents a high-level category in the architecture (e.g. Core System)
type Section struct {
	ID          string    `gorm:"primaryKey"`
	Title       string    `gorm:"not null"`
	Icon        string    `gorm:"default:'📦'"`
	Description string
	Color       string    `gorm:"default:'#3b82f6'"`
	SortOrder   int       `gorm:"default:0"`
	Features    []Feature `gorm:"foreignKey:SectionID;constraint:OnDelete:CASCADE;"` // Preloads all associated features
}

// Feature represents an individual item within a section
type Feature struct {
	ID           uint   `gorm:"primaryKey"`
	SectionID    string `gorm:"index:idx_feature_unique,unique;not null"`
	Title        string `gorm:"index:idx_feature_unique,unique;not null"`
	Icon         string `gorm:"default:'📦'"`
	Status       string `gorm:"default:'live'"`
	TagExtra     string
	Subtitle     string
	HowItWorks   string
	Approach     string
	Tech         string `gorm:"default:'[]'"` // Stored as JSON string arrays
	Capabilities string `gorm:"default:'[]'"` // Stored as JSON string arrays
	Flow         string `gorm:"default:'[]'"` // Stored as JSON string arrays
	Impact       string
	SortOrder    int    `gorm:"default:0"`
}
