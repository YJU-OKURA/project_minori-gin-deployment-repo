package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

type ClassService interface {
	GetClass(classID uint) (*models.Class, error)
	CreateClass(request dto.CreateClassRequest) (uint, error)
	UpdateClassImage(classID uint, imageUrl string) error
	UpdateClass(classID uint, userID uint, request dto.UpdateClassRequest) error
	DeleteClass(classID uint, userID uint) error
	GenerateClassCode() (string, error)
}

type classServiceImpl struct {
	classRepo     repositories.ClassRepository
	classUserRepo repositories.ClassUserRepository
	classCodeRepo repositories.ClassCodeRepository
	userRepo      repositories.UserRepository
}

func NewCreateClassService(
	classRepo repositories.ClassRepository,
	classUserRepo repositories.ClassUserRepository,
	classCodeRepo repositories.ClassCodeRepository,
	userRepo repositories.UserRepository,
) ClassService {
	return &classServiceImpl{
		classRepo:     classRepo,
		classUserRepo: classUserRepo,
		classCodeRepo: classCodeRepo,
		userRepo:      userRepo,
	}
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (s *classServiceImpl) GetClass(classID uint) (*models.Class, error) {
	return s.classRepo.GetByID(classID)
}

func (s *classServiceImpl) CreateClass(request dto.CreateClassRequest) (uint, error) {

	var user *models.User
	var err error
	if user, err = s.userRepo.FindByID(request.UID); err != nil {
		return 0, err
	}

	class := models.Class{
		Name:        request.Name,
		Limitation:  request.Limitation,
		Description: request.Description,
		UID:         request.UID,
	}

	classID, err := s.classRepo.Save(&class)
	if err != nil {
		return 0, err
	}

	classUser := models.ClassUser{
		CID:        classID,
		UID:        request.UID,
		Nickname:   user.Name,
		IsFavorite: false,
		Role:       "ADMIN",
	}
	err = s.classUserRepo.Save(&classUser)
	if err != nil {
		return 0, err
	}

	code, err := s.GenerateClassCode()
	if err != nil {
		return 0, err
	}
	classCode := models.ClassCode{
		Code:   code,
		CID:    classID,
		UID:    request.UID,
		Secret: request.Secret,
	}
	if err := s.classCodeRepo.SaveClassCode(&classCode); err != nil {
		return 0, err
	}

	return classID, nil
}

func (s *classServiceImpl) UpdateClassImage(classID uint, imageUrl string) error {
	return s.classRepo.UpdateClassImage(classID, imageUrl)
}

func (s *classServiceImpl) UpdateClass(classID uint, userID uint, request dto.UpdateClassRequest) error {
	isAdmin, err := s.IsAdmin(userID, classID)
	if err != nil || !isAdmin {
		return errors.New("unauthorized: user is not an admin")
	}

	class, err := s.GetClass(classID)
	if err != nil {
		return err
	}

	if request.Name != "" {
		class.Name = request.Name
	}
	if request.Limitation != nil {
		class.Limitation = request.Limitation
	}
	if request.Description != nil {
		class.Description = request.Description
	}

	return s.classRepo.Update(class)
}

func (s *classServiceImpl) IsAdmin(userID uint, classID uint) (bool, error) {
	role, err := s.classUserRepo.GetRole(userID, classID)
	if err != nil {
		return false, err
	}
	return role == "ADMIN", nil
}

func (s *classServiceImpl) DeleteClass(classID uint, userID uint) error {
	role, err := s.classUserRepo.GetRole(userID, classID)
	if err != nil {
		return err
	}

	if role != "ADMIN" {
		return errors.New(fmt.Sprintf("unauthorized access: role %s", role))
	}

	return s.classRepo.Delete(classID)
}

func (s *classServiceImpl) GenerateClassCode() (string, error) {
	for {
		code := make([]byte, 6)
		for i := range code {
			code[i] = letters[rand.Intn(len(letters))]
		}
		existingCode, err := s.classCodeRepo.FindByCode(string(code))
		if err != nil {
			return "", err
		}
		if existingCode == nil {
			return string(code), nil
		}
	}
}
