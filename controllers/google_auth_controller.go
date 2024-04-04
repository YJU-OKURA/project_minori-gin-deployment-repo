package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
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
	c.JSON(constants.StatusOK, gin.H{"url": url})
}

// GoogleAuthCallback godoc
// @Summary Googleログイン認証のコールバック処理
// @Description Googleログイン後にコールバックで受け取ったコードを使用してユーザー情報を取得し、ユーザー情報を基にトークンを生成します。
// @Tags GoogleAuth
// @ID google-auth-callback
// @Accept json
// @Produce json
// @Param code query string true "Googleから返された認証コード"
// @Success 200 {object} map[string]interface{} "認証成功時、アクセストークン、リフレッシュトークン、ユーザー情報を返す"
// @Failure 400 {string} string "ユーザー情報の取得に失敗"
// @Failure 500 {string} string "内部サーバーエラー"
// @Router /auth/google/callback [get]
func (controller *GoogleAuthController) GoogleAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(constants.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	userInfo, err := controller.Service.GetGoogleUserInfo(code)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userInput dto.UserInput
	err = json.Unmarshal(userInfo, &userInput)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := controller.Service.UpdateOrCreateUser(userInput)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.Service.GenerateToken(user.ID)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := controller.Service.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(constants.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: false,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
	})

	c.Redirect(constants.StatusFound, "http://localhost:3000/")
}
