package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-broadcast"
	"github.com/go-redis/redis/v8"
)

// Message ユーザーとルームの識別子を持つチャットメッセージを表す
type Message struct {
	UserId     string
	RoomId     string // もしIsDMがtrueならば、RoomIdはnull
	ReceiverId string // もしIsDMがtrueならば、ReceiverIdはnullになれない
	Text       string
	IsDM       bool
}

// Listener 特定のルームの着信チャットメッセージを処理
type Listener struct {
	RoomId string
	Chan   chan interface{}
}

// Manager チャットルームの管理を行う
type Manager struct {
	roomChannels map[string]broadcast.Broadcaster
	open         chan *Listener
	close        chan *Listener
	delete       chan string
	messages     chan *Message
	redisClient  *redis.Client
}

// NewRoomManager function マネージャーを作成
func NewRoomManager(redisClient *redis.Client) *Manager {
	manager := &Manager{
		roomChannels: make(map[string]broadcast.Broadcaster),
		open:         make(chan *Listener, 100),
		close:        make(chan *Listener, 100),
		delete:       make(chan string, 100),
		messages:     make(chan *Message, 100),
		redisClient:  redisClient,
	}

	go manager.run()
	return manager
}

// run マネージャーを実行
func (m *Manager) run() {
	for {
		select {
		case listener := <-m.open:
			m.register(listener)
		case listener := <-m.close:
			m.deregister(listener)
		case roomid := <-m.delete:
			m.deleteBroadcast(roomid)
		case message := <-m.messages:
			m.room(message.RoomId).Submit(message.UserId + ": " + message.Text)
		}
	}
}

// register リスナーを登録
func (m *Manager) register(listener *Listener) {
	m.room(listener.RoomId).Register(listener.Chan)
}

// deregister リスナーを登録解除
func (m *Manager) deregister(listener *Listener) {
	m.room(listener.RoomId).Unregister(listener.Chan)
	close(listener.Chan)
}

// deleteBroadcast ブロードキャストを削除
func (m *Manager) deleteBroadcast(roomid string) {
	b, ok := m.roomChannels[roomid]
	if ok {
		err := b.Close()
		if err != nil {
			return
		}
		delete(m.roomChannels, roomid)
	}
}

// room ルームを取得
func (m *Manager) room(roomid string) broadcast.Broadcaster {
	b, ok := m.roomChannels[roomid]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		m.roomChannels[roomid] = b
	}
	return b
}

// OpenListener リスナーを開く
func (m *Manager) OpenListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	m.open <- &Listener{
		RoomId: roomid,
		Chan:   listener,
	}
	return listener
}

// CloseListener リスナーを閉じる
func (m *Manager) CloseListener(roomid string, channel chan interface{}) {
	m.close <- &Listener{
		RoomId: roomid,
		Chan:   channel,
	}
}

// Submit メッセージを送信
func (m *Manager) Submit(userid, roomid, text string) {
	msg := &Message{
		UserId: userid,
		RoomId: roomid,
		Text:   text,
	}
	m.messages <- msg

	// Redisにメッセージを保存
	key := "chat:" + roomid
	err := m.redisClient.RPush(context.Background(), "chat:"+roomid, fmt.Sprintf("%s: %s", userid, text)).Err()
	if err != nil {
		log.Printf("Redis error: %v", err)
	}

	// メッセージの有効期限を設定(e.g. , 1時間)
	msgErr := m.redisClient.Expire(context.Background(), key, time.Hour).Err()
	if msgErr != nil {
		return
	}
}

// SubmitDirectMessage ダイレクトメッセージを送信
func (m *Manager) SubmitDirectMessage(senderId, receiverId, text string) error {
	msg := &Message{
		UserId:     senderId,
		ReceiverId: receiverId,
		Text:       text,
		IsDM:       true,
	}

	messageJSON, _ := json.Marshal(msg)
	key := "dm:" + senderId + ":" + receiverId
	if err := m.pushToRedis(key, messageJSON); err != nil {
		log.Printf("Redis error: %v", err)
		return err
	}

	// メッセージの有効期限を設定(e.g. , 1時間)
	err := m.redisClient.Expire(context.Background(), key, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) pushToRedis(key string, data []byte) error {
	if err := m.redisClient.RPush(context.Background(), key, data).Err(); err != nil {
		return err
	}

	// Set the expiration of the message to 1 hour
	if err := m.redisClient.Expire(context.Background(), key, time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

// GetDirectMessages ダイレクトメッセージを取得
func (m *Manager) GetDirectMessages(senderId, receiverId string) ([]Message, error) {
	key := "dm:" + senderId + ":" + receiverId // e.g. dm:1:2
	messagesJSON, err := m.redisClient.LRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var messages []Message
	for _, mJSON := range messagesJSON {
		var msg Message
		if err := json.Unmarshal([]byte(mJSON), &msg); err != nil {
			continue // TODO: メッセージのデコードに失敗した場合は無視中なんですが修正が必要かもしれません
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (m *Manager) CreateRoom(roomID string) {
	if _, exists := m.roomChannels[roomID]; !exists {
		m.roomChannels[roomID] = broadcast.NewBroadcaster(10)
		fmt.Println("Chat room created: ", roomID)
	} else {
		log.Printf("Attempted to create an already existing room: %s", roomID)
	}
}

func (m *Manager) DeleteBroadcast(roomID string) {
	b, ok := m.roomChannels[roomID]
	if ok {
		err := b.Close()
		if err != nil {
			log.Printf("Error closing broadcaster for room %s: %v", roomID, err)
			return
		}
		delete(m.roomChannels, roomID)
		delErr := m.redisClient.Del(context.Background(), "chat:"+roomID).Err()
		if delErr != nil {
			log.Printf("Error deleting Redis key for room %s: %v", roomID, delErr)
		}
		log.Printf("Chat room deleted: %s", roomID)
	} else {
		log.Printf("Attempted to delete a non-existing room: %s", roomID)
	}
}

func (m *Manager) DeleteDirectMessages(senderId, receiverId string) error {
	key := "dm:" + senderId + ":" + receiverId
	if err := m.redisClient.Del(context.Background(), key).Err(); err != nil {
		log.Printf("Error deleting DMs from Redis: %v", err)
		return err
	}
	return nil
}
