package services

//import (
//	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
//	"github.com/gorilla/websocket"
//	"github.com/nareix/joy4/format"
//	"github.com/nareix/joy4/format/rtmp"
//	"log"
//)
//
//func init() {
//	format.RegisterAll()
//}
//
//type RTMPServer struct {
//	server *rtmp.Server
//}
//
//func NewRTMPServer() *RTMPServer {
//	server := &rtmp.Server{}
//
//	// Handle RTMP publish
//	server.HandlePublish = func(conn *rtmp.Conn) {
//		//	Handle incoming stream, for example, save it, broadcast it, etc.
//	}
//
//	// Handle RTMP play
//	server.HandlePlay = func(conn *rtmp.Conn) {
//		//	Handle playback of stream
//	}
//
//	return &RTMPServer{server: server}
//}
//
//func (s *RTMPServer) Start(address string) error {
//	log.Printf("Starting RTMP server on %s", address)
//	return s.server.ListenAndServe(address)
//}
//
//type LiveClassService interface {
//	CreateRoom() (string, error)
//	InsertIntoRoom(roomID string, host bool, conn *websocket.Conn)
//}
//
//// liveClassService インタフェースを実装
//type liveClassService struct {
//	allRooms   *controllers.RoomMap
//	rtmpServer *RTMPServer
//}
//
//// NewLiveClassService LiveClassServiceを生成
//func NewLiveClassService(allRooms *controllers.RoomMap) LiveClassService {
//	return &liveClassService{
//		allRooms:   allRooms,
//		rtmpServer: NewRTMPServer(),
//	}
//}
//
//func (s *liveClassService) CreateRoom() (string, error) {
//	roomID, err := s.allRooms.CreateRoom()
//	if err != nil {
//		// Handle error, for example, log it or wrap it with more context
//		return "", err
//	}
//	return roomID, nil
//}
//
//func (s *liveClassService) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
//	s.allRooms.InsertIntoRoom(roomID, host, conn)
//}
