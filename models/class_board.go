package models

import "time"

type ClassBoard struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:255;not null"`
	Content     string    `gorm:"type:text;not null"`
	Image       string    `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"not null;"`
	UpdatedAt   time.Time `gorm:"not null;"`
	IsAnnounced bool      `gorm:"not null;default:false"`
	CID         uint      `gorm:"column:cid;not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UID         uint      `gorm:"column:uid;not null"` // User ID
	Class       Class     `gorm:"foreignKey:CID;"`
	User        User      `gorm:"foreignKey:UID"`
}
