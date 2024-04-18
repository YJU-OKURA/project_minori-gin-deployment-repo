package controllers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

type GoogleAuthController struct {
	Service    services.GoogleAuthService
	JWTService services.JWTService
}

func NewGoogleAuthController(service services.GoogleAuthService, jwtService services.JWTService) *GoogleAuthController {
	return &GoogleAuthController{
		Service:    service,
		JWTService: jwtService,
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
		handleServiceError(c, fmt.Errorf("Invalid request"))
		return
	}

	authCode, ok := requestBody["authCode"]
	if !ok {
		handleServiceError(c, fmt.Errorf("Auth code is required"))
		return
	}

	userInfo, err := controller.Service.GetGoogleUserInfo(authCode)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	var userInput dto.UserInput
	if err := json.Unmarshal(userInfo, &userInput); err != nil {
		handleServiceError(c, err)
		return
	}

	user, err := controller.Service.UpdateOrCreateUser(userInput)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	accessToken, err := controller.JWTService.GenerateToken(user.ID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	refreshToken, err := controller.JWTService.GenerateRefreshToken(user.ID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": map[string]interface{}{
			"name":  user.Name,
			"image": user.Image,
			"id":    user.ID,
		},
	})
}

// RefreshAccessTokenHandler godoc
// @Summary アクセストークンの更新
// @Description 提供されたリフレッシュトークンを使用してアクセストークンを更新します
// @Tags GoogleAuth
// @Accept  json
// @Produce  json
// @Param   refresh_token     body    string  true  "リフレッシュトークン"
// @Success 200 {object} map[string]interface{} "アクセストークンと有効期限が返されます"
// @Failure 400 {object} map[string]interface{} "JSON形式が不正、またはリフレッシュトークンが提供されていない場合のエラー"
// @Failure 401 {object} map[string]interface{} "リフレッシュトークンが無効または期限切れの場合の認証エラー"
// @Failure 500 {object} map[string]interface{} "未処理のエラーによる内部サーバーエラー"
// @Router /auth/google/refresh-token [post]
func (controller *GoogleAuthController) RefreshAccessTokenHandler(c *gin.Context) {
	var requestBody map[string]string
	if err := c.BindJSON(&requestBody); err != nil {
		handleServiceError(c, fmt.Errorf("Invalid JSON format or structure: %v", err))
		return
	}

	refreshToken, ok := requestBody["refresh_token"]
	if !ok {
		handleServiceError(c, fmt.Errorf("Refresh token is required"))
		return
	}

	refreshToken = strings.TrimSpace(strings.TrimPrefix(refreshToken, "Bearer "))

	tokenDetails, err := controller.JWTService.RefreshAccessToken(refreshToken)
	if err != nil {
		handleServiceError(c, fmt.Errorf("Failed to refresh access token: %v", err))
		return
	}

	respondWithSuccess(c, constants.StatusOK, gin.H{
		"access_token": tokenDetails.Raw,
		"expires_in":   tokenDetails.Claims,
	})
}
