package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// GroupBoardRepository インタフェース
type GroupBoardRepository interface {
	InsertGroupBoard(b *models.GroupBoard) (*models.GroupBoard, error)
	FindByID(ID uint) (*models.GroupBoard, error)
	FindAll() ([]models.GroupBoard, error)
	FindAnnounced(isAnnounced bool) ([]models.GroupBoard, error)
	UpdateGroupBoard(b *models.GroupBoard) error
	DeleteGroupBoard(id uint) error
}

// groupBoardConnection グループ掲示板リポジトリ
type groupBoardConnection struct {
	DB *gorm.DB
}

// NewGroupBoardRepository グループ掲示板リポジトリを生成
func NewGroupBoardRepository(db *gorm.DB) GroupBoardRepository {
	return &groupBoardConnection{DB: db}
}

// InsertGroupBoard グループ掲示板を作成
func (db *groupBoardConnection) InsertGroupBoard(b *models.GroupBoard) (*models.GroupBoard, error) {
	result := db.DB.Create(b)
	return b, result.Error
}

// FindByID IDでグループ掲示板を取得
func (db *groupBoardConnection) FindByID(ID uint) (*models.GroupBoard, error) {
	var groupBoard models.GroupBoard
	result := db.DB.First(&groupBoard, ID)
	return &groupBoard, result.Error
}

// FindAnnounced 公開されたグループ掲示板を取得
func (db *groupBoardConnection) FindAnnounced(isAnnounced bool) ([]models.GroupBoard, error) {
	var groupBoards []models.GroupBoard
	result := db.DB.Where("is_announced = ?", isAnnounced).Find(&groupBoards)
	return groupBoards, result.Error
}

// FindAll 全てのグループ掲示板を取得
func (db *groupBoardConnection) FindAll() ([]models.GroupBoard, error) {
	var groupBoards []models.GroupBoard
	result := db.DB.Find(&groupBoards)
	return groupBoards, result.Error
}

// UpdateGroupBoard グループ掲示板を更新
func (db *groupBoardConnection) UpdateGroupBoard(b *models.GroupBoard) error {
	result := db.DB.Save(b)
	return result.Error
}

// DeleteGroupBoard グループ掲示板を削除
func (db *groupBoardConnection) DeleteGroupBoard(id uint) error {
	var groupBoard models.GroupBoard
	result := db.DB.Delete(&groupBoard, id)
	return result.Error
}
