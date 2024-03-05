package controllers

import (
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

// GetAllAttendances godoc
// @Summary クラスの全ての出席情報を取得
// @Description クラスの全ての出席情報を取得
// @Tags attendance
// @Accept json
// @Produce json
// @Param cid path int true "Class ID"
// @Success 200 {array} models.Attendance "List of attendances"
// @Failure 500 {string} string "Server error"
// @Router /c/{cid} [get]
func (controller *AttendanceController) GetAllAttendances(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 32)
	attendances, err := controller.attendanceService.GetAllAttendancesByCID(uint(cid))
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constants.StatusOK, attendances)
}

// GetAttendance godoc
// @Summary 出席情報を取得
// @Description 指定されたIDの出席情報を取得
// @Tags attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} models.Attendance "Attendance"
// @Failure 500 {string} string "Server error"
// @Router /attendance/{id} [get]
func (controller *AttendanceController) GetAttendance(c *gin.Context) {
	id := c.Param("id")
	attendance, err := controller.attendanceService.GetAttendanceByID(id)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constants.StatusOK, attendance)
}

// CreateOrUpdateAttendance godoc
// @Summary 出席情報を作成または更新
// @Description 出席情報を作成または更新
// @Tags attendance
// @Accept json
// @Produce json
// @Param cid path int true "Class ID"
// @Param uid path int true "User ID"
// @Param status query string true "Status"
// @Success 200 {string} string "success"
// @Failure 500 {string} string "Server error"
// @Router /attendance/{cid}/{uid} [post]
func (controller *AttendanceController) CreateOrUpdateAttendance(c *gin.Context) {
	cid, _ := strconv.ParseUint(c.Param("cid"), 10, 32)
	uid, _ := strconv.ParseUint(c.Param("uid"), 10, 32)
	status := c.Query("status")

	err := controller.attendanceService.CreateOrUpdateAttendance(uint(cid), uint(uid), status)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constants.StatusOK, gin.H{"message": "success"})
}

// DeleteAttendance godoc
// @Summary 出席情報を削除
// @Description 指定されたIDの出席情報を削除
// @Tags attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {string} string "Attendance deleted successfully"
// @Failure 500 {string} string "Server error"
// @Router /attendance/{id} [delete]
func (controller *AttendanceController) DeleteAttendance(c *gin.Context) {
	id := c.Param("id")
	err := controller.attendanceService.DeleteAttendance(id)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(constants.StatusOK, gin.H{"message": "Attendance deleted successfully"})
}
