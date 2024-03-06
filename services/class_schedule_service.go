package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
)

// ClassScheduleService インタフェース
type ClassScheduleService interface {
	CreateClassSchedule(classSchedule *models.ClassSchedule) (*models.ClassSchedule, error)
	GetClassScheduleByID(cid uint) (*models.ClassSchedule, error)
	GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error)
	UpdateClassSchedule(id uint, dto *dto.UpdateClassScheduleDTO) (*models.ClassSchedule, error)
	DeleteClassSchedule(id uint) error
	GetLiveClassSchedules(cid uint) ([]models.ClassSchedule, error)
	GetClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error)
}

// classScheduleService インタフェースを実装
type classScheduleService struct {
	repo repositories.ClassScheduleRepository
}

// NewClassScheduleService ClassScheduleServiceを生成
func NewClassScheduleService(repo repositories.ClassScheduleRepository) ClassScheduleService {
	return &classScheduleService{
		repo: repo,
	}
}

// GetClassScheduleByID クラススケジュールを取得
func (s *classScheduleService) GetClassScheduleByID(cid uint) (*models.ClassSchedule, error) {
	return s.repo.GetClassScheduleByID(cid)
}

// GetAllClassSchedules 全てのクラススケジュールを取得
func (s *classScheduleService) GetAllClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	return s.repo.GetAllClassSchedules(cid)
}

// CreateClassSchedule 新しいクラススケジュールを作成
func (s *classScheduleService) CreateClassSchedule(classSchedule *models.ClassSchedule) (*models.ClassSchedule, error) {
	err := s.repo.CreateClassSchedule(classSchedule)
	return classSchedule, err
}

// UpdateClassSchedule クラススケジュールを更新
func (s *classScheduleService) UpdateClassSchedule(id uint, dto *dto.UpdateClassScheduleDTO) (*models.ClassSchedule, error) {
	classSchedule, err := s.repo.GetClassScheduleByID(id)
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

	err = s.repo.UpdateClassSchedule(classSchedule)
	if err != nil {
		return nil, err
	}

	return classSchedule, nil
}

// DeleteClassSchedule クラススケジュールを削除
func (s *classScheduleService) DeleteClassSchedule(id uint) error {
	return s.repo.DeleteClassSchedule(id)
}

// GetLiveClassSchedules ライブ中のクラススケジュールを取得
func (s *classScheduleService) GetLiveClassSchedules(cid uint) ([]models.ClassSchedule, error) {
	return s.repo.FindLiveClassSchedules(cid)
}

// GetClassSchedulesByDate 日付でクラススケジュールを取得
func (s *classScheduleService) GetClassSchedulesByDate(cid uint, date string) ([]models.ClassSchedule, error) {
	return s.repo.FindClassSchedulesByDate(cid, date)
}
