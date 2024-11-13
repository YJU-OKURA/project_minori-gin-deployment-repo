package services

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"gorm.io/gorm"
)

type AttendanceService interface {
	CreateOrUpdateAttendance(cid uint, uid uint, csid uint, status string) error
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAllAttendancesByCSID(csid uint) ([]models.Attendance, error)
	DeleteAttendance(id uint) error
}

type attendanceService struct {
	repo repositories.AttendanceRepository
}

func NewAttendanceService(repo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{repo: repo}
}

func (s *attendanceService) CreateOrUpdateAttendance(cid uint, uid uint, csid uint, status string) error {
	attendance, err := s.repo.GetAttendanceByUIDAndCSID(uid, csid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAttendance := models.Attendance{
				CID:          cid,
				UID:          uid,
				CSID:         csid,
				IsAttendance: models.AttendanceType(status),
			}
			return s.repo.CreateAttendance(&newAttendance)
		}
		return err
	}

	attendance.IsAttendance = models.AttendanceType(status)
	return s.repo.UpdateAttendance(attendance)
}

func (s *attendanceService) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	return s.repo.GetAllAttendancesByCID(cid)
}

func (s *attendanceService) GetAllAttendancesByCSID(csid uint) ([]models.Attendance, error) {
	return s.repo.GetAllAttendancesByCSID(csid)
}

func (s *attendanceService) DeleteAttendance(id uint) error {
	return s.repo.DeleteAttendance(id)
}
