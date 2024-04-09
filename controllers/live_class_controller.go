package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Participant struct {
	Host bool
	Conn *websocket.Conn
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[string][]Participant
}

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomID]
}

func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomID := uuid.New().String()
	r.Map[roomID] = []Participant{}

	r.Map[roomID] = []Participant{}

	return roomID
}

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{host, conn}

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}

// LiveClassController インタフェースを実装
type LiveClassController struct {
	liveClassService services.LiveClassService
}

// NewLiveClassController LiveClassControllerを生成
func NewLiveClassController(service services.LiveClassService) *LiveClassController {
	return &LiveClassController{
		liveClassService: service,
	}
}

func (c *LiveClassController) CreateRoomRequestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := c.liveClassService.CreateRoom()
		ctx.JSON(http.StatusOK, gin.H{"room_id": roomID})
	}
}

func (c *LiveClassController) JoinRoomRequestHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
