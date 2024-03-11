package repositories

import (
	"fmt"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type GoogleAuthRepository interface {
	UpdateOrCreateUser(userInput dto.UserInput) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

type googleAuthRepository struct {
	db *gorm.DB
}

func NewGoogleAuthRepository(db *gorm.DB) GoogleAuthRepository {
	return &googleAuthRepository{db: db}
}

func (repo *googleAuthRepository) UpdateOrCreateUser(userInput dto.UserInput) (models.User, error) {
	var user models.User
	result := repo.db.Where("p_id = ?", fmt.Sprint(userInput.ID)).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user = models.User{
				PID:   fmt.Sprint(userInput.ID),
				Name:  userInput.Name,
				Image: userInput.Picture,
			}
			result = repo.db.Create(&user)
		}
	}
	return user, result.Error
}

func (repo *googleAuthRepository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	result := repo.db.First(&user, id)
	return user, result.Error
}
