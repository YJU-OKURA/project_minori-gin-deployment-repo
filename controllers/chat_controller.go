package controllers

import (
	"context"
	"io"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// ChatController チャットコントローラ
type ChatController struct {
	chatManager *services.Manager
	redisClient *redis.Client
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
// @Param scheduleId path string true "スケジュールID"
// @Param userId path string true "User ID"
// @Success 200 {string} string "チャットルームが正常にハンドルされました"
// @Router /chat/room/{scheduleId}/{userId} [get]
// @Security Bearer
func (c *ChatController) HandleChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	userId := ctx.Param("userId")

	ctx.HTML(constants.StatusOK, "chat_room", gin.H{
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
// @Param scheduleId path string true "スケジュールID"
// @Success 200 {object} map[string]interface{} "Chat room created successfully."
// @Failure 400 {object} map[string]interface{} "Failed to create chat room."
// @Router /chat/create-room/{scheduleId} [post]
// @Security Bearer
func (c *ChatController) CreateChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	c.chatManager.CreateRoom(scheduleId)
	respondWithSuccess(ctx, constants.StatusOK, "Chat room created successfully.")
}

// PostToChatRoom godoc
// @Summary チャットルームに投稿
// @Description チャットルームにメッセージを投稿する。
// @Tags Chat Room
// @Accept multipart/form-data
// @Produce json
// @Param scheduleId path int true "スケジュールID"
// @Param user formData string true "ユーザーID"
// @Param message formData string true "メッセージ"
// @Success 200 {object} map[string]interface{} "Message posted successfully."
// @Router /chat/room/{scheduleId} [post]
// @Security Bearer
func (c *ChatController) PostToChatRoom(ctx *gin.Context) {
	user, message := ctx.PostForm("user"), ctx.PostForm("message")
	if user == "" || message == "" {
		respondWithError(ctx, constants.StatusBadRequest, "User and message must be provided.")
		return
	}
	scheduleId := ctx.Param("scheduleId")
	c.chatManager.Submit(user, scheduleId, message)
	respondWithSuccess(ctx, constants.StatusOK, "Message posted successfully.")
}

// DeleteChatRoom godoc
// @Summary チャットルームを削除
// @Description チャットルームを削除する。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param scheduleId path string true "スケジュールID"
// @Success 200 {object} string "Chat room deleted successfully."
// @Router /chat/room/{scheduleId} [delete]
// @Security Bearer
func (c *ChatController) DeleteChatRoom(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	c.chatManager.DeleteBroadcast(scheduleId)
	respondWithSuccess(ctx, constants.StatusOK, "Chat room deleted successfully.")
}

// StreamChat godoc
// @Summary チャットをストリーム
// @Description チャットをストリームする。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param scheduleId path int true "スケジュールID"
// @Router /chat/stream/{scheduleId} [get]
// @Security Bearer
func (c *ChatController) StreamChat(ctx *gin.Context) {
	scheduleId := ctx.Param("scheduleId")
	listener := c.chatManager.OpenListener(scheduleId)
	defer c.chatManager.CloseListener(scheduleId, listener)

	ctx.Stream(func(w io.Writer) bool {
		select {
		case message := <-listener:
			ctx.SSEvent("message", message)
			return true
		case <-ctx.Request.Context().Done():
			return false
		}
	})
}

// GetChatMessages godoc
// @Summary チャットメッセージを取得
// @Description チャットメッセージを取得する。
// @Tags Chat Room
// @Accept json
// @Produce json
// @Param roomid path string true "ルームID"
// @Success 200 {object} string "success"
// @Failure 404 {object} string "Chat room not found"
// @Router /chat/messages/{roomid} [get]
// @Security Bearer
func (c *ChatController) GetChatMessages(ctx *gin.Context) {
	roomid := ctx.Param("roomid")
	exists, err := c.redisClient.Exists(context.Background(), "chat:"+roomid).Result()
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Error checking room existence.")
		return
	}
	if exists == 0 {
		respondWithError(ctx, constants.StatusNotFound, "Chat room not found.")
		return
	}
	messages, err := c.redisClient.LRange(context.Background(), "chat:"+roomid, 0, -1).Result()
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Failed to load messages.")
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, messages)
}

// SendDirectMessage godoc
// @Summary DMを送信
// @Description 特定のユーザーにDMを送信
// @Tags Direct Message
// @Accept json
// @Produce json
// @Param senderId path string true "送信者ID"
// @Param receiverId path string true "受信者ID"
// @Param message formData string true "Message"
// @Success 200 {object} string "Message sent successfully"
// @Router /chat/dm/{senderId}/{receiverId} [post]
// @Security Bearer
func (c *ChatController) SendDirectMessage(ctx *gin.Context) {
	senderId, receiverId, message := ctx.Param("senderId"), ctx.Param("receiverId"), ctx.PostForm("message")
	if senderId == "" || receiverId == "" || message == "" {
		respondWithError(ctx, constants.StatusBadRequest, "Sender, receiver and message must be provided and non-empty.")
		return
	}
	if err := c.chatManager.SubmitDirectMessage(senderId, receiverId, message); err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Failed to send message.")
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, "Message sent successfully.")
}

// GetDirectMessages godoc
// @Summary DM履歴を取得
// @Description 特定のユーザー間のDM履歴を取得
// @Tags Direct Message
// @Accept json
// @Produce json
// @Param senderId path string true "送信者ID"
// @Param receiverId path string true "受信者ID"
// @Success 200 {object} string "Messages fetched successfully"
// @Router /chat/dm/{senderId}/{receiverId} [get]
// @Security Bearer
func (c *ChatController) GetDirectMessages(ctx *gin.Context) {
	senderId, receiverId := ctx.Param("senderId"), ctx.Param("receiverId")
	messages, err := c.chatManager.GetDirectMessages(senderId, receiverId)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Failed to fetch messages.")
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, messages)
}

// DeleteDirectMessages godoc
// @Summary DM履歴を削除
// @Description 特定のユーザー間のDM履歴を削除
// @Tags Direct Message
// @Accept json
// @Produce json
// @Param senderId path string true "送信者ID"
// @Param receiverId path string true "受信者ID"
// @Success 200 {object} string "Messages deleted successfully"
// @Router /chat/dm/{senderId}/{receiverId} [delete]
// @Security Bearer
func (c *ChatController) DeleteDirectMessages(ctx *gin.Context) {
	senderId, receiverId := ctx.Param("senderId"), ctx.Param("receiverId")
	if err := c.chatManager.DeleteDirectMessages(senderId, receiverId); err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Failed to delete messages.")
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, "Messages deleted successfully.")
}
