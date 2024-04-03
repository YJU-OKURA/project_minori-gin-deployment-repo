package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
)

// ClassBoardService インタフェース
type ClassBoardService interface {
	CreateClassBoard(b dto.ClassBoardCreateDTO) (*models.ClassBoard, error)
	GetAllClassBoards(cid uint, page int, pageSize int) ([]models.ClassBoard, error)
	GetClassBoardByID(id uint) (*models.ClassBoard, error)
	GetAnnouncedClassBoards(cid uint) ([]models.ClassBoard, error)
	UpdateClassBoard(id uint, b dto.ClassBoardUpdateDTO, imageUrl string) (*models.ClassBoard, error) // Added imageUrl parameter
	DeleteClassBoard(id uint) error
}

// classBoardService インタフェースを実装
type classBoardService struct {
	repo     repositories.ClassBoardRepository
	uploader utils.Uploader
}

// NewClassBoardService ClassClassServiceを生成
func NewClassBoardService(repo repositories.ClassBoardRepository) ClassBoardService {
	return &classBoardService{
		repo:     repo,
		uploader: utils.NewAwsUploader(),
	}
}

// CreateClassBoard 新しいグループ掲示板を作成
func (s *classBoardService) CreateClassBoard(b dto.ClassBoardCreateDTO) (*models.ClassBoard, error) {
	var imageUrl string
	var err error
	if b.Image != nil {
		imageUrl, err = s.uploader.UploadImage(b.Image, b.CID, false)
		if err != nil {
			return nil, err
		}
	}

	classBoard := models.ClassBoard{
		Title:       b.Title,
		Content:     b.Content,
		Image:       imageUrl,
		IsAnnounced: b.IsAnnounced,
		CID:         b.CID,
		UID:         b.UID,
	}
	return s.repo.InsertClassBoard(&classBoard)
}

// GetAllClassBoards 全てのグループ掲示板を取得
func (s *classBoardService) GetAllClassBoards(cid uint, page int, pageSize int) ([]models.ClassBoard, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindAllPaged(cid, pageSize, offset)
}

// GetClassBoardByID IDでグループ掲示板を取得
func (s *classBoardService) GetClassBoardByID(id uint) (*models.ClassBoard, error) {
	return s.repo.FindByID(id)
}

// GetAnnouncedClassBoards 公開されたグループ掲示板を取得
func (s *classBoardService) GetAnnouncedClassBoards(cid uint) ([]models.ClassBoard, error) {
	return s.repo.FindAnnounced(true, cid)
}

// UpdateClassBoard 更新
func (s *classBoardService) UpdateClassBoard(id uint, b dto.ClassBoardUpdateDTO, imageUrl string) (*models.ClassBoard, error) {
	classBoard, err := s.GetClassBoardByID(id)
	if err != nil {
		return nil, err
	}

	if imageUrl != "" {
		classBoard.Image = imageUrl
	}
	if b.Title != "" {
		classBoard.Title = b.Title
	}
	if b.Content != "" {
		classBoard.Content = b.Content
	}

	classBoard.IsAnnounced = b.IsAnnounced

	err = s.repo.UpdateClassBoard(classBoard)
	if err != nil {
		return nil, err
	}

	return classBoard, nil
}

// DeleteClassBoard 削除
func (s *classBoardService) DeleteClassBoard(id uint) error {
	return s.repo.DeleteClassBoard(id)
}
