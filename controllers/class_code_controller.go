package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ClassCodeController struct {
	Service          services.ClassCodeService
	ClassUserService services.ClassUserService
}

func NewClassCodeController(classCodeService services.ClassCodeService, classUserService services.ClassUserService) *ClassCodeController {
	return &ClassCodeController{
		Service:          classCodeService,
		ClassUserService: classUserService,
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
func (controller *ClassCodeController) CheckSecretExists(c *gin.Context) {
	code := c.Query("code")

	secretExists, err := controller.Service.CheckSecretExists(code)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{"secretExists": secretExists})
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
func (controller *ClassCodeController) VerifyClassCode(c *gin.Context) {
	code := c.Query("code")
	secret := c.Query("secret")
	uid, err := strconv.ParseUint(c.Query("uid"), 10, 64)
	if err != nil {
		respondWithError(c, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classCode, err := controller.Service.VerifyClassCode(code, secret)
	if err != nil {
		respondWithError(c, constants.StatusUnauthorized, constants.SecretMismatch)
		return
	}

	if classCode == nil {
		respondWithError(c, constants.StatusNotFound, constants.CodeNotFound)
		return
	}

	err = controller.ClassUserService.AssignRole(uint(uid), classCode.CID, "APPLICANT")
	if err != nil {
		respondWithError(c, constants.StatusInternalServerError, constants.AssignError)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{"message": constants.ClassMemberRegistration})
}
