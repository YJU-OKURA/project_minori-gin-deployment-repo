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
type classScheduleRepository struct {
	db *gorm.DB
}

// NewClassScheduleRepository クラススケジュールリポジトリを生成
func NewClassScheduleRepository(db *gorm.DB) ClassScheduleRepository {
	return &classScheduleRepository{db: db}
}

// GetClassScheduleByID クラススケジュールを取得
func (repo *classScheduleRepository) GetClassScheduleByID(id uint) (*models.ClassSchedule, error) {
	var classSchedule models.ClassSchedule
	err := repo.db.First(&classSchedule, id).Error
	return &classSchedule, err
}

// GetAllClassSchedules 全てのクラススケジュールを取得
func (repo *classScheduleRepository) GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	err := repo.db.Where("cid = ?", cid).Find(&classSchedules).Error
	return classSchedules, err
}

// CreateClassSchedule 新しいクラススケジュールを作成
func (repo *classScheduleRepository) CreateClassSchedule(classSchedule *models.ClassSchedule) error {
	return repo.db.Create(classSchedule).Error
}

// UpdateClassSchedule クラススケジュールを更新
func (repo *classScheduleRepository) UpdateClassSchedule(classSchedule *models.ClassSchedule) error {
	return repo.db.Save(classSchedule).Error
}

// DeleteClassSchedule クラススケジュールを削除
func (repo *classScheduleRepository) DeleteClassSchedule(id uint) error {
	return repo.db.Delete(&models.ClassSchedule{}, id).Error
}

// FindLiveClassSchedules ライブ中のクラススケジュールを取得
func (repo *classScheduleRepository) FindLiveClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	err := repo.db.Where("cid = ? AND is_live = true AND end_time > NOW()", cid).Find(&classSchedules).Error
	return classSchedules, err
}

// FindClassSchedulesByDate 日付でクラススケジュールを取得
func (repo *classScheduleRepository) FindClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error) {
	var classSchedules []models.ClassSchedule
	err := repo.db.Where("cid = ? AND DATE(start_time) = ?", cid, date).Find(&classSchedules).Error
	return classSchedules, err
}
