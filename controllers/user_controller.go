package controllers

import (
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewCreateUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetApplyingClasses godoc
// @Summary 申し込んだクラスを取得
// @Description ユーザーが申し込んだクラスを取得します。
// @Tags User
// @Accept json
// @Produce json
// @Param userID path int true "ユーザーID"
// @Success 200 {array} models.ClassUser
// @Failure 400 {object} string "無効なユーザーID"
// @Failure 404 {object} string "申請中のクラスが見つかりません"
// @Failure 500 {object} string "内部サーバーエラー"
// @Router /u/{userID}/applying-classes [get]
// @Security Bearer
func (uc *UserController) GetApplyingClasses(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 64)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.ErrNoUserID)
		return
	}

	classes, err := uc.userService.GetApplyingClasses(uint(userID))
	if err != nil {
		if err.Error() == services.ErrUserNotFound {
			respondWithError(ctx, constants.StatusNotFound, constants.UserNotFound)
		} else {
			handleServiceError(ctx, err)
		}
		return
	}

	if len(classes) == 0 {
		respondWithError(ctx, constants.StatusNotFound, constants.ApplyingClassNotFound)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"applyingClasses": classes})
}

// SearchByName godoc
// @Summary 名前でユーザーを検索
// @Description 名前でユーザーを検索します。
// @Tags User
// @Accept json
// @Produce json
// @Param name query string true "ユーザー名"
// @Success 200 {array} models.User
// @Failure 400 {object} string "Nameパラメーターが必要です"
// @Failure 500 {object} string "サーバーエラー"
// @Router /u/search [get]
func (uc *UserController) SearchByName(ctx *gin.Context) {
	name := ctx.Query("name")
	if name == "" {
		respondWithError(ctx, constants.StatusBadRequest, "Name parameter is required")
		return
	}

	users, err := uc.userService.SearchUsersByName(name)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, err.Error())
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"users": users})
}

// RemoveUserFromService godoc
// @Summary ユーザー削除
// @Description ユーザーIDによってサービスからユーザーを削除します。
// @Tags User
// @Accept  json
// @Produce  json
// @Param   userID   path    int  true  "ユーザーID"
// @Success 200 {object} map[string]interface{} "message: ユーザーが正常に削除されました。"
// @Failure 400 {object} map[string]interface{} "error: 不正なリクエスト、無効なユーザーIDです。"
// @Failure 404 {object} map[string]interface{} "error: ユーザーが見つかりません。"
// @Failure 500 {object} map[string]interface{} "error: サーバー内部エラーです。"
// @Router /u/{userID}/delete [delete]
func (c *UserController) RemoveUserFromService(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.ErrNoUserID)
		return
	}

	err = c.userService.RemoveUserFromService(uint(userID))
	if err != nil {
		if err.Error() == services.ErrUserNotFound {
			respondWithError(ctx, constants.StatusNotFound, constants.UserNotFound)
		} else {
			handleServiceError(ctx, err)
		}
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"deletedUserID": userID})
}
