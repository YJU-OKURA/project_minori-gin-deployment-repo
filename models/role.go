package models

const RoleApplicantID = 4

type Role struct {
	ID   int    `gorm:"primaryKey"`
	Role string `gorm:"type:enum('USER', 'ADMIN', 'ASSISTANT', 'APPLICANT', 'BLACKLIST');not null"`
}
