package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"log"
	"sync"
)

var webrtcConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
	},
}

type Room struct {
	ID             string
	Members        map[*websocket.Conn]bool
	PeerConnection *webrtc.PeerConnection
}

type RoomMap struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewRoomMap() *RoomMap {
	return &RoomMap{
		rooms: make(map[string]*Room),
	}
}

type LiveClassService interface {
	CreateRoom() (string, error)
	InsertIntoRoom(roomID string, conn *websocket.Conn) error
	StartScreenShare(roomID string) (*webrtc.PeerConnection, error)
	StopScreenShare(roomID string) error
}

type liveClassService struct {
	roomMap *RoomMap
}

func NewLiveClassService(roomMap *RoomMap) LiveClassService {
	return &liveClassService{
		roomMap: roomMap,
	}
}

func (s *liveClassService) CreateRoom() (string, error) {
	return s.roomMap.CreateRoom()
}

func (s *liveClassService) InsertIntoRoom(roomID string, conn *websocket.Conn) error {
	return s.roomMap.InsertIntoRoom(roomID, conn)
}

func (m *RoomMap) CreateRoom() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	roomID := uuid.New().String()
	m.rooms[roomID] = &Room{
		ID:      roomID,
		Members: make(map[*websocket.Conn]bool),
	}
	log.Printf("Created new room with ID: %s", roomID)
	return roomID, nil
}
func (m *RoomMap) InsertIntoRoom(roomID string, conn *websocket.Conn) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, ok := m.rooms[roomID]
	if !ok {
		log.Printf("Attempted to access non-existing room: %s", roomID)
		return fmt.Errorf("no room found with ID %s", roomID)
	}

	room.Members[conn] = true
	log.Printf("Added connection to room: %s", roomID)
	return nil
}

func (s *liveClassService) StartScreenShare(roomID string) (*webrtc.PeerConnection, error) {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	room, exists := s.roomMap.rooms[roomID]
	if !exists {
		log.Printf("Room not found: %s", roomID)
		return nil, fmt.Errorf("room not found: %s", roomID)
	}

	pc, err := webrtc.NewPeerConnection(webrtcConfig)
	if err != nil {
		log.Printf("Failed to create peer connection: %v", err)
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		pc.Close() // Ensure resources are freed on error
		log.Printf("Failed to create video track: %v", err)
		return nil, fmt.Errorf("failed to create video track: %v", err)
	}

	if _, err = pc.AddTrack(videoTrack); err != nil {
		pc.Close() // Ensure resources are freed on error
		log.Printf("Failed to add track: %v", err)
		return nil, fmt.Errorf("failed to add track: %v", err)
	}

	room.PeerConnection = pc
	log.Printf("Started screen sharing in room ID: %s", roomID)
	return pc, nil
}

func (s *liveClassService) StopScreenShare(roomID string) error {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	room, exists := s.roomMap.rooms[roomID]
	if !exists {
		log.Printf("Attempt to stop screen sharing in non-existing room: %s", roomID)
		return fmt.Errorf("no room found with ID %s", roomID)
	}

	if room.PeerConnection != nil {
		if err := room.PeerConnection.Close(); err != nil {
			log.Printf("Failed to close peer connection for room %s: %v", roomID, err)
			return err
		}
		room.PeerConnection = nil
		log.Printf("Stopped screen sharing in room ID: %s", roomID)
	}

	return nil
}
