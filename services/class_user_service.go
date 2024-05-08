package services

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"gorm.io/gorm"
)

// ClassUserService はグループコードのサービスです。
type ClassUserService interface {
	GetClassMembers(cid uint, roleNames ...string) ([]dto.ClassMemberDTO, error)
	GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error)
	GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetRole(uid uint, cid uint) (string, error)
	GetFavoriteClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetUserClassesByRole(uid uint, roleName string, page int, limit int) ([]dto.UserClassInfoDTO, error)
	AssignRole(uid uint, cid uint, roleName string) error
	UpdateUserName(uid uint, cid uint, newName string) error
	ToggleFavorite(uid uint, cid uint) error
	RemoveUserFromClass(uid uint, cid uint) error
	SearchUserClassesByName(uid uint, name string) ([]dto.UserClassInfoDTO, error)
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

func (s *classUserServiceImpl) GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClasses(uid, page, limit)
}

func (s *classUserServiceImpl) GetClassMembers(cid uint, roleNames ...string) ([]dto.ClassMemberDTO, error) {
	if len(roleNames) > 0 {
		return s.classUserRepo.GetClassMembers(cid, roleNames[0])
	}
	return s.classUserRepo.GetClassMembers(cid)
}

func (s *classUserServiceImpl) GetFavoriteClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetFavoriteClasses(uid, page, limit)
}

func (s *classUserServiceImpl) GetUserClassesByRole(uid uint, roleName string, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClassesByRole(uid, roleName, page, limit)
}

func (s *classUserServiceImpl) GetRole(uid uint, cid uint) (string, error) {
	roleName, err := s.classUserRepo.GetRole(uid, cid)
	if err != nil {
		return "", err
	}
	return roleName, nil
}

func (s *classUserServiceImpl) AssignRole(uid uint, cid uint, roleName string) error {
	exists, err := s.classUserRepo.RoleExists(uid, cid)
	if err != nil {
		return err
	}
	if exists {
		return s.classUserRepo.UpdateUserRole(uid, cid, roleName)
	} else {
		return s.classUserRepo.CreateUserRole(uid, cid, roleName)
	}
}

func (s *classUserServiceImpl) UpdateUserName(uid uint, cid uint, newName string) error {
	return s.classUserRepo.UpdateUserName(uid, cid, newName)
}

func (s *classUserServiceImpl) ToggleFavorite(uid uint, cid uint) error {
	err := s.classUserRepo.ToggleFavorite(uid, cid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}

		return err
	}
	return nil
}

func (s *classUserServiceImpl) RemoveUserFromClass(uid uint, cid uint) error {
	return s.classUserRepo.DeleteClassUser(uid, cid)
}

func (s *classUserServiceImpl) SearchUserClassesByName(uid uint, name string) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.SearchUserClassesByName(uid, name)
}
