package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// ClassUserService はグループコードのサービスです。
type ClassUserService interface {
	AssignRole(uid uint, cid uint, roleName string) error
	GetRole(uid uint, cid uint) (int, error)
	UpdateUserName(uid uint, cid uint, newName string) error
	GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error)
}

// classUserServiceImpl はClassCodeServiceの実装です。
type classUserServiceImpl struct {
	roleRepo      repositories.RoleRepository
	classUserRepo repositories.ClassUserRepository
}

func NewClassUserService(classUserRepo repositories.ClassUserRepository, roleRepo repositories.RoleRepository) ClassUserService {
	return &classUserServiceImpl{
		classUserRepo: classUserRepo,
		roleRepo:      roleRepo,
	}
}

func (s *classUserServiceImpl) AssignRole(uid uint, cid uint, roleName string) error {
	var role models.Role
	err := s.roleRepo.FindByRoleName(roleName, &role)
	if err != nil {
		return err
	}

	return s.classUserRepo.UpdateUserRole(uid, cid, role.ID)
}

func (s *classUserServiceImpl) GetRole(uid uint, cid uint) (int, error) {
	roleID, err := s.classUserRepo.GetRole(uid, cid)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

func (s *classUserServiceImpl) UpdateUserName(uid uint, cid uint, newName string) error {
	return s.classUserRepo.UpdateUserName(uid, cid, newName)
}

func (s *classUserServiceImpl) GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClasses(uid)
}
