package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// ClassUserService はグループコードのサービスです。
type ClassUserService interface {
	GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error)
	GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error)
	GetRole(uid uint, cid uint) (int, error)
	GetFavoriteClasses(uid uint) ([]dto.UserClassInfoDTO, error)
	GetUserClassesByRole(uid uint, roleID int) ([]dto.UserClassInfoDTO, error)
	AssignRole(uid uint, cid uint, roleName string) error
	UpdateUserName(uid uint, cid uint, newName string) error
	GetClassMembers(cid uint) ([]dto.ClassMemberDTO, error)
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

func (s *classUserServiceImpl) GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error) {
	return s.classUserRepo.GetClassUserInfo(uid, cid)
}

func (s *classUserServiceImpl) GetUserClasses(uid uint) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClasses(uid)
}

func (s *classUserServiceImpl) GetClassMembers(cid uint) ([]dto.ClassMemberDTO, error) {
	return s.classUserRepo.GetClassMembers(cid)
}

func (s *classUserServiceImpl) GetFavoriteClasses(uid uint) ([]dto.UserClassInfoDTO, error) {
	classes, err := s.classUserRepo.GetUserClasses(uid)
	if err != nil {
		return nil, err
	}

	if len(classes) == 0 {
		return nil, ErrNotFound
	}

	var favoriteClasses []dto.UserClassInfoDTO
	for _, class := range classes {
		if class.IsFavorite {
			favoriteClasses = append(favoriteClasses, class)
		}
	}

	return favoriteClasses, nil
}

func (s *classUserServiceImpl) GetUserClassesByRole(uid uint, roleID int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClassesByRole(uid, roleID)
}

func (s *classUserServiceImpl) GetRole(uid uint, cid uint) (int, error) {
	roleID, err := s.classUserRepo.GetRole(uid, cid)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

func (s *classUserServiceImpl) AssignRole(uid uint, cid uint, roleName string) error {
	var role models.Role
	err := s.roleRepo.FindByRoleName(roleName, &role)
	if err != nil {
		return err
	}

	return s.classUserRepo.UpdateUserRole(uid, cid, role.ID)
}

func (s *classUserServiceImpl) UpdateUserName(uid uint, cid uint, newName string) error {
	return s.classUserRepo.UpdateUserName(uid, cid, newName)
}
