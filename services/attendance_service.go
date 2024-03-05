package services

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"gorm.io/gorm"
)

// AttendanceService インタフェース
type AttendanceService interface {
	CreateOrUpdateAttendance(cid, uid uint, status string) error
	GetAllAttendancesByCID(cid uint) ([]models.Attendance, error)
	GetAttendanceByID(id string) (*models.Attendance, error)
	DeleteAttendance(id string) error
}

// attendanceService インタフェースを実装
type attendanceService struct {
	AttendanceRepository repositories.AttendanceRepository
}

// attendanceServiceImpl はAttendanceの実装です。
type attendanceServiceImpl struct {
	Repo *repositories.AttendanceRepository
}

// NewAttendanceService AttendanceServiceを生成
func NewAttendanceService(repo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{
		AttendanceRepository: repo,
	}
}

// GetAllAttendancesByCID CIDによって全ての出席情報を取得
func (service *attendanceService) GetAllAttendancesByCID(cid uint) ([]models.Attendance, error) {
	return service.AttendanceRepository.GetAllAttendancesByCID(cid)
}

// GetAttendanceByID IDによって出席情報を取得
func (service *attendanceService) GetAttendanceByID(id string) (*models.Attendance, error) {
	return service.AttendanceRepository.GetAttendanceByID(id)
}

// CreateOrUpdateAttendance 出席情報を作成または更新
func (service *attendanceService) CreateOrUpdateAttendance(cid, uid uint, status string) error {
	attendance, err := service.AttendanceRepository.GetAttendanceByUIDAndCID(uid, cid)
	if err != nil {
		// レコードが見つからない場合は新規作成
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newAttendance := models.Attendance{
				CID:          cid,
				UID:          uid,
				IsAttendance: status,
			}
			return service.AttendanceRepository.CreateAttendance(&newAttendance)
		}
		return err
	}

	// レコードが見つかった場合は更新
	attendance.IsAttendance = status
	return service.AttendanceRepository.UpdateAttendance(attendance)
}

// DeleteAttendance 出席情報を削除
func (service *attendanceService) DeleteAttendance(id string) error {
	return service.AttendanceRepository.DeleteAttendance(id)
}
