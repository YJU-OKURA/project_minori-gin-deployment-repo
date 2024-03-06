package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"strconv"
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
// @Tags class_code
// @Accept json
// @Produce json
// @Param code query string true "Code to check"
// @Success 200 {object} bool "secretExists" "シークレットが存在します"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Router /cc/checkSecretExists [get]
func (c *ClassCodeController) CheckSecretExists(ctx *gin.Context) {
	code := ctx.Query("code")

	secretExists, err := c.classCodeService.CheckSecretExists(code)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"secretExists": secretExists})
}

// VerifyClassCode godoc
// @Summary グループコードとシークレットを検証＆ユーザーに役割を割り当てる
// @Description グループコードと、該当する場合はそのシークレットを確認する。また、指定されたユーザーに役割を割り当てる。
// @Tags class_code
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
func (c *ClassCodeController) VerifyClassCode(ctx *gin.Context) {
	code := ctx.Query("code")
	secret := ctx.Query("secret")
	uid, err := parseUintQueryParam(ctx, "uid")
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classCode, err := c.classCodeService.VerifyClassCode(code, secret)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	err = c.classUserService.AssignRole(uint(uid), classCode.CID, "APPLICANT")
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.AssignError)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.ClassMemberRegistration})
}

func parseUintQueryParam(ctx *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(ctx.Query(param), 10, 64)
}
