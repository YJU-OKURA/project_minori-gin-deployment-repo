package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:50;not null"`
	Image     string    `gorm:"size:255;not null;"`
	PID       string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"not null;"`
}
