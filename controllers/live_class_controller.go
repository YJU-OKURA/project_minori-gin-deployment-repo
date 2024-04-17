package controllers

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
)

// LiveClassController handles all web socket operations for live classroom interactions
// including creating a room, starting and stopping screen sharing
type LiveClassController struct {
	liveClassService services.LiveClassService
	upgrader         websocket.Upgrader
}

// RoomResponse encapsulates the response structure for room creation.
type RoomResponse struct {
	RoomID string `json:"roomID"`
}

// ScreenShareResponse contains the SDP information necessary for establishing
// a WebRTC connection for screen sharing.
type ScreenShareResponse struct {
	SDP string `json:"sdp"`
}

// StandardResponse provides a generic response structure for simple messages.
type StandardResponse struct {
	Message string `json:"message"`
}

// ErrorResponse provides a structured error message for API responses.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// NewLiveClassController creates a new controller instance with the necessary dependencies.
func NewLiveClassController(service services.LiveClassService) *LiveClassController {
	return &LiveClassController{
		liveClassService: service,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// CreateRoomHandler godoc
// @Summary ルームを生成します。
// @Description ルームを生成します。
// @Tags Live Class
// @Produce json
// @Success 200 {object} RoomResponse
// @Failure 400 {object} ErrorResponse
// @Router /live/create-room [post]
func (c *LiveClassController) CreateRoomHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID, err := c.liveClassService.CreateRoom()
		if err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Failed to create room"+err.Error())
			return
		}
		respondWithSuccess(ctx, constants.StatusOK, RoomResponse{RoomID: roomID})
	}
}

// StartScreenShareHandler godoc
// @Summary 画面共有を開始します。
// @Description 画面共有を開始します。
// @Tags Live Class
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param roomID path string true "ルームID"
// @Param userID path string true "ユーザーID"
// @Success 200 {object} ScreenShareResponse
// @Failure 400 {object} ErrorResponse
// @Router /live/start-screen-share/{roomID} [get]
func (c *LiveClassController) StartScreenShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !authenticateUser(authHeader) {
			respondWithError(ctx, constants.StatusUnauthorized, "Unauthorized access")
			return
		}

		roomID, userID := ctx.Param("roomID"), ctx.Param("userID")
		pc, err := c.liveClassService.StartScreenShare(roomID, userID)
		if err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Screen sharing could not be started: "+err.Error())
			return
		}

		offer, err := pc.CreateOffer(nil)
		if err != nil {
			pc.Close()
			respondWithError(ctx, constants.StatusInternalServerError, "Failed to create offer: "+err.Error())
			return
		}

		if err := pc.SetLocalDescription(offer); err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Failed to set local description: "+err.Error())
			return
		}

		respondWithSuccess(ctx, constants.StatusOK, ScreenShareResponse{SDP: pc.LocalDescription().SDP})
		log.Printf("Screen sharing started by user %s in room %s", userID, roomID)
	}
}

// StopScreenShareHandler godoc
// @Summary 画面共有を停止します。
// @Description 画面共有を停止します。
// @Tags Live Class
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param roomID path string true "ルームID"
// @Param userID path string true "ユーザーID"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Router /live/stop-screen-share/{roomID} [get]
func (c *LiveClassController) StopScreenShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("roomID")
		userID := ctx.Param("userID")
		err := c.liveClassService.StopScreenShare(roomID, userID)
		if err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Screen sharing could not be stopped: "+err.Error())
			return
		}
		respondWithSuccess(ctx, constants.StatusOK, StandardResponse{Message: "Screen share stopped successfully"})
	}
}

// ViewScreenShareHandler godoc
// @Summary 画面共有情報を取得します。
// @Description 画面共有情報を取得します。
// @Tags Live Class
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Success 200 {object} ScreenShareResponse "SDP information"
// @Failure 400 {object} ErrorResponse "Error message and details"
// @Failure 401 "Unauthorized if the user is not authenticated or not part of the room"
// @Router /live/view-screen-share/{roomID} [get]
func (c *LiveClassController) ViewScreenShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("roomID")
		userID, exists := ctx.Get("userID")
		if !exists {
			respondWithError(ctx, constants.StatusUnauthorized, "User ID not provided")
			return
		}

		if !c.liveClassService.IsUserInRoom(userID.(string), roomID) {
			respondWithError(ctx, constants.StatusUnauthorized, "Access denied")
			return
		}

		sdp, err := c.liveClassService.GetScreenShareInfo(roomID)
		if err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Failed to retrieve screen share info: "+err.Error())
			return
		}

		respondWithSuccess(ctx, constants.StatusOK, ScreenShareResponse{SDP: sdp})
	}
}

// authenticateUser checks if the provided JWT token is valid and authorized to access the system.
func authenticateUser(tokenString string) bool {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("Failed to authenticate user: %v", err)
		return false
	}

	return token.Valid
}
