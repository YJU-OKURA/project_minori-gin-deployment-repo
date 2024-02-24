package controllers

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Param group_code_check_request body dto.GroupCodeCheckRequest true "Check Secret"
// @Success 200 {object} bool "secretExists" "シークレットが存在します"
// @Failure 400 {object} string "error" "無効なリクエストです"
// @Failure 404 {object} string "error" "コードが見つかりません"
// @Router /gc/checkSecretExists [post]
func (controller *GroupCodeController) CheckSecretExists(c *gin.Context) {
	var req dto.GroupCodeCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constants.StatusBadRequest, gin.H{"error": constants.InvalidRequest})
		return
	}

	secretExists, err := controller.Service.CheckSecretExists(req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(constants.StatusNotFound, gin.H{"error": constants.CodeNotFound})
		} else {
			c.JSON(constants.StatusInternalServerError, gin.H{"error": constants.InternalServerError})
		}
		return
	}

	c.JSON(constants.StatusOK, gin.H{"secretExists": secretExists})
}

// VerifyGroupCode godoc
// @Summary グループコードとシークレットを検証
// @Description グループコードと、該当する場合はそのシークレットを確認する。
// @Tags group_code
// @Accept json
// @Produce json
// @Param group_code_request body dto.GroupCodeRequest true "グループコードの確認"
// @Success 200 {object} string "message" "グループコードが検証されました"
// @Failure 400 {object} string "error" "無効なリクエストです"
// @Failure 401 {object} string "error" "シークレットが一致しません"
// @Failure 404 {object} string "error" "コードが見つかりません"
// @Router /gc/verifyGroupCode [post]
func (controller *GroupCodeController) VerifyGroupCode(c *gin.Context) {
	var req dto.GroupCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(constants.StatusBadRequest, gin.H{"error": constants.InvalidRequest})
		return
	}

	verified, err := controller.Service.VerifyGroupCode(req.Code, req.Secret)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(constants.StatusNotFound, gin.H{"error": constants.CodeNotFound})
		} else {
			c.JSON(constants.StatusInternalServerError, gin.H{"error": constants.InternalServerError})
		}
		return
	}

	if !verified {
		c.JSON(constants.StatusUnauthorized, gin.H{"error": constants.SecretMismatch})
		return
	}

	c.JSON(constants.StatusOK, gin.H{"message": constants.GroupCodeVerified})
}
