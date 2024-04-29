package services

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"net/http"
	"sync"
)

// ClassBoardService インタフェース
type ClassBoardService interface {
	CreateClassBoard(b dto.ClassBoardCreateDTO) (*models.ClassBoard, error)
	GetAllClassBoards(cid uint, page int, pageSize int) ([]models.ClassBoard, error)
	GetClassBoardByID(id uint) (*models.ClassBoard, error)
	GetAnnouncedClassBoards(cid uint) ([]models.ClassBoard, error)
	UpdateClassBoard(id uint, b dto.ClassBoardUpdateDTO, imageUrl string) (*models.ClassBoard, error) // Added imageUrl parameter
	DeleteClassBoard(id uint) error
	GetUpdateNotifier() *UpdateNotifier
}

// classBoardService インタフェースを実装
type classBoardService struct {
	repo     repositories.ClassBoardRepository
	uploader utils.Uploader
	notifier *UpdateNotifier
}

// NewClassBoardService ClassClassServiceを生成
func NewClassBoardService(repo repositories.ClassBoardRepository) ClassBoardService {
	notifier := NewUpdateNotifier()
	return &classBoardService{
		repo:     repo,
		uploader: utils.NewAwsUploader(),
		notifier: notifier,
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

type UpdateNotifier struct {
	Register   chan http.ResponseWriter
	Unregister chan http.ResponseWriter
	Broadcast  chan []byte
	clients    map[http.ResponseWriter]struct{}
	mu         sync.Mutex
}

func NewUpdateNotifier() *UpdateNotifier {
	notifier := &UpdateNotifier{
		Register:   make(chan http.ResponseWriter),
		Unregister: make(chan http.ResponseWriter),
		Broadcast:  make(chan []byte),
		clients:    make(map[http.ResponseWriter]struct{}),
	}
	go notifier.run()
	return notifier
}

func (u *UpdateNotifier) run() {
	for {
		select {
		case s := <-u.Register:
			u.mu.Lock()
			u.clients[s] = struct{}{}
			u.mu.Unlock()
		case s := <-u.Unregister:
			u.mu.Lock()
			delete(u.clients, s)
			u.mu.Unlock()
		case msg := <-u.Broadcast:
			u.mu.Lock()
			for s := range u.clients {
				_, _ = s.Write(msg)
				if f, ok := s.(http.Flusher); ok {
					f.Flush()
				}
			}
			u.mu.Unlock()
		}
	}
}

func (s *classBoardService) GetUpdateNotifier() *UpdateNotifier {
	return s.notifier
}
