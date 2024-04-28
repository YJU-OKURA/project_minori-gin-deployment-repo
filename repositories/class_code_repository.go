package repositories

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type ClassCodeRepository interface {
	FindByCode(code string) (*models.ClassCode, error)
	SaveClassCode(classCode *models.ClassCode) error
}

// ClassCodeRepository はグループコードのリポジトリです。
type classCodeRepository struct {
	db *gorm.DB
}

// NewClassCodeRepository はClassCodeRepositoryを生成します。
func NewClassCodeRepository(db *gorm.DB) ClassCodeRepository {
	return &classCodeRepository{db: db}
}

// FindByCode は指定されたコードのグループコードを取得します。
func (r *classCodeRepository) FindByCode(code string) (*models.ClassCode, error) {
	var classCode models.ClassCode
	result := r.db.Where("code = ?", code).First(&classCode)
	if result.Error != nil {
		// レコードが見つからない場合、nilを返します。
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// それ以外のエラーの場合はエラーを返します。
		return nil, result.Error
	}
	return &classCode, nil
}

func (r *classCodeRepository) SaveClassCode(classCode *models.ClassCode) error {
	return r.db.Create(classCode).Error
}
