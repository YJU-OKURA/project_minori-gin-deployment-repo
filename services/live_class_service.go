package services

import (
	"errors"
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var webrtcConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
}

type Room struct {
	ID             string
	Members        map[string]*websocket.Conn
	PeerConnection *webrtc.PeerConnection
	IsSharing      bool
	ClassID        uint
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
	CreateRoom(classID uint, adminID uint) (string, error)
	StartScreenShare(roomID string, userID string) error
	StopScreenShare(roomID string, adminID string) error
	GetScreenShareSDP(roomID string) (string, error)
	JoinScreenShare(roomID string, userID uint) (string, error)
	IsUserInClass(userID uint, classID uint) (bool, error)
}

type liveClassServiceImpl struct {
	roomMap       *RoomMap
	config        webrtc.Configuration
	classUserRepo repositories.ClassUserRepository
}

func NewLiveClassService(roomMap *RoomMap, classUserRepo repositories.ClassUserRepository) LiveClassService {
	return &liveClassServiceImpl{
		roomMap:       roomMap,
		config:        webrtcConfig,
		classUserRepo: classUserRepo,
	}
}

func (s *liveClassServiceImpl) setupPeerConnection() (*webrtc.PeerConnection, error) {
	m := &webrtc.MediaEngine{}
	// Register codecs
	if err := m.RegisterDefaultCodecs(); err != nil {
		return nil, err
	}

	// Create a new API with a MediaEngine containing the default codecs
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// Create a new PeerConnection with the configuration
	pc, err := api.NewPeerConnection(s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to create peer connection: %v", err)
	}

	return pc, nil
}

func (s *liveClassServiceImpl) CreateRoom(classID uint, userID uint) (string, error) {
	isAdmin, err := s.classUserRepo.IsAdmin(userID, classID)
	if err != nil {
		return "", err
	}
	if !isAdmin {
		return "", errors.New("unauthorized: only admins can create rooms")
	}

	s.roomMap.mu.Lock()
	defer s.roomMap.mu.Unlock()

	roomID := uuid.NewString()
	s.roomMap.rooms[roomID] = &Room{
		ID:        roomID,
		ClassID:   classID,
		Members:   make(map[string]*websocket.Conn),
		IsSharing: false,
	}
	return roomID, nil
}

func (s *liveClassServiceImpl) StartScreenShare(roomID string, userID string) error {
	s.roomMap.mu.Lock()
	room, exists := s.roomMap.rooms[roomID]
	s.roomMap.mu.Unlock()

	if !exists {
		log.Println("Attempt to access non-existent room:", roomID)
		return fmt.Errorf("room not found")
	}

	pc, err := s.setupPeerConnection()
	if err != nil {
		log.Println("Error setting up peer connection:", err)
		return err
	}

	//	Setup track
	track, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		pc.Close()
		log.Println("Failed to create track:", err)
		return err
	}

	if _, err = pc.AddTrack(track); err != nil {
		pc.Close()
		log.Println("Failed to add track:", err)
		return err
	}

	// Ensure peer connection is fully established before continuing
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		pc.Close()
		return fmt.Errorf("failed to create offer: %v", err)
	}

	if err = pc.SetLocalDescription(offer); err != nil {
		pc.Close()
		return fmt.Errorf("failed to set local description: %v", err)
	}

	// Handling ICE gathering
	done := make(chan bool)
	go awaitIceGatheringComplete(pc, done)

	if success := <-done; !success {
		pc.Close()
		return fmt.Errorf("failed to gather ICE candidates")
	}

	// Setup connection state monitoring to ensure the connection is ready
	connectionReady := make(chan bool)
	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		if state == webrtc.PeerConnectionStateConnected {
			close(connectionReady)
		}
	})

	select {
	case <-connectionReady:
		log.Println("Connection is fully established.")
	case <-time.After(30 * time.Second): // Wait for up to 30 seconds
		pc.Close()
		return fmt.Errorf("peer connection setup timed out")
	}

	room.PeerConnection = pc
	room.IsSharing = true
	log.Println("Started screen sharing in room:", roomID)
	return nil
}

