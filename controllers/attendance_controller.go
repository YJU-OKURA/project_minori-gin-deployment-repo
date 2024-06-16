package controllers

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// AttendanceController インタフェースを実装
type AttendanceController struct {
	attendanceService services.AttendanceService
}

type AttendanceInput struct {
	UID    uint   `json:"uid"`
	CID    uint   `json:"cid"`
	CSID   uint   `json:"csid"`
	Status string `json:"status"`
}

// NewAttendanceController AttendanceControllerを生成
func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{
		attendanceService: service,
	}
}

// CreateOrUpdateAttendance godoc
// @Summary 複数の出席情報を作成または更新
// @Description 複数の出席情報を作成または更新します。'ATTENDANCE', 'TARDY', 'ABSENCE'のいずれかのステータスを持つことができます。
// @Tags Attendance
// @Accept json
// @Produce json
// @Param attendances body []AttendanceInput true "出席情報"
// @Success 200 {string} string "作成または更新に成功しました"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at [post]
// @Security Bearer
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

	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}

// GetAllAttendances godoc
// @Summary クラスの全ての出席情報を取得
// @Description クラスの全ての出席情報を取得
// @Tags Attendance
// @Accept json
// @Produce json
// @Param cid path int true "Class ID"
// @Success 200 {array} models.Attendance "Attendance"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/{cid} [get]
// @Security Bearer
func (ac *AttendanceController) GetAllAttendances(ctx *gin.Context) {
	log.Println("GetAllAttendances: Request received")

	classID, err := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if err != nil {
		log.Printf("GetAllAttendances: Invalid classID: %v", err)
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}
	log.Printf("GetAllAttendances: Parsed classID: %d", classID)

	attendances, serviceErr := ac.attendanceService.GetAllAttendancesByCID(uint(classID))
	if serviceErr != nil {
		log.Printf("GetAllAttendances: Error retrieving attendances: %v", serviceErr)
		handleServiceError(ctx, serviceErr)
		return
	}

	if len(attendances) == 0 {
		log.Println("GetAllAttendances: No attendances found")
		respondWithError(ctx, constants.StatusNotFound, "No attendance found")
		return
	}
	log.Printf("GetAllAttendances: Found %d attendances", len(attendances))
	respondWithSuccess(ctx, constants.StatusOK, attendances)
}

// GetAttendance godoc
// @Summary 出席情報を取得
// @Description 指定されたIDの出席情報を取得
// @Tags Attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} models.Attendance "Attendance"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/attendance/{id} [get]
// @Security Bearer
func (ac *AttendanceController) GetAttendance(ctx *gin.Context) {
	attendanceID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}

	attendance, err := ac.attendanceService.GetAttendanceByID(strconv.Itoa(int(attendanceID)))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	if attendance == nil {
		respondWithError(ctx, constants.StatusNotFound, "Attendance not found")
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, attendance)
}

// DeleteAttendance godoc
// @Summary 出席情報を削除
// @Description 指定されたIDの出席情報を削除
// @Tags Attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {string} string "削除に成功しました"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /at/attendance/{id} [delete]
// @Security Bearer
func (ac *AttendanceController) DeleteAttendance(ctx *gin.Context) {
	attendanceID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}

	if err := ac.attendanceService.DeleteAttendance(strconv.Itoa(int(attendanceID))); err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.DeleteSuccess})
}
