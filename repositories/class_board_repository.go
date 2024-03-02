package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// ClassBoardRepository インタフェース
type ClassBoardRepository interface {
	InsertClassBoard(b *models.ClassBoard) (*models.ClassBoard, error)
	FindByID(ID uint) (*models.ClassBoard, error)
	FindAll(cid uint) ([]models.ClassBoard, error)
	FindAnnounced(isAnnounced bool, cid uint) ([]models.ClassBoard, error)
	UpdateClassBoard(b *models.ClassBoard) error
	DeleteClassBoard(id uint) error
}

// classBoardConnection グループ掲示板リポジトリ
type classBoardConnection struct {
	DB *gorm.DB
}

// NewClassBoardRepository グループ掲示板リポジトリを生成
func NewClassBoardRepository(db *gorm.DB) ClassBoardRepository {
	return &classBoardConnection{DB: db}
}

// InsertClassBoard グループ掲示板を作成
func (db *classBoardConnection) InsertClassBoard(b *models.ClassBoard) (*models.ClassBoard, error) {
	result := db.DB.Create(b)
	return b, result.Error
}

// FindByID IDでグループ掲示板を取得
func (db *classBoardConnection) FindByID(ID uint) (*models.ClassBoard, error) {
	var classBoard models.ClassBoard
	result := db.DB.First(&classBoard, ID)
	return &classBoard, result.Error
}

// FindAnnounced 公開されたグループ掲示板を取得
func (db *classBoardConnection) FindAnnounced(isAnnounced bool, cid uint) ([]models.ClassBoard, error) {
	var classBoards []models.ClassBoard
	result := db.DB.Where("is_announced = ? AND cid = ?", isAnnounced, cid).Find(&classBoards)
	return classBoards, result.Error
}

// FindAll 全てのグループ掲示板を取得
func (db *classBoardConnection) FindAll(cid uint) ([]models.ClassBoard, error) {
	var classBoards []models.ClassBoard
	result := db.DB.Where("cid = ?", cid).Find(&classBoards)
	return classBoards, result.Error
}

// UpdateClassBoard グループ掲示板を更新
func (db *classBoardConnection) UpdateClassBoard(b *models.ClassBoard) error {
	result := db.DB.Save(b)
	return result.Error
}

// DeleteClassBoard グループ掲示板を削除
func (db *classBoardConnection) DeleteClassBoard(id uint) error {
	var classBoard models.ClassBoard
	result := db.DB.Delete(&classBoard, id)
	return result.Error
}
