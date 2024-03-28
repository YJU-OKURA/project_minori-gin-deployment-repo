package repositories

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassUserRepository interface {
	GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error)
	GetRole(uid uint, cid uint) (int, error)
	UpdateUserRole(uid uint, cid uint, rid int) error
	UpdateUserName(uid uint, cid uint, newName string) error
	GetClassMembers(cid uint) ([]dto.ClassMemberDTO, error)
}

type classUserRepository struct {
	db *gorm.DB
}

func NewClassUserRepository(db *gorm.DB) ClassUserRepository {
	return &classUserRepository{db: db}
}

// GetUserClasses はユーザーが所属しているクラスの情報を取得します。
func (r *classUserRepository) GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error) {
	var userClassesInfo []dto.UserClassInfoDTO
	err := r.db.Table("classes").
		Select("classes.id, classes.name, classes.limitation, classes.description, classes.image, class_users.is_favorite, class_users.role_id").
		Joins("INNER JOIN class_users ON classes.id = class_users.cid").
		Where("class_users.uid = ?", uid).
		Scan(&userClassesInfo).Error

	return userClassesInfo, err
}

// GetClassMembers はクラスのメンバー情報を取得します。
func (r *classUserRepository) GetClassMembers(cid uint) ([]dto.ClassMemberDTO, error) {
	var members []dto.ClassMemberDTO
	err := r.db.Table("class_users").
		Select("class_users.uid, class_users.nickname, class_users.role_id, users.image").
		Joins("join users on class_users.uid = users.id").
		Where("class_users.cid = ?", cid).
		Scan(&members).Error
	return members, err
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

// UpdateUserName はユーザーの名前を更新します。
func (r *classUserRepository) UpdateUserName(uid uint, cid uint, newName string) error {
	var classUser models.ClassUser
	result := r.db.First(&classUser, "uid = ? AND cid = ?", uid, cid)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	} else if result.Error != nil {
		return result.Error
	}

	return r.db.Model(&classUser).Update("nickname", newName).Error
}
