package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// ClassScheduleService インタフェース
type ClassScheduleService interface {
	GetClassScheduleByID(cid uint) (*models.ClassSchedule, error)
	GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error)
	CreateClassSchedule(classSchedule *models.ClassSchedule) (*models.ClassSchedule, error)
	UpdateClassSchedule(id uint, dto *dto.UpdateClassScheduleDTO) (*models.ClassSchedule, error)
	DeleteClassSchedule(id uint) error
	GetLiveClassSchedules(cid uint) ([]models.ClassSchedule, error)
	GetClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error)
}

// classScheduleService インタフェースを実装
type classScheduleService struct {
	ClassScheduleRepository repositories.ClassScheduleRepository
}

// classScheduleServiceImpl はClassScheduleServiceの実装です。
type classScheduleServiceImpl struct {
	Repo *repositories.ClassScheduleRepository
}

// NewClassScheduleService ClassScheduleServiceを生成
func NewClassScheduleService(repo repositories.ClassScheduleRepository) ClassScheduleService {
	return &classScheduleService{
		ClassScheduleRepository: repo,
	}
}

// GetClassScheduleByID クラススケジュールを取得
func (service *classScheduleService) GetClassScheduleByID(cid uint) (*models.ClassSchedule, error) {
	return service.ClassScheduleRepository.GetClassScheduleByID(cid)
}

// GetAllClassSchedules 全てのクラススケジュールを取得
func (service *classScheduleService) GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	return service.ClassScheduleRepository.GetAllClassSchedules(cid)
}

// CreateClassSchedule 新しいクラススケジュールを作成
func (service *classScheduleService) CreateClassSchedule(classSchedule *models.ClassSchedule) (*models.ClassSchedule, error) {
	err := service.ClassScheduleRepository.CreateClassSchedule(classSchedule)
	return classSchedule, err
}

// UpdateClassSchedule クラススケジュールを更新
func (service *classScheduleService) UpdateClassSchedule(id uint, dto *dto.UpdateClassScheduleDTO) (*models.ClassSchedule, error) {
	classSchedule, err := service.ClassScheduleRepository.GetClassScheduleByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Title != nil {
		classSchedule.Title = *dto.Title
	}
	if dto.StartedAt != nil {
		classSchedule.StartedAt = *dto.StartedAt
	}
	if dto.EndedAt != nil {
		classSchedule.EndedAt = *dto.EndedAt
	}
	if dto.IsLive != nil {
		classSchedule.IsLive = *dto.IsLive
	}

	err = service.ClassScheduleRepository.UpdateClassSchedule(classSchedule)
	if err != nil {
		return nil, err
	}

	return classSchedule, nil
}

// DeleteClassSchedule クラススケジュールを削除
func (service *classScheduleService) DeleteClassSchedule(id uint) error {
	return service.ClassScheduleRepository.DeleteClassSchedule(id)
}

// GetLiveClassSchedules ライブ中のクラススケジュールを取得
func (service *classScheduleService) GetLiveClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	return service.ClassScheduleRepository.FindLiveClassSchedules(cid)
}

// GetClassSchedulesByDate 日付でクラススケジュールを取得
func (service *classScheduleService) GetClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error) {
	return service.ClassScheduleRepository.FindClassSchedulesByDate(cid, date)
}
