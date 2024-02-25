package models

type GroupCode struct {
	ID     string  `gorm:"primaryKey;size:255"`
	Code   string  `gorm:"size:10;not null"`
	Secret *string `gorm:"size:20"`
	CID    uint    `gorm:"column:cid;not null"` // Class ID
	UID    uint    `gorm:"column:uid;not null"` // User ID
	Class  Class   `gorm:"foreignKey:CID"`
	User   User    `gorm:"foreignKey:UID"`
}
