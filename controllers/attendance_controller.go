package controllers

import (
	"log"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Invalid request format"`
}

type SuccessResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

type AttendanceController struct {
	attendanceService services.AttendanceService
}

type AttendanceInput struct {
	UID    uint   `json:"uid"`
	CID    uint   `json:"cid"`
	CSID   uint   `json:"csid"`
	Status string `json:"status"`
}

type AttendanceResponse struct {
	Attendances []models.Attendance           `json:"attendances"`
	Statistics  map[models.AttendanceType]int `json:"statistics"`
}

func NewAttendanceController(service services.AttendanceService) *AttendanceController {
	return &AttendanceController{attendanceService: service}
}

// CreateOrUpdateAttendances godoc
// @Summary 出席情報を作成または更新
// @Description 出席情報を作成または更新する。
// @Tags Attendance
// @Security Bearer
// @Accept json
// @Produce json
// @Param attendance body []AttendanceInput true "Attendance records to create or update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendances [post]
func (ac *AttendanceController) CreateOrUpdateAttendances(ctx *gin.Context) {
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

// GetAttendancesByClass godoc
// @Summary 全ての出席情報を取得
// @Description 全ての出席情報を取得する。
// @Tags Attendance
// @Accept json
// @Produce json
// @Param cid path uint true "Class ID"
// @Param csid path uint true "Class Schedule ID"
// @Success 200 {object} AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendances/class/{classId} [get]
func (ac *AttendanceController) GetAttendancesByClass(ctx *gin.Context) {
	classId, err := strconv.ParseUint(ctx.Param("classId"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid class ID")
		return
	}

	attendances, err := ac.attendanceService.GetAllAttendancesByCID(uint(classId))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	if len(attendances) == 0 {
		respondWithError(ctx, constants.StatusNotFound, "No attendance records found")
		return
	}

	response := generateAttendanceResponse(attendances)
	respondWithSuccess(ctx, constants.StatusOK, response)
}

// GetAttendancesBySchedule godoc
// @Summary クラススケジュールIDで出席統計情報を取得
// @Description 指定されたクラススケジュールIDの出席統計情報を取得する。
// @Tags Attendance
// @Security Bearer
// @Accept json
// @Produce json
// @Param scheduleId path integer true "Class Schedule ID"
// @Success 200 {object} AttendanceResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendances/schedule/{scheduleId} [get]
func (ac *AttendanceController) GetAttendancesBySchedule(ctx *gin.Context) {
	scheduleId, err := strconv.ParseUint(ctx.Param("scheduleId"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid schedule ID")
		return
	}

	attendances, err := ac.attendanceService.GetAllAttendancesByCSID(uint(scheduleId))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	if len(attendances) == 0 {
		respondWithError(ctx, constants.StatusNotFound, "No attendance records found")
		return
	}

	response := generateAttendanceResponse(attendances)
	respondWithSuccess(ctx, constants.StatusOK, response)
}

// DeleteAttendance godoc
// @Summary IDで出席情報を削除
// @Description 指定されたIDの出席情報を削除する。
// @Tags Attendance
// @Security Bearer
// @Accept json
// @Produce json
// @Param id path integer true "Attendance ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /attendances/{id} [delete]
func (ac *AttendanceController) DeleteAttendance(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid attendance ID")
		return
	}

	if err := ac.attendanceService.DeleteAttendance(uint(id)); err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, map[string]string{"message": "Attendance deleted successfully"})
}

func generateAttendanceResponse(attendances []models.Attendance) AttendanceResponse {
	statistics := map[models.AttendanceType]int{
		models.AttendanceStatus: 0,
		models.TardyStatus:      0,
		models.AbsenceStatus:    0,
	}

	for _, attendance := range attendances {
		statistics[attendance.IsAttendance]++
	}

	return AttendanceResponse{
		Attendances: attendances,
		Statistics:  statistics,
	}
}
