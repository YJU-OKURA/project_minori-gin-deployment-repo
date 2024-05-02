package repositories

import (
	"gorm.io/gorm"
)

// RoleRepository はロールのリポジトリです。
type RoleRepository interface {
	FindByRoleName(roleName string) (string, error) // 변경된 메서드 시그니처
}

// roleRepository はRoleRepositoryの実装です。
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository はRoleRepositoryを生成します。
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) FindByRoleName(roleName string) (string, error) {
	var role string
	result := r.db.Table("class_users").Select("role").Where("role = ?", roleName).Limit(1).Scan(&role)
	if result.Error != nil {
		return "", result.Error
	}
	return role, nil
}
