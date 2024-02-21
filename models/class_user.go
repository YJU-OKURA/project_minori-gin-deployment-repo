package models

type ClassUser struct {
	CID        uint   `gorm:"primaryKey"`
	UID        uint   `gorm:"primaryKey"`
	Nickname   string `gorm:"size:50;not null"`
	IsFavorite bool
	RoleID     int   `gorm:"not null"`
	Class      Class `gorm:"foreignKey:CID"`
	User       User  `gorm:"foreignKey:UID"`
}
