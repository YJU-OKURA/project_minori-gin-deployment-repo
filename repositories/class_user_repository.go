package repositories

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassUserRepository interface {
	UpdateUserRole(uid uint, cid uint, rid int) error
}

type classUserConnection struct {
	DB *gorm.DB
}

func NewClassUserRepository(db *gorm.DB) ClassUserRepository {
	return &classUserConnection{DB: db}
}

// UpdateUserRole はユーザーのロールを更新します。
func (r *classUserConnection) UpdateUserRole(uid uint, cid uint, rid int) error {
	var classUser models.ClassUser
	result := r.DB.First(&classUser, "uid = ? AND cid = ?", uid, cid)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newUser := models.ClassUser{
			UID:    uid,
			CID:    cid,
			RoleID: rid,
		}
		return r.DB.Create(&newUser).Error
	} else if result.Error != nil {
		return result.Error
	}

	return r.DB.Model(&classUser).Update("role_id", rid).Error
}
