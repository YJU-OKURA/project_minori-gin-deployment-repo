package models

type Role struct {
	ID   int    `gorm:"primaryKey"`
	Role string `gorm:"type:enum('USER', 'ADMIN', 'ASSISTANT', 'APPLICANT', 'BLACKLIST');not null"`
}
