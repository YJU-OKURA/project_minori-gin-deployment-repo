package services

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
	"log"
	"strconv"
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
	IsSharing      bool
	AdminID        string
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
	StartScreenShare(roomID string, userID string) (*webrtc.PeerConnection, error)
	StopScreenShare(roomID string, userID string) error
}

// liveClassService implements the LiveClassService interface to manage live classroom sessions.

type liveClassService struct {
	roomMap       *RoomMap
	classUserRepo repositories.ClassUserRepository
}

func NewLiveClassService(roomMap *RoomMap, repo repositories.ClassUserRepository) LiveClassService {
	return &liveClassService{
		roomMap:       roomMap,
		classUserRepo: repo,
	}
}

func (s *liveClassService) CreateRoom() (string, error) {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	roomID := uuid.New().String()
	s.roomMap.rooms[roomID] = &Room{
		ID:      roomID,
		Members: make(map[*websocket.Conn]bool),
	}
	log.Printf("Created new room with ID: %s", roomID)
	return roomID, nil
}

func (s *liveClassService) InsertIntoRoom(roomID string, conn *websocket.Conn) error {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	room, ok := s.roomMap.rooms[roomID]
	if !ok {
		return fmt.Errorf("no room found with ID %s", roomID)
	}

	room.Members[conn] = true
	return nil
}

func (s *liveClassService) StartScreenShare(roomID string, userID string) (*webrtc.PeerConnection, error) {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	room, exists := s.roomMap.rooms[roomID]
	if !exists || !s.isRoomAdmin(userID, roomID) {
		return nil, fmt.Errorf("room %s does not exist", roomID)
	}
	if !s.isRoomAdmin(userID, roomID) {
		return nil, fmt.Errorf("user %s is not an admin of room %s", userID, roomID)
	}
	if room.IsSharing {
		return nil, fmt.Errorf("screen sharing is already in progress in room %s", roomID)
	}

	pc, err := webrtc.NewPeerConnection(webrtcConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		pc.Close()
		return nil, fmt.Errorf("failed to create video track: %v", err)
	}

	if _, err = pc.AddTrack(videoTrack); err != nil {
		pc.Close()
		return nil, fmt.Errorf("failed to add track: %v", err)
	}

	room.PeerConnection = pc
	room.IsSharing = true
	log.Printf("Screen sharing started in room %s by user %s", roomID, userID)
	return pc, nil
}

func (s *liveClassService) StopScreenShare(roomID string, userID string) error {
	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	room, exists := s.roomMap.rooms[roomID]
	if !exists || !s.isRoomAdmin(userID, roomID) {
		return fmt.Errorf("room %s does not exist", roomID)
	}

	if room.PeerConnection != nil {
		if err := room.PeerConnection.Close(); err != nil {
			log.Printf("Failed to close peer connection for room %s: %v", roomID, err)
			return err
		}
		room.PeerConnection = nil
		log.Printf("Stopped screen sharing in room ID: %s", roomID)
	}

	room.IsSharing = false
	return nil
}

func (s *liveClassService) isRoomAdmin(userID string, roomID string) bool {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		log.Printf("Error converting userID to uint: %v", err)
		return false
	}

	cid, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		log.Printf("Error converting roomID to uint: %v", err)
		return false
	}

	roleID, err := s.classUserRepo.GetRole(uint(uid), uint(cid))
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		return false
	}
	return roleID == 2
}
