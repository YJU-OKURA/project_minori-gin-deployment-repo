package models

type GroupCode struct {
	ID     string `gorm:"primaryKey;size:255"`
	Code   string `gorm:"size:10;not null"`
	Secret string `gorm:"size:20"`
	CID    uint   `gorm:"not null"` // Class ID
	UID    uint   `gorm:"not null"` // User ID
	Class  Class  `gorm:"foreignKey:CID"`
	User   User   `gorm:"foreignKey:UID"`
}
