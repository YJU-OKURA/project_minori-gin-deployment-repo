package controllers

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"strconv"
)

// ClassBoardController インタフェースを実装
type ClassBoardController struct {
	classBoardService services.ClassBoardService
	uploader          utils.Uploader
}

// NewClassBoardController ClassBoardControllerを生成
func NewClassBoardController(service services.ClassBoardService, uploader utils.Uploader) *ClassBoardController {
	return &ClassBoardController{
		classBoardService: service,
		uploader:          uploader,
	}
}

// CreateClassBoard godoc
// @Summary クラス掲示板を作成
// @Description クラス掲示板を作成します。
// @Tags Class Board
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Class board title"
// @Param content formData string true "Class board content"
// @Param cid formData int true "Class ID"
// @Param uid formData int true "User ID"
// @Param is_announced formData boolean false "Is announced"
// @Param image formData file false "Upload image file"
// @Success 200 {object} models.ClassBoard "Class board created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /cb [post]
func (c *ClassBoardController) CreateClassBoard(ctx *gin.Context) {
	var createDTO dto.ClassBoardCreateDTO
	if err := ctx.ShouldBindWith(&createDTO, binding.FormMultipart); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	imageUrl, err := c.handleImageUpload(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	createDTO.ImageURL = imageUrl

	result, err := c.classBoardService.CreateClassBoard(createDTO)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetClassBoardByID godoc
// @Summary IDでグループ掲示板を取得
// @Description 指定されたIDのグループ掲示板の詳細を取得します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Success 200 {object} models.ClassBoard "グループ掲示板が取得されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id} [get]
func (c *ClassBoardController) GetClassBoardByID(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	result, err := c.classBoardService.GetClassBoardByID(uint(ID))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetAllClassBoards godoc
// @Summary 全てのグループ掲示板を取得
// @Description cidに基づいて、グループの全ての掲示板を取得します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param cid query int true "Class ID"
// @Success 200 {array} []models.ClassBoard "全てのグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb [get]
func (c *ClassBoardController) GetAllClassBoards(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Query("cid"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	result, err := c.classBoardService.GetAllClassBoards(uint(cid))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetAnnouncedClassBoards godoc
// @Summary 公告されたグループ掲示板を取得
// @Description cidに基づいて、公告されたグループの掲示板を取得します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param cid query int true "Class ID"
// @Success 200 {array} []models.ClassBoard "公告されたグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/announced [get]
func (c *ClassBoardController) GetAnnouncedClassBoards(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Query("cid"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	result, err := c.classBoardService.GetAnnouncedClassBoards(uint(cid))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// UpdateClassBoard godoc
// @Summary グループ掲示板を更新
// @Description 指定されたIDのグループ掲示板の詳細を更新します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Param cid query int true "Class ID"
// @Param uid query int true "User ID"
// @Param class_board_update body dto.ClassBoardUpdateDTO true "クラス掲示板の更新"
// @Success 200 {object} models.ClassBoard "グループ掲示板が正常に更新されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id} [patch]
func (c *ClassBoardController) UpdateClassBoard(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	var updateDTO dto.ClassBoardUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	imageUrl, err := c.handleImageUpload(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	result, err := c.classBoardService.UpdateClassBoard(uint(ID), updateDTO, imageUrl)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// DeleteClassBoard godoc
// @Summary グループ掲示板を削除
// @Description 指定されたIDのグループ掲示板を削除します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Param cid query int true "Class ID"
// @Param uid query int true "User ID"
// @Success 200 {object} string "クラス掲示板が正常に削除されました"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id} [delete]
func (c *ClassBoardController) DeleteClassBoard(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	err = c.classBoardService.DeleteClassBoard(uint(ID))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, constants.DeleteSuccess)
}

// respondWithError エラーレスポンスを返す
func (c *ClassBoardController) handleImageUpload(ctx *gin.Context) (string, error) {
	cid, err := strconv.ParseUint(ctx.PostForm("cid"), 10, 64)
	if err != nil {
		log.Printf("Error parsing class ID: %v", err)
		return "", err
	}

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil // No file was uploaded
		}
		return "", err
	}

	imageUrl, err := c.uploader.UploadImage(fileHeader, uint(cid), false)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}
