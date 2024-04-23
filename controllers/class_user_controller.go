package controllers

import (
	"errors"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"gorm.io/gorm"
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

// GetUserClassUserInfo godoc
// @Summary ユーザーに関連するクラスユーザー情報を取得
// @Description 特定のユーザーIDに基づいて、クラスユーザー情報を取得します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Success 200 {object} dto.ClassMemberDTO "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "情報が見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/info [get]
// @Security Bearer
func (c *ClassUserController) GetUserClassUserInfo(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	cidStr := ctx.Param("cid")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	classUserInfo, err := c.classUserService.GetClassUserInfo(uint(uid), uint(cid))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.UserNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, classUserInfo)
}

// GetUserClasses godoc
// @Summary ユーザーが参加しているクラスのリストを取得
// @Description 特定のユーザーが参加している全てのクラスの情報を取得します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {array} models.Class "成功"
// @Router /cu/{uid}/classes [get]
// @Security Bearer
func (c *ClassUserController) GetUserClasses(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	classes, err := c.classUserService.GetUserClasses(uint(uid), page, limit)
	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}
	//if err != nil {
	//	if errors.Is(err, services.ErrNotFound) {
	//		respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
	//	} else {
	//		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
	//	}
	//	return
	//}

	if len(classes) == 0 {
		respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, classes)
}

// GetClassMembers godoc
// @Summary クラスメンバーの情報を取得します
// @Description 指定されたcidのクラスに所属しているメンバーの情報を取得します。
// @Tags Class User
// @Accept  json
// @Produce  json
// @Param cid path int true "クラスID"
// @Param roleID query int false "ロールID"
// @Success 200 {array} dto.ClassMemberDTO "成功時、クラスメンバーの情報を返します"
// @Failure 400 {object} map[string]interface{} "無効なクラスIDが指定された場合のエラーメッセージ"
// @Failure 500 {object} map[string]interface{} "サーバー内部エラー"
// @Router /cu/class/{cid}/{role}/members [get]
// @Security Bearer
func (c *ClassUserController) GetClassMembers(ctx *gin.Context) {
	cid, err := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	roleIDStr := ctx.DefaultQuery("roleID", "")
	var roleID int
	if roleIDStr != "" {
		roleID, err = strconv.Atoi(roleIDStr)
		if err != nil {
			respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
			return
		}
	}

	var members []dto.ClassMemberDTO
	if roleIDStr == "" {
		members, err = c.classUserService.GetClassMembers(uint(cid))
	} else {
		members, err = c.classUserService.GetClassMembers(uint(cid), roleID)
	}

	if err != nil {
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		return
	}

	if len(members) == 0 {
		respondWithSuccess(ctx, constants.StatusOK, []dto.ClassMemberDTO{})
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
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {array} dto.UserClassInfoDTO "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "クラスが見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/favorite-classes [get]
// @Security Bearer
func (c *ClassUserController) GetFavoriteClasses(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	favoriteClasses, err := c.classUserService.GetFavoriteClasses(uint(uid), page, limit)
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
// @Param uid path int true "ユーザーID"
// @Param roleID path int true "ロールID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {array} dto.UserClassInfoDTO "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "クラスが見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/classes/{roleID} [get]
// @Security Bearer
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

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	classes, err := c.classUserService.GetUserClassesByRole(uint(uid), roleID, page, limit)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.ClassNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
	}

	if len(classes) == 0 {
		respondWithSuccess(ctx, constants.StatusOK, []dto.UserClassInfoDTO{})
	}

	respondWithSuccess(ctx, constants.StatusOK, classes)
}

// ChangeUserRole godoc
// @Summary ユーザーのロールを変更します。
// @Description ユーザーのロールを変更します。
// @Tags Class User
// @ID change-user-role
// @Accept  json
// @Produce  json
// @Param uid path string true "ユーザーID"
// @Param cid path string true "クラスID"
// @Param roleID path string true "ロールID"
// @Success 200 {string} string "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/{roleID} [patch]
// @Security Bearer
func (c *ClassUserController) ChangeUserRole(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	cidStr := ctx.Param("cid")
	roleIDStr := ctx.Param("roleID")

	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil || !isValidRoleID(roleID) {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	err = c.classUserService.AssignRole(uint(uid), uint(cid), roleID)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, "ユーザーまたはクラスが見つかりません")
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}

func isValidRoleID(roleID int) bool {
	validRoleIDs := map[int]bool{
		1: true, // USER
		2: true, // ADMIN
		3: true, // ASSISTANT
		4: true, // APPLICANT
		5: true, // BLACKLIST
		6: true, // INVITE
	}

	_, isValid := validRoleIDs[roleID]
	return isValid
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
// @Security Bearer
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

// ToggleFavorite godoc
// @Summary お気に入りのクラスを切り替えます
// @Description ユーザーIDとクラスIDに基づいて、お気に入りのクラスを切り替えます。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Success 200 {string} string "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "ユーザーまたはクラスが見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/toggle-favorite [patch]
// @Security Bearer
func (c *ClassUserController) ToggleFavorite(ctx *gin.Context) {
	uid, uidErr := strconv.ParseUint(ctx.Param("uid"), 10, 32)
	cid, cidErr := strconv.ParseUint(ctx.Param("cid"), 10, 32)
	if uidErr != nil || cidErr != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	err := c.classUserService.ToggleFavorite(uint(uid), uint(cid))
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.UserNClassNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, constants.Success)
}

// RemoveUserFromClass godoc
// @Summary ユーザーをクラスから削除します。
// @Description 指定したユーザーIDとクラスIDに基づいて、ユーザーをクラスから削除します。
// @Tags Class User
// @Accept json
// @Produce json
// @Param uid path int true "ユーザーID"
// @Param cid path int true "クラスID"
// @Success 200 {string} string "成功"
// @Failure 400 {string} string "無効なリクエスト"
// @Failure 404 {string} string "ユーザーまたはクラスが見つかりません"
// @Failure 500 {string} string "サーバーエラーが発生しました"
// @Router /cu/{uid}/{cid}/remove [delete]
// @Security Bearer
func (c *ClassUserController) RemoveUserFromClass(ctx *gin.Context) {
	uidStr := ctx.Param("uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	cidStr := ctx.Param("cid")
	cid, err := strconv.ParseUint(cidStr, 10, 32)
	if err != nil {
		respondWithError(ctx, constants.StatusBadRequest, constants.InvalidRequest)
		return
	}

	err = c.classUserService.RemoveUserFromClass(uint(uid), uint(cid))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respondWithError(ctx, constants.StatusNotFound, constants.UserNotFound)
		} else {
			respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
		}
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, constants.DeleteSuccess)
}
