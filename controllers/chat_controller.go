package controllers

import (
	"context"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
)

// ChatController インタフェースを実装
type ChatController struct {
	chatManager *services.Manager
	redisClient *redis.Client
}

type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// NewChatController ChatControllerを生成
func NewChatController(chatMgr *services.Manager, redisClient *redis.Client) *ChatController {
	return &ChatController{
		chatManager: chatMgr,
		redisClient: redisClient,
	}
}

// HandleChatRoom godoc
// @Summary チャットルームをハンドル
// @Description チャットルームをハンドルする。
// @Tags Chat Room
// @Accept html
// @Produce html
// @Param scheduleId path string true "Schedule ID"
// @Param userId path string true "User ID"
// @Success 200 {string} string "チャットルームが正常にハンドルされました"
// @Router /chat/room/{scheduleId}/{userId} [get]
func (controller *ChatController) HandleChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	userId := ctx.Param("userId")

	ctx.HTML(http.StatusOK, "chat_room", gin.H{
		"scheduleId": scheduleId,
		"userid":     userId,
	})
}

// CreateChatRoom godoc
// @Summary チャットルームを作成
// @Description チャットルームを作成する。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {object} string "Chat room created successfully"
// @Failure 400 {object} string "Failed to create chat room"
// @Router /chat/create-room/{scheduleId} [post]
func (controller *ChatController) CreateChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")

	controller.chatManager.CreateRoom(scheduleId)

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"status": "Chat room created successfully",
	})
}

// PostToChatRoom godoc
// @Summary チャットルームに投稿
// @Description チャットルームにメッセージを投稿する。
// @Tags Chat Room
// @Accept multipart/form-data
// @Produce json
// @Param scheduleId path int true "Schedule ID"
// @Param user formData string true "User ID"
// @Param message formData string true "Message"
// @Success 200 {object} string "success"
// @Router /chat/room/{scheduleId} [post]
func (controller *ChatController) PostToChatRoom(ctx *gin.Context) {
	user := ctx.PostForm("user")
	message := ctx.PostForm("message")
	scheduleId := ctx.Param("scheduleId")

	if user == "" || message == "" {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid request: user and message must be provided")
		return
	}

	controller.chatManager.Submit(user, scheduleId, message)

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"status":  "Message sent",
		"message": message,
	})
}

// DeleteChatRoom godoc
// @Summary チャットルームを削除
// @Description チャットルームを削除する。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param scheduleId path string true "Schedule ID"
// @Success 200 {object} string "success"
// @Router /chat/room/{scheduleId} [delete]
func (controller *ChatController) DeleteChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	controller.chatManager.DeleteBroadcast(scheduleId)

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"status": constants.DeleteSuccess,
	})
}

// StreamChat godoc
// @Summary チャットをストリーム
// @Description チャットをストリームする。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param scheduleId path int true "Schedule ID"
// @Router /chat/stream/{scheduleId} [get]
func (controller *ChatController) StreamChat(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	listener := controller.chatManager.OpenListener(scheduleId)
	defer controller.chatManager.CloseListener(scheduleId, listener)

	clientGone := ctx.Request.Context().Done()
	ctx.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			ctx.SSEvent("message", message)
			return true
		}
	})
}

// GetChatMessages godoc
// @Summary チャットメッセージを取得
// @Description チャットメッセージを取得する。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param roomid path string true "Room ID"
// @Success 200 {object} string "success"
// @Router /chat/messages/{roomid} [get]
func (controller *ChatController) GetChatMessages(ctx *gin.Context) {
	roomid := ctx.Param("roomid")
	messages, err := controller.redisClient.LRange(context.Background(), "chat:"+roomid, 0, -1).Result()
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.ErrLoadMessage)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"messages": messages,
	})
}

// SendDirectMessage godoc
// @Summary Send a direct message
// @Description Send a direct message to a specific user
// @Tags Direct Message
// @Accept json
// @Produce json
// @Param senderId path string true "Sender ID"
// @Param receiverId path string true "Receiver ID"
// @Param message formData string true "Message"
// @Success 200 {object} string "Message sent successfully"
// @Router /chat/dm/{senderId}/{receiverId} [post]
func (controller *ChatController) SendDirectMessage(ctx *gin.Context) {
	senderId := ctx.Param("senderId")
	receiverId := ctx.Param("receiverId")
	message := ctx.PostForm("message")

	if err := controller.chatManager.SubmitDirectMessage(senderId, receiverId, message); err != nil {
		respondWithError(ctx, http.StatusInternalServerError, constants.ErrSendMessage)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"status":  constants.MessageSent,
		"message": message,
	})
}

// GetDirectMessages godoc
// @Summary Get direct messages
// @Description Get direct messages history between two users
// @Tags Direct Message
// @Accept json
// @Produce json
// @Param userId1 path string true "User ID 1"
// @Param userId2 path string true "User ID 2"
// @Success 200 {object} string "Messages fetched successfully"
// @Router /chat/dm/{userId1}/{userId2} [get]
func (controller *ChatController) GetDirectMessages(ctx *gin.Context) {
	userId1 := ctx.Param("userId1")
	userId2 := ctx.Param("userId2")

	messages, err := controller.chatManager.GetDirectMessages(userId1, userId2)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.ErrLoadMessage)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"messages": messages,
	})
}
