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

type LINEAuthController struct {
	Service    services.LINEAuthService
	JWTService services.JWTService
}

func NewLINEAuthController(service services.LINEAuthService, jwtService services.JWTService) *LINEAuthController {
	return &LINEAuthController{
		Service:    service,
		JWTService: jwtService,
	}
}

func (controller *LINEAuthController) LINELoginHandler(c *gin.Context) {
	oauthStateString := controller.Service.GenerateStateOauthCookie(c.Writer)
	url := controller.Service.OauthConfig().AuthCodeURL(oauthStateString)
	respondWithSuccess(c, constants.StatusOK, gin.H{"url": url})
}

func (controller *LINEAuthController) ProcessAuthCode(c *gin.Context) {
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

	userInfo, err := controller.Service.GetLINEUserInfo(authCode)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	var userInput dto.LineUserInput
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

func (controller *LINEAuthController) RefreshAccessTokenHandler(c *gin.Context) {
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
