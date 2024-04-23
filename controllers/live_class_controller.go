package controllers

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type LiveClassController struct {
	liveClassService services.LiveClassService
}

func NewLiveClassController(liveClassService services.LiveClassService) *LiveClassController {
	return &LiveClassController{
		liveClassService: liveClassService,
	}
}

// GetScreenShareInfo godoc
// @Summary スクリーン共有情報を取得
// @Description 特定のクラスのスクリーン共有情報を取得します。ユーザーがそのクラスのメンバーである必要があります。
// @Tags Live Class
// @Accept  json
// @Produce  json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /live/screen_share/{uid}/{cid} [get]
func (ctrl *LiveClassController) GetScreenShareInfo(c *gin.Context) {
	uid, _ := strconv.ParseUint(c.Param("uid"), 10, 64)
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 64)

	info, err := ctrl.liveClassService.GetScreenShareInfo(c, uint(uid), uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// StartScreenShare godoc
// @Summary スクリーン共有を開始
// @Description 特定のクラスのスクリーン共有を開始します。
// @Tags Live Class
// @Accept  json
// @Produce  json
// @Param cid path int true "クラスID"
// @Success 200 {object} map[string]interface{} "Message and stream URL indicating the screen sharing has started successfully"
// @Failure 500 {object} map[string]string "Internal server error with error message"
// @Router /live/screen_share/start/{cid} [post]
func (ctrl *LiveClassController) StartScreenShare(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 64)

	// Simulate generating a unique stream URL using the class ID
	streamURL := fmt.Sprintf("https://streaming.service.com/live/%d/%d", cid, time.Now().Unix())

	info := map[string]interface{}{
		"streamURL": streamURL,
		"startedAt": time.Now(),
	}

	err := ctrl.liveClassService.SaveScreenShareInfo(c.Request.Context(), uint(cid), info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Screen sharing started successfully", "streamURL": streamURL})
}
