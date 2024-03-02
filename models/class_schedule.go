package models

import "time"

type ClassSchedule struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	StartedAt time.Time `gorm:"not null"`
	EndedAt   time.Time `gorm:"not null"`
	CID       uint      `gorm:"column:cid;not null"` // Class ID
	IsLive    bool      `gorm:"not null;default:false"`
	Class     Class     `gorm:"foreignKey:CID"`
}
