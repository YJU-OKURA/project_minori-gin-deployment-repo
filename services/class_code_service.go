package services

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

const ErrClassNotFound = "class not found"

// ClassCodeService はグループコードのサービスです。
type ClassCodeService interface {
	CheckSecretExists(code string) (bool, error)
	VerifyClassCode(code, secret string) (bool, error)
	FindClassCode(code string) (*models.ClassCode, error)
}

// classCodeServiceImpl はClassCodeServiceの実装です。
type classCodeServiceImpl struct {
	repo repositories.ClassCodeRepository
}

// NewClassCodeService はClassCodeServiceを生成します。
func NewClassCodeService(repo repositories.ClassCodeRepository) ClassCodeService {
	return &classCodeServiceImpl{repo: repo}
}

// FindClassCode findClassCode は指定されたグループコードを取得します。
func (s *classCodeServiceImpl) FindClassCode(code string) (*models.ClassCode, error) {
	classCode, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}
	if classCode == nil {
		return nil, errors.New(ErrClassNotFound)
	}
	return classCode, nil
}

// CheckSecretExists は指定されたグループコードにシークレットがあるかどうかをチェックします。
func (s *classCodeServiceImpl) CheckSecretExists(code string) (bool, error) {
	classCode, err := s.FindClassCode(code)
	if err != nil {
		return false, err
	}
	if classCode == nil {
		return false, errors.New(ErrClassNotFound)
	}
	if classCode.Secret == nil || *classCode.Secret == "" {
		return false, nil // シークレットが存在しない場合はfalseを返す
	}

	return true, nil
}

// VerifyClassCode はグループコードと、該当する場合はそのシークレットを確認します。
func (s *classCodeServiceImpl) VerifyClassCode(code string, secret string) (bool, error) {
	classCode, err := s.FindClassCode(code)
	if err != nil {
		return false, err
	}

	if classCode.Secret == nil || *classCode.Secret != secret {
		return false, nil
	}

	return true, nil
}
