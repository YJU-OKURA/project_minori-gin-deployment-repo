package controllers

import (
	"log"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	attendanceService services.AttendanceService
}

type AttendanceInput struct {
	UID    uint   `json:"uid"`
	CID    uint   `json:"cid"`
	CSID   uint   `json:"csid"`
	Status string `json:"status"`
}

func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{attendanceService: service}
}

// CreateOrUpdateAttendance godoc
// @Summary 出席情報を作成または更新
// @Description 出席情報を作成または更新する。
// @Tags Attendance
// @Security ApiKeyAuth
// @CrossOrigin
// @Accept json
// @Produce json
// @Param attendance body []AttendanceInput true "Attendance records to create or update"
// @Success 200 {object} map[string]string "出席情報が正常に作成または更新されました"
// @Failure 400 {object} map[string]string "リクエストが不正です"
// @Failure 500 {object} map[string]string "サーバーエラーが発生しました"
// @Router /at [post]
func (ac *AttendanceController) CreateOrUpdateAttendance(ctx *gin.Context) {
	var attendances []AttendanceInput
	if err := ctx.ShouldBindJSON(&attendances); err != nil {
		log.Printf("Error binding JSON: %v", err)
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	for _, attendance := range attendances {
		if attendance.Status != string(models.AttendanceStatus) && attendance.Status != string(models.TardyStatus) && attendance.Status != string(models.AbsenceStatus) {
			log.Printf("Invalid attendance status: %s", attendance.Status)
			respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
			return
		}

		err := ac.attendanceService.CreateOrUpdateAttendance(attendance.CID, attendance.UID, attendance.CSID, attendance.Status)
		if err != nil {
			log.Printf("Error creating or updating attendance: %v", err)
			handleServiceError(ctx, err)
			return
		}
	}

	respondWithSuccess(ctx, constants.StatusOK, map[string]string{"message": constants.Success})
}

// GetAllAttendances godoc
// @Summary 全ての出席情報を取得
// @Description 全ての出席情報を取得する。
// @Tags Attendance
// @Accept json
// @Produce json
// @Param cid path uint true "Class ID"
// @Param csid path uint true "Class Schedule ID"
// @Success 200 {object} []models.Attendance "出席情報が見つかりました"
// @Failure 400 {object} string "無効なクラスIDです"
// @Failure 404 {object} string "出席情報が見つかりません"
// @Router /attendance/{cid} [get]
func (ac *AttendanceController) GetAllAttendances(ctx *gin.Context) {
	cid, cidErr := strconv.ParseUint(ctx.Query("cid"), 10, 32)
	csid, csidErr := strconv.ParseUint(ctx.Query("csid"), 10, 32)

	var attendances []models.Attendance
	var err error

	if cidErr == nil {
		attendances, err = ac.attendanceService.GetAllAttendancesByCID(uint(cid))
	} else if csidErr == nil {
		attendances, err = ac.attendanceService.GetAllAttendancesByCSID(uint(csid))
	} else {
		log.Printf("Invalid ID: %v, %v", cidErr, csidErr)
		handleServiceError(ctx, err)
		return
	}

	if err != nil {
		log.Printf("Error retrieving attendances: %v", err)
		handleServiceError(ctx, err)
		return
	}

	if len(attendances) == 0 {
		respondWithError(ctx, constants.StatusNotFound, "No attendance found")
		return
	}

	statistics := make(map[models.AttendanceType]int)
	for _, attendance := range attendances {
		statistics[attendance.IsAttendance]++
	}

	response := gin.H{
		"attendances": attendances,
		"statistics":  statistics,
	}

	respondWithSuccess(ctx, constants.StatusOK, response)
}

// DeleteAttendance godoc
// @Summary IDで出席情報を削除
// @Description 指定されたIDの出席情報を削除する。
// @Tags Attendance
// @Accept json
// @Produce json
// @Param id path uint true "Attendance ID"
// @Success 200 {object} map[string]string "出席情報が正常に削除されました"
// @Failure 400 {object} string "無効なID形式です"
// @Failure 404 {object} string "出席情報が見つかりません"
// @Router /attendance/{id} [delete]
func (ac *AttendanceController) DeleteAttendance(ctx *gin.Context) {
	attendanceID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	if err := ac.attendanceService.DeleteAttendance(uint(attendanceID)); err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, "Attendance deleted successfully")
}

// GetAttendanceStatisticsByCSID godoc
// @Summary クラススケジュールIDで出席統計情報を取得
// @Description 指定されたクラススケジュールIDの出席統計情報を取得する。
// @Tags Attendance
// @Accept json
// @Produce json
// @Param csid path uint true "Class Schedule ID"
// @Success 200 {object} map[string]int "出席統計情報が見つかりました"
// @Failure 400 {object} string "無効なクラススケジュールIDです"
// @Failure 404 {object} string "出席統計情報が見つかりません"
// @Router /attendance/statistics/schedule/{csid} [get]
func (ac *AttendanceController) GetAttendanceStatisticsByCSID(ctx *gin.Context) {
	classScheduleID, err := strconv.ParseUint(ctx.Param("csid"), 10, 32)
	if err != nil {
		log.Printf("Invalid classScheduleID: %v", err)
		ctx.JSON(400, gin.H{"error": "Invalid class schedule ID"})
		return
	}

	attendances, err := ac.attendanceService.GetAllAttendancesByCSID(uint(classScheduleID))
	if err != nil {
		log.Printf("Error retrieving statistics: %v", err)
		ctx.JSON(500, gin.H{"error": "Server error"})
		return
	}

	if len(attendances) == 0 {
		ctx.JSON(404, gin.H{"error": "No statistics found"})
		return
	}

	statistics := make(map[models.AttendanceType]int)
	for _, attendance := range attendances {
		statistics[attendance.IsAttendance]++
	}

	response := gin.H{
		"attendances": attendances,
		"statistics":  statistics,
	}

	ctx.JSON(200, response)
}
