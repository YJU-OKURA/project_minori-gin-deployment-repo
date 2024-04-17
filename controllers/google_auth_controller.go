package controllers

import (
	"encoding/json"
	"fmt"

	"os"
	"strings"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type GoogleAuthController struct {
	Service services.GoogleAuthService
}

func NewGoogleAuthController(service services.GoogleAuthService) *GoogleAuthController {
	return &GoogleAuthController{
		Service: service,
	}
}

// GoogleLoginHandler godoc
// @Summary Googleのログインページへリダイレクトします。
// @Description ユーザーをGoogleのログインページへリダイレクトして認証を行います。
// @Tags GoogleAuth
// @ID google-login-handler
// @Produce html
// @Success 302 "Googleのログインページへのリダイレクト"
// @Router /auth/google/login [get]
func (controller *GoogleAuthController) GoogleLoginHandler(c *gin.Context) {
	oauthStateString := controller.Service.GenerateStateOauthCookie(c.Writer)

	url := controller.Service.OauthConfig().AuthCodeURL(oauthStateString)
	respondWithSuccess(c, constants.StatusOK, gin.H{"url": url})
}

// ProcessAuthCode godoc
// @Summary 認可コードを処理します。
// @Description ユーザーがGoogleログイン後に受け取った認可コードを使って、ユーザー情報を照会し、トークンを生成します。
// @Tags GoogleAuth
// @Accept json
// @Produce json
// @Param authCode body string true "Googleから受け取った認可コード"
// @Success 200 {object} map[string]interface{} "ユーザー情報及びトークン情報"
// @Router /auth/google/process [post]
func (controller *GoogleAuthController) ProcessAuthCode(c *gin.Context) {
	var requestBody map[string]string
	if err := c.BindJSON(&requestBody); err != nil {
		handleServiceError(c, fmt.Errorf(constants.InvalidRequest))
		return
	}

	authCode, ok := requestBody["authCode"]
	if !ok {
		handleServiceError(c, fmt.Errorf(constants.AuthCodeRequired))
		return
	}

	userInfo, err := controller.Service.GetGoogleUserInfo(authCode)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	var userInput dto.UserInput
	err = json.Unmarshal(userInfo, &userInput)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	user, err := controller.Service.UpdateOrCreateUser(userInput)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	accessToken, err := controller.Service.GenerateToken(user.ID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	refreshToken, err := controller.Service.GenerateRefreshToken(user.ID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"name":  user.Name,
			"image": user.Image,
			"id":    user.ID,
		},
	})
}

func (ac *GoogleAuthController) RefreshAccessTokenHandler(ctx *gin.Context) {
	var requestBody map[string]string
	if err := ctx.BindJSON(&requestBody); err != nil {
		handleServiceError(ctx, fmt.Errorf(constants.InvalidRequest))
		return
	}

	refreshToken, ok := requestBody["refresh_token"]
	if !ok {
		handleServiceError(ctx, fmt.Errorf(constants.RefreshTokenRequired))
		return
	}

	refreshToken = strings.TrimSpace(strings.TrimPrefix(refreshToken, "Bearer "))

	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		handleServiceError(ctx, err)
		return
	}

	newToken, err := ac.Service.RefreshAccessToken(refreshToken)
	if err != nil {
		handleServiceError(ctx, fmt.Errorf("invalid refresh token: %v", err))
		return
	}

	respondWithSuccess(ctx, constants.StatusOK, gin.H{
		"access_token": newToken.AccessToken,
		"expires_in":   newToken.Expiry.Unix(),
	})
}

// func (ac *GoogleAuthController) RefreshAccessTokenHandler(ctx *gin.Context) {
// 	var requestBody map[string]string
// 	if err := ctx.BindJSON(&requestBody); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format or structure."})
// 		ctx.Error(err)
// 		return
// 	}

// 	refreshToken, ok := requestBody["refresh_token"]
// 	if !ok {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required."})
// 		ctx.Error(fmt.Errorf("refresh token missing in request body"))
// 		return
// 	}

// 	refreshToken = strings.TrimSpace(strings.TrimPrefix(refreshToken, "Bearer "))

// 	_, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(os.Getenv("JWT_SECRET")), nil
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token.", "details": err.Error()})
// 		ctx.Error(err)
// 		return
// 	}

// 	newToken, err := ac.Service.RefreshAccessToken(refreshToken)
// 	if err != nil {
// 		detailedError := fmt.Sprintf("Failed to refresh access token: %v. Ensure your refresh token is valid and not expired.", err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh access token.", "details": detailedError})
// 		ctx.Error(err)
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"access_token": newToken.AccessToken,
// 		"expires_in":   newToken.Expiry.Unix(),
// 	})
// }
