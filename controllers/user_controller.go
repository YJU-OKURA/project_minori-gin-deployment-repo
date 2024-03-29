package controllers

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"strconv"
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
