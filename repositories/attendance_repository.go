package repositories

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/gorm"
)

type AttendanceRepository interface {
	CreateAttendance(attendance *models.Attendance) error
	GetAttendanceByUIDAndCSID(uid uint, cid uint) (*models.Attendance, error)
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAllAttendancesByCSID(csid uint) ([]models.Attendance, error)
	UpdateAttendance(attendance *models.Attendance) error
	DeleteAttendance(id uint) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (repo *attendanceRepository) CreateAttendance(attendance *models.Attendance) error {
	return repo.db.Create(attendance).Error
}

func (repo *attendanceRepository) GetAttendanceByUIDAndCSID(uid uint, cid uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := repo.db.
		Where("uid = ? AND cid = ?", uid, cid).
		First(&attendance).Error
	return &attendance, err
}

func (repo *attendanceRepository) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := repo.db.
		Preload("ClassUser.User").
		Preload("ClassUser.Class").
		Preload("ClassSchedule.Class").
		Where("cid = ?", cid).
		Find(&attendances).Error
	return attendances, err
}

func (repo *attendanceRepository) GetAllAttendancesByCSID(csid uint) ([]models.Attendance, error) {
	var attendances []models.Attendance
	err := repo.db.
		Preload("ClassUser.User").
		Preload("ClassUser.Class").
		Preload("ClassSchedule.Class").
		Where("csid = ?", csid).
		Find(&attendances).Error
	return attendances, err
}

func (repo *attendanceRepository) UpdateAttendance(attendance *models.Attendance) error {
	return repo.db.Save(attendance).Error
}

func (repo *attendanceRepository) DeleteAttendance(id uint) error {
	return repo.db.Delete(&models.Attendance{}, id).Error
}
