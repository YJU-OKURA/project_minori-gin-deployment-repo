package services

import (
	"github.com/gorilla/websocket"
)

type LiveClassService interface {
	CreateRoom() string
	InsertIntoRoom(roomID string, host bool, conn *websocket.Conn)
}

// liveClassService インタフェースを実装
type liveClassService struct {
	allRooms *RoomMap
}

// NewLiveClassService LiveClassServiceを生成
func NewLiveClassService(allRooms *RoomMap) LiveClassService {
	return &liveClassService{
		allRooms: allRooms,
	}
}

func (s *liveClassService) CreateRoom() string {
	return s.allRooms.CreateRoom()
}

func (s *liveClassService) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	s.allRooms.InsertIntoRoom(roomID, host, conn)
}
