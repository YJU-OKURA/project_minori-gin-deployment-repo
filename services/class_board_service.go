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
	GetAllClassBoards(cid uint) ([]models.ClassBoard, error)
	GetClassBoardByID(id uint) (*models.ClassBoard, error)
	GetAnnouncedClassBoards(cid uint) ([]models.ClassBoard, error)
	UpdateClassBoard(id uint, b dto.ClassBoardUpdateDTO, imageUrl string) (*models.ClassBoard, error) // Added imageUrl parameter
	DeleteClassBoard(id uint) error
}

// classBoardService インタフェースを実装
type classBoardService struct {
	ClassBoardRepository repositories.ClassBoardRepository
	uploader             utils.Uploader
}

// classBoardServiceImpl はClassBoardServiceの実装です。
type classBoardServiceImpl struct {
	Repo *repositories.ClassBoardRepository
}

// NewClassBoardService ClassClassServiceを生成
func NewClassBoardService(repo repositories.ClassBoardRepository) ClassBoardService {
	return &classBoardService{
		ClassBoardRepository: repo,
		uploader:             utils.NewAwsUploader(),
	}
}

// CreateClassBoard 新しいグループ掲示板を作成
// func (service *classBoardService) CreateClassBoard(b dto.ClassBoardCreateDTO, imageUrl string) (*models.ClassBoard, error) {
func (service *classBoardService) CreateClassBoard(b dto.ClassBoardCreateDTO) (*models.ClassBoard, error) {
	var imageUrl string
	var err error
	if b.Image != nil {
		imageUrl, err = service.uploader.UploadImage(b.Image)
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
	return service.ClassBoardRepository.InsertClassBoard(&classBoard)
}

// GetClassBoardByID IDでグループ掲示板を取得
func (service *classBoardService) GetClassBoardByID(id uint) (*models.ClassBoard, error) {
	return service.ClassBoardRepository.FindByID(id)
}

// GetAllClassBoards 全てのグループ掲示板を取得
func (service *classBoardService) GetAllClassBoards(cid uint) ([]models.ClassBoard, error) {
	return service.ClassBoardRepository.FindAll(cid)
}

// GetAnnouncedClassBoards 公開されたグループ掲示板を取得
func (service *classBoardService) GetAnnouncedClassBoards(cid uint) ([]models.ClassBoard, error) {
	return service.ClassBoardRepository.FindAnnounced(true, cid)
}

// UpdateClassBoard 更新
func (service *classBoardService) UpdateClassBoard(id uint, b dto.ClassBoardUpdateDTO, imageUrl string) (*models.ClassBoard, error) {
	classBoard, err := service.GetClassBoardByID(id)
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

	err = service.ClassBoardRepository.UpdateClassBoard(classBoard)
	if err != nil {
		return nil, err
	}

	return classBoard, nil
}

// DeleteClassBoard 削除
func (service *classBoardService) DeleteClassBoard(id uint) error {
	return service.ClassBoardRepository.DeleteClassBoard(id)
}
