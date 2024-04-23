package controllers

//
//import (
//	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strconv"
//)
//
//type LiveClassController struct {
//	Service services.LiveClassService
//}
//
//func NewLiveClassController(service services.LiveClassService) *LiveClassController {
//	return &LiveClassController{
//		Service: service,
//	}
//}
//
//// CreateRoom godoc
//// @Summary 新しいルームを作成します。
//// @Description 新しいルームを作成します。
//// @Tags Live Class
//// @Accept json
//// @Produce json
//// @Param   classID path uint true "Class ID"
//// @Param   userID path uint true "User ID"
//// @Success 200 {object} map[string]interface{} "roomID returned on successful creation"
//// @Failure 400 {object} map[string]interface{} "Invalid class ID"
//// @Failure 401 {object} map[string]interface{} "Unauthorized to create room"
//// @Failure 500 {object} map[string]interface{} "Internal server error"
//// @Router /live/create-room/{classID}/{userID} [post]
//func (c *LiveClassController) CreateRoom(ctx *gin.Context) {
//	userID, _ := strconv.ParseUint(ctx.Param("userID"), 10, 32)
//	classID, err := strconv.ParseUint(ctx.Param("classID"), 10, 32)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid classID"})
//		return
//	}
//
//	roomID, err := c.Service.CreateRoom(uint(classID), uint(userID))
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{"roomID": roomID})
//}
//
//// StartScreenShare godoc
//// @Summary 画面共有を開始します。
//// @Description 画面共有を開始します。
//// @Tags Live Class
//// @Accept  json
//// @Produce  json
//// @Param   roomID path string true "Room ID"
//// @Param   userID path uint true "User ID"
//// @Success 200 {object} map[string]interface{} "SDP data for the screen share"
//// @Failure 400 {object} map[string]interface{} "Invalid room ID"
//// @Failure 401 {object} map[string]interface{} "Unauthorized to start screen sharing"
//// @Failure 500 {object} map[string]interface{} "Internal server error"
//// @Router /live/start-screen-share/{roomID}/{userID} [post]
//func (c *LiveClassController) StartScreenShare(ctx *gin.Context) {
//	roomID := ctx.Param("roomID")
//	userID := ctx.GetString("userID")
//	err := c.Service.StartScreenShare(roomID, userID)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	sdp, err := c.Service.GetScreenShareSDP(roomID)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{"sdp": sdp})
//}
//
//// StopScreenShare godoc
//// @Summary 画面共有を停止します。
//// @Description 画面共有を停止します。
//// @Tags Live Class
//// @Accept  json
//// @Produce  json
//// @Param   roomID path string true "Room ID"
//// @Param   userID path uint true "User ID"
//// @Success 200 {object} map[string]interface{} "Screen sharing stopped successfully"
//// @Failure 400 {object} map[string]interface{} "Invalid room ID"
//// @Failure 401 {object} map[string]interface{} "Unauthorized to stop screen sharing"
//// @Failure 500 {object} map[string]interface{} "Internal server error"
//// @Router /live/stop-screen-share/{roomID}/{userID} [post]
//func (c *LiveClassController) StopScreenShare(ctx *gin.Context) {
//	roomID := ctx.Param("roomID")
//	adminID := ctx.GetString("adminID")
//	err := c.Service.StopScreenShare(roomID, adminID)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{"message": "Screen sharing stopped"})
//}
//
//// JoinScreenShare godoc
//// @Summary 画面共有に参加します。
//// @Description 画面共有に参加します。
//// @Tags Live Class
//// @Accept  json
//// @Produce  json
//// @Param   roomID path string true "Room ID"
//// @Param   userID path uint true "User ID"
//// @Success 200 {object} map[string]interface{} "SDP data for the screen share"
//// @Failure 400 {object} map[string]interface{} "Invalid room ID or User ID"
//// @Failure 401 {object} map[string]interface{} "Unauthorized to join screen sharing"
//// @Failure 500 {object} map[string]interface{} "Internal server error"
//// @Router /live/join-screen-share/{roomID}/{userID} [get]
//func (c *LiveClassController) JoinScreenShare(ctx *gin.Context) {
//	roomID := ctx.Param("roomID")
//	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID"})
//		return
//	}
//
//	offer, err := c.Service.JoinScreenShare(roomID, uint(userID))
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	ctx.JSON(http.StatusOK, gin.H{"offer": offer})
//}
