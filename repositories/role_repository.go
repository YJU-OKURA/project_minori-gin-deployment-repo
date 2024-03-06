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
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository はRoleRepositoryを生成します。
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// FindByRoleName は指定されたロール名のロールを取得します。
func (r *roleRepository) FindByRoleName(roleName string, role *models.Role) error {
	return r.db.Where("role = ?", roleName).First(role).Error
}
