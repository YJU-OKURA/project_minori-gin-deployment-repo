package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ClassScheduleController インタフェースを実装
type ClassScheduleController struct {
	classScheduleService services.ClassScheduleService
}

// NewClassScheduleController ClassScheduleControllerを生成
func NewClassScheduleController(service services.ClassScheduleService) *ClassScheduleController {
	return &ClassScheduleController{
		classScheduleService: service,
	}
}

// CreateClassSchedule godoc
// @Summary クラススケジュールを作成
// @Description 新しいクラススケジュールを作成する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param classSchedule body dto.ClassScheduleDTO true "Class schedule to create"
// @Success 200 {object} models.ClassSchedule "クラススケジュールが正常に作成されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs [post]
func (controller *ClassScheduleController) CreateClassSchedule(c *gin.Context) {
	var dto dto.ClassScheduleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classSchedule := models.ClassSchedule{
		Title:     dto.Title,
		StartedAt: dto.StartedAt,
		EndedAt:   dto.EndedAt,
		CID:       dto.CID,
		IsLive:    dto.IsLive,
	}

	createdClassSchedule, err := controller.classScheduleService.CreateClassSchedule(&classSchedule)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	respondWithSuccess(c, constants.StatusOK, createdClassSchedule)
}

// GetClassScheduleByID godoc
// @Summary IDでクラススケジュールを取得
// @Description 指定されたIDのクラススケジュールを取得する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param id path int true "Class schedule ID"
// @Success 200 {object} models.ClassSchedule "クラススケジュールが見つかりました"
// @Failure 400 {object} string "無効なID形式です"
// @Failure 404 {object} string "クラススケジュールが見つかりません"
// @Router /cs/{id} [get]
func (controller *ClassScheduleController) GetClassScheduleByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classSchedule, err := controller.classScheduleService.GetClassScheduleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": constants.ClassNotFound})
		return
	}

	respondWithSuccess(c, constants.StatusOK, classSchedule)
}

// GetAllClassSchedules godoc
// @Summary 全てのクラススケジュールを取得
// @Description 指定されたクラスIDの全てのクラススケジュールを取得する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param cid query uint false "Class ID"
// @Success 200 {array} []models.ClassSchedule "クラススケジュールが見つかりました"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs [get]
func (controller *ClassScheduleController) GetAllClassSchedules(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.DefaultQuery("cid", "0"), 10, 32)
	classSchedules, err := controller.classScheduleService.GetAllClassSchedules(uint(cid))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	respondWithSuccess(c, constants.StatusOK, classSchedules)
}

// UpdateClassSchedule godoc
// @Summary クラススケジュールを更新
// @Description 指定されたIDのクラススケジュールを更新する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param id path int true "Class schedule ID"
// @Param classSchedule body dto.UpdateClassScheduleDTO true "Class schedule to update"
// @Success 200 {object} models.ClassSchedule "クラススケジュールが正常に更新されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs/{id} [put]
func (controller *ClassScheduleController) UpdateClassSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	var dto dto.UpdateClassScheduleDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	updatedClassSchedule, err := controller.classScheduleService.UpdateClassSchedule(uint(id), &dto)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	respondWithSuccess(c, constants.StatusOK, updatedClassSchedule)
}

// DeleteClassSchedule godoc
// @Summary クラススケジュールを削除
// @Description 指定されたIDのクラススケジュールを削除する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param id path int true "Class schedule ID"
// @Success 200 {object} string "クラススケジュールが正常に削除されました"
// @Failure 400 {object} string "無効なID形式です"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs/{id} [delete]
func (controller *ClassScheduleController) DeleteClassSchedule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	err = controller.classScheduleService.DeleteClassSchedule(uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	respondWithSuccess(c, constants.StatusOK, constants.DeleteSuccess)
}

// GetLiveClassSchedules godoc
// @Summary ライブ中のクラススケジュールを取得
// @Description 指定されたクラスIDのライブ中のクラススケジュールを取得する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param cid query uint true "Class ID"
// @Success 200 {array} []models.ClassSchedule "ライブ中のクラススケジュールが見つかりました"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs/live [get]
func (controller *ClassScheduleController) GetLiveClassSchedules(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Query("cid"), 10, 32)
	classSchedules, err := controller.classScheduleService.GetLiveClassSchedules(uint(cid))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	respondWithSuccess(c, constants.StatusOK, classSchedules)
}

// GetClassSchedulesByDate godoc
// @Summary 日付でクラススケジュールを取得
// @Description 指定されたクラスIDと日付のクラススケジュールを取得する。
// @Tags Class Schedule
// @Accept json
// @Produce json
// @Param cid query uint true "Class ID"
// @Param date query string true "Date"
// @Success 200 {array} []models.ClassSchedule "指定された日付のクラススケジュールが見つかりました"
// @Failure 400 {object} string "日付が必要です"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cs/date [get]
func (controller *ClassScheduleController) GetClassSchedulesByDate(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Query("cid"), 10, 32)
	date := c.Query("date") // Expecting date in the format 'YYYY-MM-DD'

	if date == "" {
		respondWithError(c, constants.StatusBadRequest, constants.ErrNoDateJP)
		return
	}

	classSchedules, err := controller.classScheduleService.GetClassSchedulesByDate(uint(cid), date)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	respondWithSuccess(c, constants.StatusOK, classSchedules)
}
