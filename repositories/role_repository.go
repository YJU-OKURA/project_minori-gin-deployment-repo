package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// RoleRepository はロールのリポジトリです。
type RoleRepository interface {
	FindByRoleName(roleName string, role *models.Role) error
}

// roleConnection　はRoleRepositoryの実装です。
type roleConnection struct {
	DB *gorm.DB
}

// NewRoleRepository はRoleRepositoryを生成します。
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleConnection{DB: db}
}

// FindByRoleName は指定されたロール名のロールを取得します。
func (r *roleConnection) FindByRoleName(roleName string, role *models.Role) error {
	result := r.DB.Where("role = ?", roleName).First(role)
	return result.Error
}
