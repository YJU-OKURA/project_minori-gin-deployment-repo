package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
)

// GroupBoardService インタフェース
type GroupBoardService interface {
	CreateGroupBoard(b dto.GroupBoardCreateDTO, imageUrl string) (*models.GroupBoard, error)
	GetAllGroupBoards() ([]models.GroupBoard, error)
	GetGroupBoardByID(id uint) (*models.GroupBoard, error)
	GetAnnouncedGroupBoards() ([]models.GroupBoard, error)
	UpdateGroupBoard(id uint, b dto.GroupBoardUpdateDTO, imageUrl string) (*models.GroupBoard, error) // Added imageUrl parameter
	DeleteGroupBoard(id uint) error
}

// groupBoardService インタフェースを実装
type groupBoardService struct {
	GroupBoardRepository repositories.GroupBoardRepository
	uploader             utils.Uploader
}

// NewGroupBoardService GroupBoardServiceを生成
func NewGroupBoardService(repo repositories.GroupBoardRepository) GroupBoardService {
	return &groupBoardService{GroupBoardRepository: repo}
}

// CreateGroupBoard 新しいグループ掲示板を作成
func (service *groupBoardService) CreateGroupBoard(b dto.GroupBoardCreateDTO, imageUrl string) (*models.GroupBoard, error) {
	groupBoard := models.GroupBoard{
		Title:       b.Title,
		Content:     b.Content,
		Image:       imageUrl,
		IsAnnounced: b.IsAnnounced,
		CID:         b.CID,
		UID:         b.UID,
	}
	return service.GroupBoardRepository.InsertGroupBoard(&groupBoard)
}

// GetGroupBoardByID IDでグループ掲示板を取得
func (service *groupBoardService) GetGroupBoardByID(id uint) (*models.GroupBoard, error) {
	return service.GroupBoardRepository.FindByID(id)
}

// GetAllGroupBoards 全てのグループ掲示板を取得
func (service *groupBoardService) GetAllGroupBoards() ([]models.GroupBoard, error) {
	return service.GroupBoardRepository.FindAll()
}

// GetAnnouncedGroupBoards 公開されたグループ掲示板を取得
func (service *groupBoardService) GetAnnouncedGroupBoards() ([]models.GroupBoard, error) {
	return service.GroupBoardRepository.FindAnnounced(true)
}

// UpdateGroupBoard 更新
func (service *groupBoardService) UpdateGroupBoard(id uint, b dto.GroupBoardUpdateDTO, imageUrl string) (*models.GroupBoard, error) {
	groupBoard, err := service.GetGroupBoardByID(id)
	if err != nil {
		return nil, err
	}

	if imageUrl != "" {
		groupBoard.Image = imageUrl
	}
	if b.Title != "" {
		groupBoard.Title = b.Title
	}
	if b.Content != "" {
		groupBoard.Content = b.Content
	}
	groupBoard.IsAnnounced = b.IsAnnounced

	err = service.GroupBoardRepository.UpdateGroupBoard(groupBoard)
	if err != nil {
		return nil, err
	}

	return groupBoard, nil
}

// DeleteGroupBoard 削除
func (service *groupBoardService) DeleteGroupBoard(id uint) error {
	return service.GroupBoardRepository.DeleteGroupBoard(id)
}
