package repositories

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	GetByID(classID uint) (*models.Class, error)
	Create(class *models.Class) error
	Save(class *models.Class) (uint, error)
	UpdateClassImage(classID uint, imageUrl string) error
	Update(class *models.Class) error
	Delete(classID uint) error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) GetByID(classID uint) (*models.Class, error) {
	var class models.Class
	result := r.db.First(&class, classID)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &class, nil
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

func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(classID uint) error {
	return r.db.Delete(&models.Class{}, classID).Error
}
