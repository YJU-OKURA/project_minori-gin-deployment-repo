package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// AttendanceRepository インタフェース
type AttendanceRepository interface {
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceByUIDAndCID(uid uint, cid uint) (*models.Attendance, error)
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAttendanceByID(id string) (*models.Attendance, error)
	UpdateAttendance(attendance *models.Attendance) error
	DeleteAttendance(id string) error
}

// attendanceConnection グループ掲示板リポジトリ
type attendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository グループ掲示板リポジトリを生成
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

// CreateAttendance 出席情報を作成
func (repo *attendanceRepository) CreateAttendance(attendance *models.Attendance) error {
	return repo.db.Create(attendance).Error
}

// GetAttendanceByUIDAndCID UIDとCIDによって出席情報を取得
func (repo *attendanceRepository) GetAttendanceByUIDAndCID(uid uint, cid uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := repo.db.Where("uid = ? AND cid = ?", uid, cid).First(&attendance).Error
	return &attendance, err
}

// GetAllAttendancesByCID CIDによって全ての出席情報を取得
func (repo *attendanceRepository) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := repo.db.Where("cid = ?", cid).Find(&attendances).Error
	return attendances, err
}

// GetAttendanceByID IDによって出席情報を取得
func (repo *attendanceRepository) GetAttendanceByID(id string) (*models.Attendance, error) {
	var attendance models.Attendance
	err := repo.db.Where("id = ?", id).First(&attendance).Error
	return &attendance, err
}

// UpdateAttendance 出席情報を更新
func (repo *attendanceRepository) UpdateAttendance(attendance *models.Attendance) error {
	return repo.db.Save(attendance).Error
}

// DeleteAttendance 出席情報を削除
func (repo *attendanceRepository) DeleteAttendance(id string) error {
	return repo.db.Delete(&models.Attendance{}, id).Error
}
