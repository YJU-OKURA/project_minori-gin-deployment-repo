package repositories

import (
	"fmt"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type LINEAuthRepository interface {
	UpdateOrCreateUser(userInput dto.LineUserInput) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

type lineAuthRepository struct {
	db *gorm.DB
}

func NewLINEAuthRepository(db *gorm.DB) LINEAuthRepository {
	return &lineAuthRepository{db: db}
}

func (repo *lineAuthRepository) UpdateOrCreateUser(userInput dto.LineUserInput) (models.User, error) {
	var user models.User
	result := repo.db.Where("p_id = ?", fmt.Sprint(userInput.UserID)).First(&user)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {

		pidPrefix := userInput.UserID[:4]
		uniqueName := fmt.Sprintf("%s#%s", userInput.DisplayName, pidPrefix)

		user = models.User{
			PID:   fmt.Sprint(userInput.UserID),
			Name:  uniqueName,
			Image: userInput.PictureURL,
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
