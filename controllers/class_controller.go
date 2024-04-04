package controllers

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"

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
// @Tags Classes
// @Accept  json
// @Produce  json
// @Param cid path int true "クラスID"
// @Success 200 {object} models.Class "成功時、クラスの情報を返します"
// @Failure 400 {object} map[string]interface{} "error: リクエストが不正です"
// @Failure 404 {object} map[string]interface{} "error: クラスが見つかりません"
// @Failure 500 {object} map[string]interface{} "error: サーバーエラーが発生しました"
// @Router /cl/{cid} [get]
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
// @Router /cl/create [post]
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
