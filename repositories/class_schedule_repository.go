package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// ClassScheduleRepository インタフェース
type ClassScheduleRepository interface {
	GetClassScheduleByID(id uint) (*models.ClassSchedule, error)
	GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error)
	CreateClassSchedule(classSchedule *models.ClassSchedule) error
	UpdateClassSchedule(classSchedule *models.ClassSchedule) error
	DeleteClassSchedule(id uint) error
	FindLiveClassSchedules(cid uint) ([]models.ClassSchedule, error)
	FindClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error)
}

// classScheduleConnection クラススケジュールリポジトリ
type classScheduleConnection struct {
	DB *gorm.DB
}

// NewClassScheduleRepository クラススケジュールリポジトリを生成
func NewClassScheduleRepository(db *gorm.DB) ClassScheduleRepository {
	return &classScheduleConnection{DB: db}
}

// GetClassScheduleByID クラススケジュールを取得
func (db *classScheduleConnection) GetClassScheduleByID(id uint) (*models.ClassSchedule, error) {
	var classSchedule models.ClassSchedule
	err := db.DB.First(&classSchedule, id).Error
	return &classSchedule, err
}

// GetAllClassSchedules 全てのクラススケジュールを取得
func (db *classScheduleConnection) GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	result := db.DB.Where("cid = ?", cid).Find(&classSchedules)
	return classSchedules, result.Error
}

// CreateClassSchedule 新しいクラススケジュールを作成
func (db *classScheduleConnection) CreateClassSchedule(classSchedule *models.ClassSchedule) error {
	return db.DB.Create(classSchedule).Error
}

// UpdateClassSchedule クラススケジュールを更新
func (db *classScheduleConnection) UpdateClassSchedule(classSchedule *models.ClassSchedule) error {
	return db.DB.Save(classSchedule).Error
}

// DeleteClassSchedule クラススケジュールを削除
func (db *classScheduleConnection) DeleteClassSchedule(id uint) error {
	return db.DB.Delete(&models.ClassSchedule{}, id).Error
}

// FindLiveClassSchedules ライブ中のクラススケジュールを取得
func (db *classScheduleConnection) FindLiveClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	result := db.DB.Where("cid = ? AND is_live = true AND end_time > NOW()", cid).Find(&classSchedules)
	return classSchedules, result.Error
}

// FindClassSchedulesByDate 日付でクラススケジュールを取得
func (db *classScheduleConnection) FindClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	result := db.DB.Where("cid = ? AND DATE(start_time) = ?", cid, date).Find(&classSchedules)
	return classSchedules, result.Error
}
