package controllers

import (
	"errors"
	"strconv"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

// ClassUserController インタフェースを実装
type ClassUserController struct {
	classUserService services.ClassUserService
}

// NewClassUserController ClassScheduleControllerを生成
func NewClassUserController(service services.ClassUserService) *ClassUserController {
	return &ClassUserController{
		classUserService: service,
	}
}

// UpdateUserNameRequest ユーザー名更新リクエスト
type UpdateUserNameRequest struct {
	NewName string `json:"new_name"`
}

// ChangeUserRole godoc
// @Summary ユーザーのロールを変更します。
// @Description ユーザーのロールを変更します。
// @Tags Class User
// @ID change-user-role
// @Accept  json
// @Produce  json
// @Param uid path int true "User ID"
// @Param cid path int true "Class ID"
// @Param role path string true "Role Name"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/{role} [patch]
func (c *ClassUserController) ChangeUserRole(ctx *gin.Context) {
	uid, _ := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	cid, _ := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	roleName := ctx.Param("role")

	err := c.classUserService.AssignRole(uint(uid), uint(cid), roleName)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}
	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}

// UpdateUserName godoc
// @Summary ユーザーの名前を更新します。
// @Description 特定のユーザーIDとグループIDに対してユーザーの名前を更新します。
// @Tags Class User
// @ID update-user-name
// @Accept json
// @Produce json
// @Param uid path int true "User ID"
// @Param cid path int true "Class ID"
// @Param body body UpdateUserNameRequest true "新しいニックネーム"
// @Success 200 {string} string "成功"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/rename [put]
func (c *ClassUserController) UpdateUserName(ctx *gin.Context) {
	uid, err := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	cid, err := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	var requestBody UpdateUserNameRequest
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	err = c.classUserService.UpdateUserName(uint(uid), uint(cid), requestBody.NewName)
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{"message": constants.Success})
}

// GetUserClasses godoc
// @Summary ユーザーが参加しているクラスのリストを取得
// @Description 特定のユーザーが参加している全てのクラスの情報を取得します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Success 200 {array} models.Class "成功"
// @Router /cu/{uid}/classes [get]
func (c *ClassUserController) GetUserClasses(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classes, err := c.classUserService.GetUserClasses(uint(uid))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	if len(classes) == 0 {
		respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, classes)
}

// GetClassMembers godoc
// @Summary クラスメンバーの情報を取得します
// @Description 指定されたcidのクラスに所属しているメンバーの情報を取得します。
// @Tags classes
// @Accept  json
// @Produce  json
// @Param cid path int true "クラスID"
// @Success 200 {array} dto.ClassMemberDTO "成功時、クラスメンバーの情報を返します"
// @Failure 400 {object} map[string]interface{} "無効なクラスIDが指定された場合のエラーメッセージ"
// @Failure 500 {object} map[string]interface{} "サーバー内部エラー"
// @Router /cm/{cid}/members [get]
func (c *ClassUserController) GetClassMembers(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	members, err := c.classUserService.GetClassMembers(uint(cid))
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, members)
}

// GetFavoriteClasses godoc
// @Summary お気に入りのクラス情報を取得
// @Description ユーザーIDに基づいて、お気에入りに設定されたクラスの情報を取得します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Success 200 {array} dto.UserClassInfoDTO "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "クラスが見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/favorite-classes [get]
func (c *ClassUserController) GetFavoriteClasses(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	favoriteClasses, err := c.classUserService.GetFavoriteClasses(uint(uid))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	if len(favoriteClasses) == 0 {
		respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, favoriteClasses)
}

// GetUserClassesByRole godoc
// @Summary ユーザーとロールに関連するクラス情報を取得
// @Description ユーザーIDとロールIDに基づいて、ユーザーが所属しているクラスの情報を取得します。ロールIDが2の場合は自分が作ったクラスリスト、ロールIDが4の場合はユーザーから申し込んだクラスリスト、ロールIDが6の場合はクラスの管理者から招待されたクラスリストを取得します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "User ID"
// @Param roleID path int true "Role ID"
// @Success 200 {array} dto.UserClassInfoDTO "Success"
// @Router /cu/{uid}/classes/{roleID} [get]
func (c *ClassUserController) GetUserClassesByRole(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	roleIDStr := ctx.Param("roleID")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classes, err := c.classUserService.GetUserClassesByRole(uint(uid), roleID)
	if err != nil {
		// handle error
		return
	}

	if len(classes) == 0 {
		// handle no classes found
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, classes)
}