func (s *liveClassServiceImpl) StopScreenShare(roomID string, adminID string) error {
	s.roomMap.mu.Lock()
	room, exists := s.roomMap.rooms[roomID]
	s.roomMap.mu.Unlock()

	if !exists {
		return fmt.Errorf("room not found")
	}

	if room.PeerConnection != nil {
		err := room.PeerConnection.Close()
		if err != nil {
			return err
		}
		room.PeerConnection = nil
		room.IsSharing = false
	}
	return nil
}

func (s *liveClassServiceImpl) GetScreenShareSDP(roomID string) (string, error) {
	s.roomMap.mu.RLock()
	room, exists := s.roomMap.rooms[roomID]
	s.roomMap.mu.RUnlock()

	if !exists || !room.IsSharing || room.PeerConnection == nil || room.PeerConnection.CurrentRemoteDescription() == nil {
		log.Println("Attempt to fetch SDP from inactive session:", roomID)
		return "", errors.New("screen share not active or connection closed")
	}
	return room.PeerConnection.LocalDescription().SDP, nil
}

func (s *liveClassServiceImpl) JoinScreenShare(roomID string, userID uint) (string, error) {
	s.roomMap.mu.RLock()
	room, exists := s.roomMap.rooms[roomID]
	s.roomMap.mu.RUnlock()

	if !exists || !room.IsSharing {
		return "", fmt.Errorf("room not found or no active sharing")
	}

	// Check if user is in the class associated with the room
	classID, err := s.getClassIDFromRoomID(roomID) // Implement this method based on your app's logic
	if err != nil {
		return "", err
	}

	isMember, err := s.classUserRepo.IsMember(userID, classID)
	if err != nil || !isMember {
		return "", fmt.Errorf("user is not a member of the class")
	}

	// Assume a simplified scenario where the viewer also sets up a peer connection
	pc, err := s.setupPeerConnection()
	if err != nil {
		return "", err
	}

	// Create an offer to send to the admin
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		pc.Close()
		return "", fmt.Errorf("failed to create offer: %v", err)
	}

	err = pc.SetLocalDescription(offer)
	if err != nil {
		pc.Close()
		return "", fmt.Errorf("failed to set local description: %v", err)
	}

	return offer.SDP, nil
}

func (s *liveClassServiceImpl) IsUserInClass(userID uint, classID uint) (bool, error) {
	// Use the IsAdmin method from the repository to check if the user is a member.
	// Assuming IsAdmin method checks for any user role in the class, not just admin.
	// You might need to adjust the logic based on your specific role requirements.
	return s.classUserRepo.IsAdmin(userID, classID)
}

func (s *liveClassServiceImpl) getClassIDFromRoomID(roomID string) (uint, error) {
	s.roomMap.mu.RLock()
	defer s.roomMap.mu.RUnlock()

	room, exists := s.roomMap.rooms[roomID]
	if !exists {
		return 0, fmt.Errorf("no room found for ID: %s", roomID)
	}
	return room.ClassID, nil
}

// awaitIceGatheringComplete waits for ICE candidates to be gathered before continuing.
func awaitIceGatheringComplete(pc *webrtc.PeerConnection, done chan bool) {
	gatherComplete := make(chan struct{})

	// Register the handler for ICE gathering state change
	pc.OnICEGatheringStateChange(func(state webrtc.ICEGatheringState) {
		log.Printf("ICE Gathering State has changed to %s\n", state.String())
		if state == webrtc.ICEGatheringStateComplete {
			close(gatherComplete)
		}
	})

	// Check if already complete before the handler
	if pc.ICEGatheringState() == webrtc.ICEGatheringStateComplete {
		close(gatherComplete)
	}

	select {
	case <-gatherComplete:
		log.Println("ICE gathering complete")
		done <- true
	case <-time.After(30 * time.Second): // Consider a longer timeout or make it configurable
		log.Println("ICE gathering timed out.")
		done <- false
	}
}
