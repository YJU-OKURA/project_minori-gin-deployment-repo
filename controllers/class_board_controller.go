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

// classBoardController インタフェース
type classBoardController interface {
	CreateClassBoard(ctx *gin.Context)
	GetClassBoardByID(ctx *gin.Context)
	GetAllClassBoards(ctx *gin.Context)
	GetAnnouncedClassBoards(ctx *gin.Context)
	UpdateClassBoard(ctx *gin.Context)
	DeleteClassBoard(ctx *gin.Context)
}

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
// @Summary Create a new class board
// @Description Create a new class board with the provided information, including image upload.
// @Tags class_board
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
// @Router /gb [post]
func (c *ClassBoardController) CreateClassBoard(ctx *gin.Context) {
	var createDTO dto.ClassBoardCreateDTO

	log.Printf("Form data: %v", ctx.Request.Form)
	if err := ctx.ShouldBindWith(&createDTO, binding.FormMultipart); err != nil {
		log.Printf("Error in ShouldBindWith: %v", err)
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	imageUrl, err := c.handleImageUpload(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	createDTO.ImageURL = imageUrl
	// Proceed with service call
	result, err := c.classBoardService.CreateClassBoard(createDTO)
	if err != nil {
		log.Printf("Error in CreateClassBoard service: %v", err)
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetClassBoardByID godoc
// @Summary IDでグループ掲示板を取得
// @Description 指定されたIDのグループ掲示板の詳細を取得します。
// @Tags class_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Success 200 {object} models.ClassBoard "グループ掲示板が取得されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [get]
func (c *ClassBoardController) GetClassBoardByID(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
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
// @Tags class_board
// @Accept json
// @Produce json
// @Param cid query int true "クラスID"
// @Success 200 {array} []models.ClassBoard "全てのグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb [get]
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
// @Tags class_board
// @Accept json
// @Produce json
// @Param cid query int true "クラスID"
// @Success 200 {array} []models.ClassBoard "公告されたグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/announced [get]
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
// @Tags class_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Param class_board_update body dto.ClassBoardUpdateDTO true "グループ掲示板の更新"
// @Success 200 {object} models.ClassBoard "グループ掲示板が正常に更新されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [patch]
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
// @Tags class_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Success 200 {object} string "グループ掲示板が正常に削除されました"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [delete]
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
	respondWithSuccess(ctx, constants.StatusOK, "Class board deleted successfully")
}

// respondWithError エラーレスポンスを返す
func (c *ClassBoardController) handleImageUpload(ctx *gin.Context) (string, error) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil // No file was uploaded
		}
		return "", err // Error handling file upload
	}

	imageUrl, err := c.uploader.UploadImage(fileHeader)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}
