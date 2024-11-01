package repositories

import (
	"fmt"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type LINEAuthRepository interface {
	UpdateOrCreateUser(userInput dto.UserInput) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

type lineAuthRepository struct {
	db *gorm.DB
}

func NewLINEAuthRepository(db *gorm.DB) LINEAuthRepository {
	return &lineAuthRepository{db: db}
}

func (repo *lineAuthRepository) UpdateOrCreateUser(userInput dto.UserInput) (models.User, error) {
	var user models.User
	result := repo.db.Where("p_id = ?", fmt.Sprint(userInput.ID)).First(&user)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {

		pidPrefix := userInput.ID[:4]
		uniqueName := fmt.Sprintf("%s#%s", userInput.Name, pidPrefix)

		user = models.User{
			PID:   fmt.Sprint(userInput.ID),
			Name:  uniqueName,
			Image: userInput.Picture,
		}
		result = repo.db.Create(&user)
	}
	return user, result.Error
}

func (repo *lineAuthRepository) GetUserByID(id uint) (models.User, error) {
	var user models.User
	result := repo.db.First(&user, id)
	return user, result.Error
}
