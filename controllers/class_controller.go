package controllers

import (
	"errors"
	"net/http"

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

// CreateClass godoc
// @Summary 新しいクラスを作成します
// @Description 名前、定員、説明、画像URLを持つ新しいクラスを作成します。画像はオプショナルです。
// @Tags Classes
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "クラスの名前"
// @Param limitation formData int false "クラスの定員数"
// @Param description formData string false "クラスの説明"
// @Param image formData file false "クラスの画像"
// @Success 201 {object} map[string]interface{} "message: クラスが正常に作成されました"
// @Failure 400 {object} map[string]interface{} "error: 不正なリクエストのエラーメッセージ"
// @Failure 500 {object} map[string]interface{} "error: サーバー内部エラー"
// @Router /cs/create [post]
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
