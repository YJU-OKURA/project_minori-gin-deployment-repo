package repositories

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassUserRepository interface {
	UpdateUserRole(uid uint, cid uint, rid int) error
	GetRole(uid uint, cid uint) (int, error)
}

type classUserRepository struct {
	db *gorm.DB
}

func NewClassUserRepository(db *gorm.DB) ClassUserRepository {
	return &classUserRepository{db: db}
}

// UpdateUserRole はユーザーのロールを更新します。
func (r *classUserRepository) UpdateUserRole(uid uint, cid uint, rid int) error {
	var classUser models.ClassUser
	result := r.db.First(&classUser, "uid = ? AND cid = ?", uid, cid)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newUser := models.ClassUser{
			UID:    uid,
			CID:    cid,
			RoleID: rid,
		}
		return r.db.Create(&newUser).Error
	} else if result.Error != nil {
		return result.Error
	}

	return r.db.Model(&classUser).Update("role_id", rid).Error
}

// GetRole はユーザーのロールを取得します。
func (r *classUserRepository) GetRole(uid uint, cid uint) (int, error) {
	var classUser models.ClassUser
	result := r.db.First(&classUser, "uid = ? AND cid = ?", uid, cid)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, result.Error
	} else if result.Error != nil {
		return 0, result.Error
	}

	return classUser.RoleID, nil
}
