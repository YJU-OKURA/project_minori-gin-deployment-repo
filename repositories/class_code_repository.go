package repositories

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// ClassCodeRepository はグループコードのリポジトリです。
type ClassCodeRepository struct {
	DB *gorm.DB
}

// NewClassCodeRepository はClassCodeRepositoryを生成します。
func NewClassCodeRepository(db *gorm.DB) *ClassCodeRepository {
	return &ClassCodeRepository{DB: db}
}

// FindByCode は指定されたコードのグループコードを取得します。
func (r *ClassCodeRepository) FindByCode(code string) (*models.ClassCode, error) {
	var classCode models.ClassCode
	result := r.DB.Where("code = ?", code).First(&classCode)
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
