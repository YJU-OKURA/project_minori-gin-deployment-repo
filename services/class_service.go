package services

import (
	"errors"
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

type ClassService interface {
	GetClass(classID uint) (*models.Class, error)
	CreateClass(request dto.CreateClassRequest) (uint, error)
	UpdateClassImage(classID uint, imageUrl string) error
	DeleteClass(classID uint, userID uint) error
}

type classServiceImpl struct {
	classRepo     repositories.ClassRepository
	classUserRepo repositories.ClassUserRepository
}

func NewCreateClassService(classRepo repositories.ClassRepository, classUserRepo repositories.ClassUserRepository) ClassService {
	return &classServiceImpl{
		classRepo:     classRepo,
		classUserRepo: classUserRepo,
	}
}

func (s *classServiceImpl) GetClass(classID uint) (*models.Class, error) {
	return s.classRepo.GetByID(classID)
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

func (s *classServiceImpl) DeleteClass(classID uint, userID uint) error {
	roleID, err := s.classUserRepo.GetRole(userID, classID)
	if err != nil {
		return err
	}

	if roleID != 2 {
		return errors.New(fmt.Sprintf("unauthorized access: roleID %d", roleID))
	}

	return s.classRepo.Delete(classID)
}
