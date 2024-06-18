package services

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"gorm.io/gorm"
)

// AttendanceService インタフェース
type AttendanceService interface {
	CreateOrUpdateAttendance(cid uint, uid uint, csid uint, status string) error
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAttendanceByID(id string) ([]models.Attendance, error)
	DeleteAttendance(id string) error
}

// attendanceService インタフェースを実装
type attendanceService struct {
	repo repositories.AttendanceRepository
}

// NewAttendanceService AttendanceServiceを生成
func NewAttendanceService(repo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{
		repo: repo,
	}
}

// CreateOrUpdateAttendance 出席情報を作成または更新
func (s *attendanceService) CreateOrUpdateAttendance(cid uint, uid uint, csid uint, status string) error {
	attendance, err := s.repo.GetAttendanceByUIDAndCID(uid, cid)
	if err != nil {
		// レコードが見つからない場合は新規作成
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

	// レコードが見つかった場合は更新
	attendance.IsAttendance = models.AttendanceType(status)
	return s.repo.UpdateAttendance(attendance)
}

// GetAllAttendancesByCID CIDによって全ての出席情報を取得
func (s *attendanceService) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	return s.repo.GetAllAttendancesByCID(cid)
}

// GetAttendanceByID IDによって出席情報を取得
func (s *attendanceService) GetAttendanceByID(id string) ([]models.Attendance, error) {
	attendances, err := s.repo.GetAttendanceByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return attendances, nil
}

// DeleteAttendance 出席情報を削除
func (s *attendanceService) DeleteAttendance(id string) error {
	return s.repo.DeleteAttendance(id)
}
