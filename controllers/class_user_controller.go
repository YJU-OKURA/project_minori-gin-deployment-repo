package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ClassUserController インタフェースを実装
type ClassUserController struct {
	classUserService services.ClassUserService
}

// NewClassUserController ClassScheduleControllerを生成
func NewClassUserController(service services.ClassUserService) *ClassUserController {
	return &ClassUserController{
		classUserService: service,
	}
}

// ChangeUserRole godoc
// @Summary ユーザーのロールを変更します。
// @Description ユーザーのロールを変更します。
// @Tags Class User
// @ID change-user-role
// @Accept  json
// @Produce  json
// @Param uid path int true "User ID"
// @Param cid path int true "Class ID"
// @Param role path string true "Role Name"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/{role} [patch]
func (c *ClassUserController) ChangeUserRole(ctx *gin.Context) {
	uid, _ := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	cid, _ := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	roleName := ctx.Param("role")

	err := c.classUserService.AssignRole(uint(uid), uint(cid), roleName)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}
