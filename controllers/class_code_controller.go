package controllers

import (
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type ClassCodeController struct {
	classCodeService services.ClassCodeService
	classUserService services.ClassUserService
}

func NewClassCodeController(classCodeService services.ClassCodeService, classUserService services.ClassUserService) *ClassCodeController {
	return &ClassCodeController{
		classCodeService: classCodeService,
		classUserService: classUserService,
	}
}

// CheckSecretExists godoc
// @Summary グループコードにシークレットが存在するかチェック
// @Description 指定されたグループコードにシークレットがあるかどうかをチェックする。
// @Tags Class Code
// @Accept json
// @Produce json
// @Param code query string true "Code to check"
// @Success 200 {object} bool "secretExists" "シークレットが存在します"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Router /cc/checkSecretExists [get]
// @Security Bearer
func (c *ClassCodeController) CheckSecretExists(ctx *gin.Context) {
	code := ctx.Query("code")

	secretExists, err := c.classCodeService.CheckSecretExists(ctx, code)
	if err != nil {
		// エラーメッセージに基づいて適切なHTTPステータスを返す
		if err.Error() == services.ErrClassNotFound {
			respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
			return
		}
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"secretExists": secretExists})
}

// VerifyClassCode godoc
// @Summary グループコードとシークレットを検証＆ユーザーに役割を割り当てる
// @Description グループコードと、該当する場合はそのシークレットを確認する。また、指定されたユーザーに役割を割り当てる。
// @Tags Class Code
// @Accept json
// @Produce json
// @Param code query string true "Code to verify"
// @Param secret query string false "Secret for the code"
// @Param uid query int true "User ID to assign role"
// @Success 200 {object} string "グループコードが検証されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 401 {object} string "シークレットが一致しません"
// @Failure 404 {object} string "コードが見つかりません"
// @Router /cc/verifyClassCode [get]
// @Security Bearer
func (c *ClassCodeController) VerifyClassCode(ctx *gin.Context) {
	code := ctx.Query("code")
	secret := ctx.Query("secret")
	uidStr := ctx.Query("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	isValid, err := c.classCodeService.VerifyClassCode(code, secret)
	if err != nil {
		if err.Error() == services.ErrClassNotFound {
			respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
			return
		}
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}

	if !isValid {
		respondWithSuccess(ctx, constants.StatusOK, gin.H{"valid": false})
		return
	}

	roleName := "APPLICANT"
	cid := uint(uid)
	err = c.classUserService.AssignRole(uint(uid), cid, roleName)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.AssignError)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"valid": true, "message": constants.ClassMemberRegistration})
}

func parseUintQueryParam(ctx *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(ctx.Query(param), 10, 64)
}
