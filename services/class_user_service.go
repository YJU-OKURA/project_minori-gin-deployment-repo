package services

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"gorm.io/gorm"
)

// ClassUserService はグループコードのサービスです。
type ClassUserService interface {
	GetClassMembers(cid uint, roleID ...int) ([]dto.ClassMemberDTO, error)
	GetClassUserInfo(uid uint, cid uint) (dto.ClassMemberDTO, error)
	GetUserClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetRole(uid uint, cid uint) (int, error)
	GetFavoriteClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error)
	GetUserClassesByRole(uid uint, roleID int, page int, limit int) ([]dto.UserClassInfoDTO, error)
	AssignRole(uid uint, cid uint, roleID int) error
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

func (s *classUserServiceImpl) GetClassMembers(cid uint, roleID ...int) ([]dto.ClassMemberDTO, error) {
	return s.classUserRepo.GetClassMembers(cid, roleID...)
}

func (s *classUserServiceImpl) GetFavoriteClasses(uid uint, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetFavoriteClasses(uid, page, limit)
}

func (s *classUserServiceImpl) GetUserClassesByRole(uid uint, roleID int, page int, limit int) ([]dto.UserClassInfoDTO, error) {
	return s.classUserRepo.GetUserClassesByRole(uid, roleID, page, limit)
}

func (s *classUserServiceImpl) GetRole(uid uint, cid uint) (int, error) {
	roleID, err := s.classUserRepo.GetRole(uid, cid)
	if err != nil {
		return 0, err
	}

	return roleID, nil
}

func (s *classUserServiceImpl) AssignRole(uid uint, cid uint, roleID int) error {
	return s.classUserRepo.UpdateUserRole(uid, cid, roleID)
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
