package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// ClassBoardRepository インタフェース
type ClassBoardRepository interface {
	InsertClassBoard(b *models.ClassBoard) (*models.ClassBoard, error)
	FindByID(id uint) (*models.ClassBoard, error)
	FindAll(cid uint) ([]models.ClassBoard, error)
	FindAnnounced(isAnnounced bool, cid uint) ([]models.ClassBoard, error)
	UpdateClassBoard(b *models.ClassBoard) error
	DeleteClassBoard(id uint) error
}

// classBoardConnection グループ掲示板リポジトリ
type classBoardRepository struct {
	db *gorm.DB
}

// NewClassBoardRepository グループ掲示板リポジトリを生成
func NewClassBoardRepository(db *gorm.DB) ClassBoardRepository {
	return &classBoardRepository{db: db}
}

// InsertClassBoard グループ掲示板を作成
func (repo *classBoardRepository) InsertClassBoard(b *models.ClassBoard) (*models.ClassBoard, error) {
	result := repo.db.Create(b)
	return b, result.Error
}

// FindByID IDでグループ掲示板を取得
func (repo *classBoardRepository) FindByID(id uint) (*models.ClassBoard, error) {
	var classBoard models.ClassBoard
	err := repo.db.First(&classBoard, id).Error
	return &classBoard, err
}

// FindAll 全てのグループ掲示板を取得
func (repo *classBoardRepository) FindAll(cid uint) ([]models.ClassBoard, error) {
	var classBoards []models.ClassBoard
	err := repo.db.Where("cid = ?", cid).Find(&classBoards).Error
	return classBoards, err
}

// FindAnnounced 公開されたグループ掲示板を取得
func (repo *classBoardRepository) FindAnnounced(isAnnounced bool, cid uint) ([]models.ClassBoard, error) {
	var classBoards []models.ClassBoard
	err := repo.db.Where("is_announced = ? AND cid = ?", isAnnounced, cid).Find(&classBoards).Error
	return classBoards, err
}

// UpdateClassBoard グループ掲示板を更新
func (repo *classBoardRepository) UpdateClassBoard(b *models.ClassBoard) error {
	return repo.db.Save(b).Error
}

// DeleteClassBoard グループ掲示板を削除
func (repo *classBoardRepository) DeleteClassBoard(id uint) error {
	return repo.db.Delete(&models.ClassBoard{}, id).Error
}
