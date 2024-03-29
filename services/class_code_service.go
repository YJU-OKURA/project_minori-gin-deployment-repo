package services

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/gin-gonic/gin"
)

const ErrClassNotFound = "class not found"
const ErrInvalidClassCodeOrSecret = "invalid class code or secret"

// ClassCodeService はグループコードのサービスです。
type ClassCodeService interface {
	CheckSecretExists(c *gin.Context, code string) (bool, error)
	VerifyClassCode(code, secret string) (bool, error)
}

// classCodeServiceImpl はClassCodeServiceの実装です。
type classCodeServiceImpl struct {
	repo *repositories.ClassCodeRepository
}

// NewClassCodeService はClassCodeServiceを生成します。
func NewClassCodeService(repo *repositories.ClassCodeRepository) ClassCodeService {
	return &classCodeServiceImpl{repo: repo}
}

// findClassCode は指定されたグループコードを取得します。
func (s *classCodeServiceImpl) findClassCode(code string) (*models.ClassCode, error) {
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
func (s *classCodeServiceImpl) CheckSecretExists(c *gin.Context, code string) (bool, error) {
	classCode, err := s.findClassCode(code)
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
	classCode, err := s.findClassCode(code)
	if err != nil {
		return false, err
	}

	if classCode.Secret == nil || *classCode.Secret != secret {
		return false, nil
	}

	return true, nil
}
