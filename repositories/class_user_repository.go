package repositories

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassUserRepository interface {
	GetClassMembers(cid uint, roleID ...int) ([]dto.ClassMemberDTO, error)
	GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error)
	GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetUserClassesByRole(uid uint, roleID int, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetRole(uid uint, cid uint) (int, error)
	UpdateUserRole(uid uint, cid uint, rid int) error
	UpdateUserName(uid uint, cid uint, newName string) error
	ToggleFavorite(uid uint, cid uint) error
	DeleteClassUser(uid uint, cid uint) error
	Save(classUser *models.ClassUser) error
}

type classUserRepository struct {
	db *gorm.DB
}

func NewClassUserRepository(db *gorm.DB) ClassUserRepository {
	return &classUserRepository{db: db}
}

// GetClassUserInfo はユーザーのクラスユーザー情報を取得します。
func (r *classUserRepository) GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error) {
	var classUser models.ClassUser
	err := r.db.Where("uid = ? AND cid = ?", uid, cid).First(&classUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ClassMemberDTO{}, errors.New(constants.UserNotFound)
		}
		return dto.ClassMemberDTO{}, err
	}
	return toClassMemberDTO(classUser), nil
}

// GetUserClasses はユーザーが所属しているクラスの情報を取得します。
func (r *classUserRepository) GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	var userClassesInfo []dto.UserClassInfoDTO
	offset := (page - 1) * limit

	err := r.db.Table("classes").
		Select("classes.id, classes.name, classes.limitation, classes.description, classes.image, class_users.is_favorite, class_users.role_id").
		Joins("INNER JOIN class_users ON classes.id = class_users.cid").
		Where("class_users.uid = ?", uid).
		Offset(offset).
		Limit(limit).
		Scan(&userClassesInfo).Error

	if err != nil {
		return nil, err
	}

	return userClassesInfo, nil
}

// GetClassMembers はクラスのメンバー情報を取得します。
func (r *classUserRepository) GetClassMembers(cid uint, roleID ...int) ([]dto.ClassMemberDTO, error) {
	var members []dto.ClassMemberDTO

	query := r.db.Table("class_users").
		Select("class_users.uid, class_users.nickname, class_users.role_id, users.image").
		Joins("join users on class_users.uid = users.id").
		Where("class_users.cid = ?", cid)

	if len(roleID) > 0 {
		query = query.Where("class_users.role_id = ?", roleID[0])
	}

	if err := query.Scan(&members).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []dto.ClassMemberDTO{}, nil
		}
		return nil, err
	}

	return members, nil
}

func (r *classUserRepository) GetUserClassesByRole(uid uint, roleID int, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	var userClassesInfo []dto.UserClassInfoDTO
	offset := (page - 1) * limit
	err := r.db.Table("classes").
		Select("classes.id, classes.name, classes.limitation, classes.description, classes.image, class_users.is_favorite, class_users.role_id").
		Joins("INNER JOIN class_users ON classes.id = class_users.cid").
		Where("class_users.uid = ? AND class_users.role_id = ?", uid, roleID).
		Offset(offset).
		Limit(limit).
		Scan(&userClassesInfo).Error

	if err != nil {
		return nil, err
	}

	return userClassesInfo, nil
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

func toClassMemberDTO(classUser models.ClassUser) dto.ClassMemberDTO {
	return dto.ClassMemberDTO{
		Uid:      classUser.UID,
		Nickname: classUser.Nickname,
		RoleId:   uint(classUser.RoleID),
		Image:    classUser.User.Image,
	}
}

func (r *classUserRepository) ToggleFavorite(uid uint, cid uint) error {
	var classUser models.ClassUser
	err := r.db.Model(&classUser).Where("uid = ? AND cid = ?", uid, cid).UpdateColumn("is_favorite", gorm.Expr("NOT is_favorite")).Error
	return err
}

func (r *classUserRepository) DeleteClassUser(uid uint, cid uint) error {
	return r.db.Where("uid = ? AND cid = ?", uid, cid).Delete(&models.ClassUser{}).Error
}

func (r *classUserRepository) Save(classUser *models.ClassUser) error {
	return r.db.Create(classUser).Error
}
