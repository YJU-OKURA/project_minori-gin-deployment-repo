package services

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/gin-gonic/gin"
)

// ClassCodeService はグループコードのサービスです。
type ClassCodeService interface {
	CheckSecretExists(c *gin.Context, code string) (bool, error)
	VerifyClassCode(code, secret string) (*models.ClassCode, error)
}

// classCodeServiceImpl はClassCodeServiceの実装です。
type classCodeServiceImpl struct {
	repo *repositories.ClassCodeRepository
}

// NewClassCodeService はClassCodeServiceを生成します。
func NewClassCodeService(repo *repositories.ClassCodeRepository) ClassCodeService {
	return &classCodeServiceImpl{repo: repo}
}

// CheckSecretExists は指定されたグループコードにシークレットがあるかどうかをチェックします。
func (s *classCodeServiceImpl) CheckSecretExists(c *gin.Context, code string) (bool, error) {
	classCode, err := s.repo.FindByCode(code)
	if err != nil {
		return false, err
	}
	if classCode == nil {
		return false, errors.New("class not found")
	}
	if classCode.Secret == nil || *classCode.Secret == "" {
		return false, nil // シークレットが存在しない場合はfalseを返す
	}

	return true, nil
}

// VerifyClassCode はグループコードと、該当する場合はそのシークレットを確認します。
func (s *classCodeServiceImpl) VerifyClassCode(code string, secret string) (*models.ClassCode, error) {
	classCode, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if classCode == nil || (classCode.Secret != nil && *classCode.Secret != secret) {
		return nil, errors.New("無効なグループコードまたはシークレットです。")
	}

	return classCode, nil
}
