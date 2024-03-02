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
	Repo          *repositories.ClassCodeRepository
	RoleRepo      repositories.RoleRepository
	ClassUserRepo repositories.ClassUserRepository
}

func NewClassUserService(classUserRepo repositories.ClassUserRepository, roleRepo repositories.RoleRepository) ClassUserService {
	return &classUserServiceImpl{ClassUserRepo: classUserRepo, RoleRepo: roleRepo}
}

func (service *classUserServiceImpl) AssignRole(uid uint, cid uint, roleName string) error {
	var role models.Role
	err := service.RoleRepo.FindByRoleName(roleName, &role)
	if err != nil {
		return err
	}

	return service.ClassUserRepo.UpdateUserRole(uid, cid, role.ID)
}
