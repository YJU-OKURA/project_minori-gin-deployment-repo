package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

// AttendanceRepository インタフェース
type AttendanceRepository interface {
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceByUIDAndCID(uid uint, cid uint) (*models.Attendance, error)
	UpdateAttendance(attendance *models.Attendance) error
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAttendanceByID(id string) (*models.Attendance, error)
	DeleteAttendance(id string) error
}

// attendanceConnection グループ掲示板リポジトリ
type attendanceConnection struct {
	DB *gorm.DB
}

// NewAttendanceRepository グループ掲示板リポジトリを生成
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceConnection{DB: db}
}

// CreateAttendance 出席情報を作成
func (db *attendanceConnection) CreateAttendance(attendance *models.Attendance) error {
	return db.DB.Create(attendance).Error
}

// GetAttendanceByUIDAndCID UIDとCIDによって出席情報を取得
func (db *attendanceConnection) GetAttendanceByUIDAndCID(uid uint, cid uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := db.DB.Where("uid = ? AND cid = ?", uid, cid).First(&attendance).Error
	return &attendance, err
}

// UpdateAttendance 出席情報を更新
func (db *attendanceConnection) UpdateAttendance(attendance *models.Attendance) error {
	return db.DB.Save(attendance).Error
}

// GetAllAttendancesByCID CIDによって全ての出席情報を取得
func (db *attendanceConnection) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := db.DB.Where("cid = ?", cid).Find(&attendances).Error
	return attendances, err
}

// GetAttendanceByID IDによって出席情報を取得
func (db *attendanceConnection) GetAttendanceByID(id string) (*models.Attendance, error) {
	var attendance models.Attendance
	err := db.DB.Where("id = ?", id).First(&attendance).Error
	return &attendance, err
}

// DeleteAttendance 出席情報を削除
func (db *attendanceConnection) DeleteAttendance(id string) error {
	return db.DB.Delete(&models.Attendance{}, id).Error
}
