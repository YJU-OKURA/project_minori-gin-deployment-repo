package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetApplyingClasses(userID uint) ([]models.ClassUser, error)
	UserExists(userID uint) (bool, error)
	FindByName(name string) ([]models.User, error)
	DeleteUser(userID uint) error
	FindByID(userID uint) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetApplyingClasses はユーザーが申請中のクラスを取得します。
func (r *userRepository) GetApplyingClasses(userID uint) ([]models.ClassUser, error) {
	var classUsers []models.ClassUser
	err := r.db.Preload("Class").Preload("User").Where("uid = ? AND role = ?", userID, "APPLICANT").Find(&classUsers).Error
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

func (r *userRepository) DeleteUser(userID uint) error {
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Delete(&models.User{}).Error
	return err
}

func (r *userRepository) FindByID(userID uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
