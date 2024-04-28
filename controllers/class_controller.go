package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ClassController struct {
	classService services.ClassService
	uploader     utils.Uploader
}

func NewCreateClassController(classService services.ClassService, uploader utils.Uploader) *ClassController {
	return &ClassController{
		classService: classService,
		uploader:     uploader,
	}
}

// GetClass godoc
// @Summary クラスの情報を取得します
// @Description 指定されたIDを持つクラスの情報を取得します。
// @Tags Class
// @Accept  json
// @Produce  json
// @Param cid path int true "クラスID"
// @Success 200 {object} models.Class "成功時、クラスの情報を返します"
// @Failure 400 {object} map[string]interface{} "error: リクエストが不正です"
// @Failure 404 {object} map[string]interface{} "error: クラスが見つかりません"
// @Failure 500 {object} map[string]interface{} "error: サーバーエラーが発生しました"
// @Router /cl/{cid} [get]
// @Security Bearer
func (cc *ClassController) GetClass(ctx *gin.Context) {
	classID, err := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	class, err := cc.classService.GetClass(uint(classID))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		return
	}
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"class": class})
}

// CreateClass godoc
// @Summary 新しいクラスを作成します
// @Description 名前、定員、説明、画像URL、作成者のUIDを持つ新しいクラスを作成します。画像はオプショナルです。
// @Tags Class
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "クラスの名前"
// @Param limitation formData int false "クラスの定員数"
// @Param description formData string false "クラスの説明"
// @Param uid formData int true "クラスを作成するユーザーのUID"
// @Param secret formData string false "クラス加入暗証番号"
// @Param image formData file false "クラスの画像"
// @Success 201 {object} map[string]interface{} "message: クラスが正常に作成されました"
// @Failure 400 {object} map[string]interface{} "error: 不正なリクエストのエラーメッセージ"
// @Failure 500 {object} map[string]interface{} "error: サーバー内部エラー"
// @Router /cl/create [post]
// @Security Bearer
func (cc *ClassController) CreateClass(ctx *gin.Context) {
	var createDTO dto.CreateClassRequest
	if err := ctx.ShouldBindWith(&createDTO, binding.FormMultipart); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	classID, err := cc.classService.CreateClass(createDTO)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	fileHeader, _ := ctx.FormFile("image")
	var imageUrl string
	if fileHeader != nil {
		imageUrl, err = cc.uploader.UploadImage(fileHeader, classID, false)
		if err != nil {
			handleServiceError(ctx, err)
			return
		}
	}

	if imageUrl != "" {
		err = cc.classService.UpdateClassImage(classID, imageUrl)
		if err != nil {
			handleServiceError(ctx, err)
			return
		}
	}

	respondWithSuccess(ctx, constants.StatusCreated, gin.H{"message": "Class created successfully", "classID": classID, "imageUrl": imageUrl})
}

func (c *ClassController) handleImageUpload(ctx *gin.Context) (string, error) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil
		}
		return "", err
	}

	tempClassID := uint(0)
	imageUrl, err := c.uploader.UploadImage(fileHeader, tempClassID, false)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

// UpdateClass godoc
// @Summary クラス情報を更新します
// @Description 指定されたIDを持つクラスの情報を更新します。
// @Tags Class
// @Accept multipart/form-data
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Param name formData string false "クラス名"
// @Param limitation formData int false "参加制限人数"
// @Param description formData string false "クラス説明"
// @Param image formData file false "クラス画像"
// @Success 200 {object} map[string]interface{} "message: クラスが正常に更新されました"
// @Failure 400 {object} map[string]interface{} "error: 不正なリクエストのエラーメッセージ"
// @Failure 401 {object} map[string]interface{} "error: 認証エラー"
// @Failure 500 {object} map[string]interface{} "error: サーバー内部エラー"
// @Router /cl/{uid}/{cid} [patch]
// @Security Bearer
func (cc *ClassController) UpdateClass(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	classID, _ := strconv.ParseUint(ctx.Param("cid"), 10, 32)

	var updateDTO dto.UpdateClassRequest
	if err := ctx.ShouldBindWith(&updateDTO, binding.Form); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	if fileHeader, _ := ctx.FormFile("image"); fileHeader != nil {
		imageUrl, fileErr := cc.uploader.UploadImage(fileHeader, uint(classID), false)
		if fileErr != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Image upload failed: "+fileErr.Error())
			return
		}

		// Call a separate method to update the image URL
		if err := cc.classService.UpdateClassImage(uint(classID), imageUrl); err != nil {
			respondWithError(ctx, constants.StatusInternalServerError, "Failed to update class image: "+err.Error())
			return
		}
	}

	if err := cc.classService.UpdateClass(uint(classID), uint(userID), updateDTO); err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, "Class update failed: "+err.Error())
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}

// DeleteClass godoc
// @Summary クラスを削除します
// @Description 指定されたIDを持つクラスを削除します。
// @Tags Class
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Success 200 {object} map[string]interface{} "message: クラスが正常に削除されました"
// @Failure 401 {object} map[string]interface{} "error: 認証エラー"
// @Failure 500 {object} map[string]interface{} "error: サーバー内部エラー"
// @Router /cl/{uid}/{cid} [delete]
// @Security Bearer
func (cc *ClassController) DeleteClass(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	classID, _ := strconv.ParseUint(ctx.Param("cid"), 10, 32)

	err := cc.classService.DeleteClass(uint(classID), uint(userID))
	if err != nil {
		respondWithError(ctx, constants.StatusUnauthorized, fmt.Sprintf("Error: %v", err))
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.DeleteSuccess})
}
