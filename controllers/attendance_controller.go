package controllers

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

// AttendanceController インタフェースを実装
type AttendanceController struct {
	attendanceService services.AttendanceService
}

// NewAttendanceController AttendanceControllerを生成
func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{
		attendanceService: service,
	}
}

// CreateOrUpdateAttendance godoc
// @Summary 出席情報を作成または更新
// @Description 出席情報を作成または更新
// @Tags attendance
// @Accept json
// @Produce json
// @Param cid path int true "Class ID"
// @Param uid path int true "User ID"
// @Param csid path int true "Class Schedule ID"
// @Param status query string true "Status"
// @Success 200 {string} string "作成または更新に成功しました"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/{cid}/{uid}/{csid} [post]
func (ac *AttendanceController) CreateOrUpdateAttendance(ctx *gin.Context) {
	classID, classIDErr := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	userID, userIDErr := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	scheduleID, scheduleIDErr := strconv.ParseUint(ctx.Param("csid"), 10, 32)

	if classIDErr != nil || userIDErr != nil || scheduleIDErr != nil {
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}

	status := ctx.Query("status")
	err := ac.attendanceService.CreateOrUpdateAttendance(uint(classID), uint(userID), uint(scheduleID), status)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.CreateOrUpdateSuccess})
}

// GetAllAttendances godoc
// @Summary クラスの全ての出席情報を取得
// @Description クラスの全ての出席情報を取得
// @Tags attendance
// @Accept json
// @Produce json
// @Param classID path int true "Class ID"
// @Success 200 {array} models.Attendance "Attendance"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/{classID} [get]
func (ac *AttendanceController) GetAllAttendances(ctx *gin.Context) {
	classID, err := strconv.ParseUint(ctx.Param("classID"), 10, 32)
	if err != nil {
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}

	attendances, serviceErr := ac.attendanceService.GetAllAttendancesByCID(uint(classID))
	if serviceErr != nil {
		handleServiceError(ctx, serviceErr)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, attendances)
}

// GetAttendance godoc
// @Summary 出席情報を取得
// @Description 指定されたIDの出席情報を取得
// @Tags attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} models.Attendance "Attendance"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/attendance/{id} [get]
func (ac *AttendanceController) GetAttendance(ctx *gin.Context) {
	attendanceID := ctx.Param("id")
	attendance, err := ac.attendanceService.GetAttendanceByID(attendanceID)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, attendance)
}

// DeleteAttendance godoc
// @Summary 出席情報を削除
// @Description 指定されたIDの出席情報を削除
// @Tags attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {string} string "削除に成功しました"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/attendance/{id} [delete]
func (ac *AttendanceController) DeleteAttendance(ctx *gin.Context) {
	attendanceID := ctx.Param("id")
	if err := ac.attendanceService.DeleteAttendance(attendanceID); err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.DeleteSuccess})
}
