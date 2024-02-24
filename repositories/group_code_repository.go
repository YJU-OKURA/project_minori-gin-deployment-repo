package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// GroupCodeRepository はグループコードのリポジトリです。
type GroupCodeRepository struct {
	DB *gorm.DB
}

// NewGroupCodeRepository はGroupCodeRepositoryを生成します。
func NewGroupCodeRepository(db *gorm.DB) *GroupCodeRepository {
	return &GroupCodeRepository{DB: db}
}

// FindByCode は指定されたコードのグループコードを取得します。
func (r *GroupCodeRepository) FindByCode(code string) (models.GroupCode, error) {
	var groupCode models.GroupCode
	result := r.DB.Where("code = ?", code).First(&groupCode)
	return groupCode, result.Error
}
