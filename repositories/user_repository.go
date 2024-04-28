package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// UserRepository は
type UserRepository interface {
	GetApplyingClasses(userID uint) ([]models.ClassUser, error)
	UserExists(userID uint) (bool, error)
	FindByName(name string) ([]models.User, error)
}

// roleConnection　はRoleRepositoryの実装です。
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository はUserRepositoryを生成します。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetApplyingClasses はユーザーが申請中のクラスを取得します。
func (r *userRepository) GetApplyingClasses(userID uint) ([]models.ClassUser, error) {
	var classUsers []models.ClassUser
	err := r.db.Preload("Class").Preload("User").Where("uid = ? AND role_id = ?", userID, models.RoleApplicantID).Find(&classUsers).Error
	return classUsers, err
}

// UserExists はユーザーが存在するかを確認します。
func (r *userRepository) UserExists(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error
	return count > 0, err
}

func (r *userRepository) FindByName(name string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("name LIKE ?", "%"+name+"%").Find(&users).Error
	return users, err
}
