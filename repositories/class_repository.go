package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	Create(class *models.Class) error
	Save(class *models.Class) (uint, error)
	UpdateClassImage(classID uint, imageUrl string) error
}

type classRepository struct {
	db *gorm.DB
}

func NewCreateClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) Save(class *models.Class) (uint, error) {
	if err := r.db.Create(&class).Error; err != nil {
		return 0, err
	}
	return class.ID, nil
}

func (r *classRepository) UpdateClassImage(classID uint, imageUrl string) error {
	return r.db.Model(&models.Class{}).Where("id = ?", classID).Update("image", imageUrl).Error
}
