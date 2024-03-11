package controllers

import (
	"encoding/json"
	"log"
	"net/http"

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

func (controller *GoogleAuthController) GoogleLoginHandler(c *gin.Context) {
	oauthStateString := controller.Service.GenerateStateOauthCookie(c.Writer)

	url := controller.Service.OauthConfig().AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (controller *GoogleAuthController) GoogleAuthCallback(c *gin.Context) {
	code := c.Query("code")

	userInfo, err := controller.Service.GetGoogleUserInfo(code)
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		handleServiceError(c, err)
	}

	var userInput dto.UserInput
	err = json.Unmarshal(userInfo, &userInput)
	if err != nil {
		log.Printf("Failed to parse user info: %v", err)
		handleServiceError(c, err)
		return
	}

	user, err := controller.Service.UpdateOrCreateUser(userInput)
	if err != nil {
		log.Printf("Failed to create or update user: %v", err)
		handleServiceError(c, err)
		return
	}

	token, err := controller.Service.GenerateToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		handleServiceError(c, err)
		return
	}

	refreshToken, err := controller.Service.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Failed to generate refreshToken: %v", err)
		handleServiceError(c, err)
		return
	}

	c.JSON(constants.StatusOK, gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
		"user":          user,
	})
}
