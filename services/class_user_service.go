package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// ClassUserService はグループコードのサービスです。
type ClassUserService interface {
	AssignRole(uid uint, cid uint, roleName string) error
}

// classUserServiceImpl はClassCodeServiceの実装です。
type classUserServiceImpl struct {
	roleRepo      repositories.RoleRepository
	classUserRepo repositories.ClassUserRepository
}

func NewClassUserService(classUserRepo repositories.ClassUserRepository, roleRepo repositories.RoleRepository) ClassUserService {
	return &classUserServiceImpl{classUserRepo: classUserRepo, roleRepo: roleRepo}
}

func (s *classUserServiceImpl) AssignRole(uid uint, cid uint, roleName string) error {
	var role models.Role
	err := s.roleRepo.FindByRoleName(roleName, &role)
	if err != nil {
		return err
	}

	return s.classUserRepo.UpdateUserRole(uid, cid, role.ID)
}
