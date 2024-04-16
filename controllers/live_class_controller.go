package controllers

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
)

// LiveClassController implements the interface
type LiveClassController struct {
	liveClassService services.LiveClassService
	upgrader         websocket.Upgrader
}

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
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"roomID": roomID})
	}
}

// StartScreenShareHandler godoc
// @Summary 画面共有を開始します。
// @Description 画面共有を開始します。
// @Tags Live Class
// @Produce json
// @Param Authorization header string true "Bearer Token"
// @Param roomID path string true "ルームID"
// @Success 200 {object} ScreenShareResponse
// @Failure 400 {object} ErrorResponse
// @Router /live/start-screen-share/{roomID} [get]
func (c *LiveClassController) StartScreenShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the JWT from the Authorization header
		authHeader := ctx.GetHeader("Authorization")

		if !authenticateUser(authHeader) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		roomID := ctx.Param("roomID")
		pc, err := c.liveClassService.StartScreenShare(roomID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		offer, err := pc.CreateOffer(nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = pc.SetLocalDescription(offer)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"sdp": pc.LocalDescription().SDP})
		log.Printf("User %s started screen share in room %s", ctx.MustGet("userID"), roomID)
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
// @Success 200 {object} StandardResponse
// @Failure 400 {object} ErrorResponse
// @Router /live/stop-screen-share/{roomID} [get]
func (c *LiveClassController) StopScreenShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("roomID")
		err := c.liveClassService.StopScreenShare(roomID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Screen share stopped successfully"})
	}
}

func authenticateUser(tokenString string) bool {
	// Assuming the token is in the Authorization header as a Bearer token
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

type RoomResponse struct {
	RoomID string `json:"roomID"`
}

type ScreenShareResponse struct {
	SDP string `json:"sdp"`
}

type StandardResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
