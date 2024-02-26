package controllers

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// groupBoardController インタフェース
type groupBoardController interface {
	CreateGroupBoard(ctx *gin.Context)
	GetGroupBoardByID(ctx *gin.Context)
	GetAllGroupBoards(ctx *gin.Context)
	GetAnnouncedGroupBoards(ctx *gin.Context)
	UpdateGroupBoard(ctx *gin.Context)
	DeleteGroupBoard(ctx *gin.Context)
}

// GroupBoardController インタフェースを実装
type GroupBoardController struct {
	groupBoardService services.GroupBoardService
	uploader          utils.Uploader
}

// NewGroupBoardController GroupBoardControllerを生成
func NewGroupBoardController(service services.GroupBoardService, uploader utils.Uploader) *GroupBoardController {
	return &GroupBoardController{
		groupBoardService: service,
		uploader:          uploader,
	}
}

// CreateGroupBoard godoc
// @Summary Create a new group board
// @Description Create a new group board with the provided information, including image upload.
// @Tags group_board
// @Accept multipart/form-data
// @Produce json
// @Param group_board_create body dto.GroupBoardCreateDTO true "Create group board"
// @Param image formData string false "Upload image file"
// @Success 200 {object} models.GroupBoard "Group board created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /gb [post]
func (c *GroupBoardController) CreateGroupBoard(ctx *gin.Context) {
	var createDTO dto.GroupBoardCreateDTO
	if err := ctx.ShouldBind(&createDTO); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	imageUrl, err := c.handleImageUpload(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	result, err := c.groupBoardService.CreateGroupBoard(createDTO, imageUrl)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetGroupBoardByID godoc
// @Summary IDでグループ掲示板を取得
// @Description 指定されたIDのグループ掲示板の詳細を取得します。
// @Tags group_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Success 200 {object} models.GroupBoard "グループ掲示板が取得されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [get]
func (c *GroupBoardController) GetGroupBoardByID(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}
	result, err := c.groupBoardService.GetGroupBoardByID(uint(ID))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetAllGroupBoards godoc
// @Summary 全てのグループ掲示板を取得
// @Description 全てのグループ掲示板のリストを取得します。
// @Tags group_board
// @Accept json
// @Produce json
// @Success 200 {array} []models.GroupBoard "全てのグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb [get]
func (c *GroupBoardController) GetAllGroupBoards(ctx *gin.Context) {
	result, err := c.groupBoardService.GetAllGroupBoards()
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetAnnouncedGroupBoards godoc
// @Summary 公告されたグループ掲示板を取得
// @Description 公告されたグループ掲示板のリストを取得します。
// @Tags group_board
// @Accept json
// @Produce json
// @Success 200 {array} []models.GroupBoard "公告されたグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/announced [get]
func (c *GroupBoardController) GetAnnouncedGroupBoards(ctx *gin.Context) {
	result, err := c.groupBoardService.GetAnnouncedGroupBoards()
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// UpdateGroupBoard godoc
// @Summary グループ掲示板を更新
// @Description 指定されたIDのグループ掲示板の詳細を更新します。
// @Tags group_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Param group_board_update body dto.GroupBoardUpdateDTO true "グループ掲示板の更新"
// @Success 200 {object} models.GroupBoard "グループ掲示板が正常に更新されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [patch]
func (c *GroupBoardController) UpdateGroupBoard(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	var updateDTO dto.GroupBoardUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	imageUrl, err := c.handleImageUpload(ctx)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	result, err := c.groupBoardService.UpdateGroupBoard(uint(ID), updateDTO, imageUrl)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, result)
}

// DeleteGroupBoard godoc
// @Summary グループ掲示板を削除
// @Description 指定されたIDのグループ掲示板を削除します。
// @Tags group_board
// @Accept json
// @Produce json
// @Param id path int true "グループ掲示板ID"
// @Success 200 {object} string "グループ掲示板が正常に削除されました"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /gb/{id} [delete]
func (c *GroupBoardController) DeleteGroupBoard(ctx *gin.Context) {
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	err = c.groupBoardService.DeleteGroupBoard(uint(ID))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, "Group board deleted successfully")
}

// respondWithError エラーレスポンスを返す
func (c *GroupBoardController) handleImageUpload(ctx *gin.Context) (string, error) {
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
