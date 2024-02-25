package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type GroupCodeController struct {
	Service services.GroupCodeService
}

func NewGroupCodeController(service services.GroupCodeService) *GroupCodeController {
	return &GroupCodeController{Service: service}
}

// CheckSecretExists godoc
// @Summary グループコードにシークレットが存在するかチェック
// @Description 指定されたグループコードにシークレットがあるかどうかをチェックする。
// @Tags group_code
// @Accept json
// @Produce json
// @Param code query string true "Code to check"
// @Success 200 {object} bool "secretExists" "シークレットが存在します"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Router /gc/checkSecretExists [get]
func (controller *GroupCodeController) CheckSecretExists(c *gin.Context) {
	code := c.Query("code")

	secretExists, err := controller.Service.CheckSecretExists(code)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{"secretExists": secretExists})
}

// VerifyGroupCode godoc
// @Summary グループコードとシークレットを検証
// @Description グループコードと、該当する場合はそのシークレットを確認する。
// @Tags group_code
// @Accept json
// @Produce json
// @Param code query string true "Code to verify"
// @Param secret query string true "Secret for the code"
// @Success 200 {object} string "グループコードが検証されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 401 {object} string "シークレットが一致しません"
// @Failure 404 {object} string "コードが見つかりません"
// @Router /gc/verifyGroupCode [get]
func (controller *GroupCodeController) VerifyGroupCode(c *gin.Context) {
	code := c.Query("code")
	secret := c.Query("secret")

	verified, err := controller.Service.VerifyGroupCode(code, secret)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	if !verified {
		respondWithError(c, constants.StatusUnauthorized, constants.SecretMismatch)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{"message": constants.GroupCodeVerified})
}
