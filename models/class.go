package models

type Class struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:30;not null"`
	Limitation  *int
	Description *string `gorm:"size:255"`
	Image       *string `gorm:"size:255"`
}
