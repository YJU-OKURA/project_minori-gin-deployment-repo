package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

type ClassService interface {
	CreateClass(request dto.CreateClassRequest) (uint, error)
	UpdateClassImage(classID uint, imageUrl string) error
}

type classServiceImpl struct {
	classRepo repositories.ClassRepository
}

func NewCreateClassService(classRepo repositories.ClassRepository) ClassService {
	return &classServiceImpl{
		classRepo: classRepo,
	}
}

func (s *classServiceImpl) CreateClass(request dto.CreateClassRequest) (uint, error) {
	class := models.Class{
		Name:        request.Name,
		Limitation:  request.Limitation,
		Description: request.Description,
	}

	classID, err := s.classRepo.Save(&class)
	if err != nil {
		return 0, err
	}

	return classID, nil
}

func (s *classServiceImpl) UpdateClassImage(classID uint, imageUrl string) error {
	return s.classRepo.UpdateClassImage(classID, imageUrl)
}
