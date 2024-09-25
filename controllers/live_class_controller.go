package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
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
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 64)
	info, err := ctrl.liveClassService.GetScreenShareInfo(c.Request.Context(), uint(cid))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func (ctrl *LiveClassController) StartScreenShare(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 64)

	// 스트리밍 서비스 API 호출을 통해 실제 스트리밍 URL 생성
	streamURL, err := ctrl.liveClassService.StartStreamingSession(uint(cid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start streaming session: " + err.Error()})
		return
	}

	// 스트리밍 정보를 저장
	info := map[string]interface{}{
		"streamURL": streamURL,
		"startedAt": time.Now(),
	}

	err = ctrl.liveClassService.SaveScreenShareInfo(c.Request.Context(), uint(cid), info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save streaming info: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Screen sharing started successfully", "streamURL": streamURL})
}
