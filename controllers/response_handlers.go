package controllers

import (
	"errors"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

// handleServiceError サービスによって返されたエラーを処理する
func handleServiceError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrNotFound):
		respondWithError(ctx, constants.StatusNotFound, constants.CodeNotFound)
	case errors.Is(err, services.ErrUnauthorized):
		respondWithError(ctx, constants.StatusUnauthorized, constants.Unauthorized)
	case errors.Is(err, services.ErrDatabase):
		respondWithError(ctx, constants.StatusInternalServerError, constants.DatabaseError)
	default:
		respondWithError(ctx, constants.StatusInternalServerError, constants.InternalServerError)
	}
}

// respondWithError エラーメッセージを返す
func respondWithError(ctx *gin.Context, statusCode int, errMsg string) {
	ctx.JSON(statusCode, gin.H{"error": errMsg})
}

// respondWithSuccess 成功時のレスポンスを返す
func respondWithSuccess(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, gin.H{"data": data})
}
