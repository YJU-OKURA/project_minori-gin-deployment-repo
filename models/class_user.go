package models

type ClassUser struct {
	CID        uint   `gorm:"column:cid;primaryKey"`
	UID        uint   `gorm:"column:uid;primaryKey"`
	Nickname   string `gorm:"size:50;not null"`
	IsFavorite bool   `gorm:"not null;default:false"`
	Role       string `gorm:"type:Role;not null"`
	Class      Class  `gorm:"foreignKey:CID;constraint:OnDelete:CASCADE"`
	User       User   `gorm:"foreignKey:UID"`
}
