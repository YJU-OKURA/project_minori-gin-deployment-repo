package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
// @Security ApiKeyAuth
// @CrossOrigin
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
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /cb [post]
// @Security Bearer
func (c *ClassBoardController) CreateClassBoard(ctx *gin.Context) {
	var createDTO dto.ClassBoardCreateDTO
	if err := ctx.ShouldBindWith(&createDTO, binding.FormMultipart); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.BadRequestMessage)
		return
	}

	cid, err := strconv.ParseUint(ctx.PostForm("cid"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid class ID")
		return
	}

	imageUrl, err := c.handleImageUpload(ctx, uint(cid))
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

	msg := fmt.Sprintf("data: %s\n\n", "New class board created")
	c.classBoardService.GetUpdateNotifier().Broadcast <- []byte(msg)

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// GetClassBoardByID godoc
// @Summary IDでグループ掲示板を取得
// @Description 指定されたIDのグループ掲示板の詳細を取得します。
// @Tags Class Board
// @CrossOrigin
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Success 200 {object} models.ClassBoard "グループ掲示板が取得されました"
// @Failure 400 {object} string "無効なリクエストです"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id} [get]
// @Security Bearer
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
// @CrossOrigin
// @Accept json
// @Produce json
// @Param cid query int true "Class ID"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {array} []models.ClassBoard "全てのグループ掲示板のリスト"
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb [get]
// @Security Bearer
func (c *ClassBoardController) GetAllClassBoards(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Query("cid"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid page number")
		return
	}

	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid page size")
		return
	}

	result, err := c.classBoardService.GetAllClassBoards(uint(cid), page, pageSize)
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
// @CrossOrigin
// @Accept json
// @Produce json
// @Param cid query int true "Class ID"
// @Success 200 {array} []models.ClassBoard "公告されたグループ掲示板のリスト"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/announced [get]
// @Security Bearer
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
// @CrossOrigin
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Param cid path int true "Class ID"
// @Param uid path int true "User ID"
// @Param class_board_update body dto.ClassBoardUpdateDTO true "クラス掲示板の更新"
// @Success 200 {object} models.ClassBoard "グループ掲示板が正常に更新されました"
// @Failure 400 {object} string "リクエストが不正です"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id}/{cid}/{uid} [patch]
// @Security Bearer
func (c *ClassBoardController) UpdateClassBoard(ctx *gin.Context) {
	log.Println("Received ID:", ctx.Param("id"))
	ID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("Error parsing ID:", err)
		respondWithError(ctx, constants.StatusBadRequest, "Invalid class board ID")
		return
	}

	var updateDTO dto.ClassBoardUpdateDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		log.Println("Error binding JSON:", err)
		respondWithError(ctx, constants.StatusBadRequest, "Invalid JSON data")
		return
	}

	imageUrl := updateDTO.Image
	if ctx.GetHeader("Content-Type") == "multipart/form-data" {
		var uploadErr error
		imageUrl, uploadErr = c.handleImageUpload(ctx, uint(ID))
		if uploadErr != nil {
			log.Println("Error handling image upload: ", uploadErr)
			handleServiceError(ctx, uploadErr)
			return
		}
	}

	result, err := c.classBoardService.UpdateClassBoard(uint(ID), updateDTO, imageUrl)
	if err != nil {
		log.Println("Error updating class board:", err)
		handleServiceError(ctx, err)
		return
	}

	msg := fmt.Sprintf("data: %s\n\n", "Class board updated")
	c.classBoardService.GetUpdateNotifier().Broadcast <- []byte(msg)

	respondWithSuccess(ctx, constants.StatusOK, result)
}

// DeleteClassBoard godoc
// @Summary グループ掲示板を削除
// @Description 指定されたIDのグループ掲示板を削除します。
// @Tags Class Board
// @CrossOrigin
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Class Board ID"
// @Param cid query int true "Class ID"
// @Param uid query int true "User ID"
// @Success 200 {object} string "クラス掲示板が正常に削除されました"
// @Failure 400 {string} string "無効なリクエストです"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {object} string "コードが見つかりません"
// @Failure 500 {object} string "サーバーエラーが発生しました"
// @Router /cb/{id} [delete]
// @Security Bearer
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
func (c *ClassBoardController) handleImageUpload(ctx *gin.Context, cid uint) (string, error) {
	// Check if there's any file part
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", nil // No file was uploaded, proceed without error
		}
		return "", err // Other errors are still considered as errors
	}

	imageUrl, err := c.uploader.UploadImage(fileHeader, cid, false)
	if err != nil {
		return "", err
	}

	return imageUrl, nil
}

// SubscribeClassBoardUpdates godoc
// @Summary クラス掲示板の更新を購読
// @Description クラス掲示板の更新を購読します。
// @Tags Class Board
// @CrossOrigin
// @Produce text/event-stream
// @Success 200 {string} string "Class board updates subscribed"
// @Failure 500 {string} string "Error setting up SSE connection."
// @Router /cb/subscribe [get]
// @Security Bearer
// @Notes Clients should reconnect automatically in case the connection closes.
func (c *ClassBoardController) SubscribeClassBoardUpdates(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	notifier := c.classBoardService.GetUpdateNotifier()

	// Adding defer to recover from any panic and avoid crashing the server
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in SubscribeClassBoardUpdates: %v", r)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		}
		notifier.Unregister <- ctx.Writer
	}()

	// Register the client for updates
	notifier.Register <- ctx.Writer

	// Keep the connection open until the client disconnects
	<-ctx.Request.Context().Done()
}

// SearchClassBoards godoc
// @Summary クラス掲示板を検索
// @Description タイトルに基づいてクラス掲示板を検索します。
// @Tags Class Board
// @Accept json
// @Produce json
// @Param cid query int true "Class ID" example="1"
// @Param title query string true "Title to search" example="Welcome to Class"
// @Success 200 {array} []models.ClassBoard "Search results"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "No class boards found"
// @Failure 500 {string} string "Server error"
// @Router /cb/search [get]
func (c *ClassBoardController) SearchClassBoards(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Query("cid"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, "Invalid class ID")
		return
	}

	title := ctx.Query("title")
	if title == "" {
		respondWithError(ctx, constants.StatusBadRequest, "Title parameter is required")
		return
	}

	result, err := c.classBoardService.SearchClassBoardsByTitle(title, uint(cid))
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	if len(result) == 0 {
		respondWithError(ctx, constants.StatusNotFound, "No class boards found")
	}

	respondWithSuccess(ctx, constants.StatusOK, result)
}
