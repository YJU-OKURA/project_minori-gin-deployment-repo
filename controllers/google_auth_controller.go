package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GoogleAuthController struct {
	Service services.GoogleAuthService
}

func NewGoogleAuthController(service services.GoogleAuthService) *GoogleAuthController {
	return &GoogleAuthController{
		Service: service,
	}
}

type UserInput struct {
	ID      string `json:"id"`
	Picture string `json:"picture"`
	Name    string `json:"name"`
}

// @Summary Google Login
// @Tags Login
// @Description GoogleLoginのためのURLを取得
// @ID google-login
// @Produce  json
// @Success 200 {string} string	"url"
// @Router /auth/google/login [get]
func (controller *GoogleAuthController) GoogleLoginHandler(c *gin.Context) {
	oauthStateString := controller.Service.GenerateStateOauthCookie(c.Writer)

	url := controller.Service.OauthConfig().AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func generateStateOauthCookie() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// @Summary Google Login
// @Tags Login
// @Description GoogleLoginのためのURLを取得
// @ID google-login
// @Produce  json
// @Success 200 {string} string	"url"
// @Router /auth/google/login [get]
func (controller *GoogleAuthController) GoogleAuthCallback(c *gin.Context) {
	// Googleからコードを取得
	code := c.Query("code")

	userInfo, err := controller.Service.GetGoogleUserInfo(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user info",
		})
		return
	}

	// Googleから取得したユーザー情報をUserInput構造体にパース
	var userInput UserInput
	err = json.Unmarshal(userInfo, &userInput)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse user info",
		})
		return
	}

	// データベースに接続
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	result := db.Where("p_id = ?", userInput.ID).First(&user)

	// ユーザーが存在しない場合、新規ユーザーを登録
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		user = models.User{
			PID:   userInput.ID,
			Name:  userInput.Name,
			Image: userInput.Picture,
		}
		result = db.Create(&user)
		if result.Error != nil {
			log.Printf("Failed to create user: %v", result.Error)
			return
		}
	} else if result.Error != nil {
		log.Printf("Database error: %v", result.Error)
		return
	}

	// ユーザーが存在する場合、名前と画像、pidを更新
	db.Model(&user).Updates(models.User{Name: userInput.Name, Image: userInput.Picture, PID: userInput.ID})
	// アクセストークンとリフレッシュトークンを生成
	accessToken, err := generateToken(user)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err.Error())
		return
	}

	refreshToken, err := generateRefreshToken(user)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err.Error())
		return
	}

	// トークンをログに出力
	log.Printf("Access Token: %s\n", accessToken)
	log.Printf("Refresh Token: %s\n", refreshToken)

	// ユーザー情報を返す
	c.JSON(http.StatusOK, user)
}

func generateToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	// claims["name"] = user.Name
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix() // expiresInを 10mに設定

	t, err := token.SignedString([]byte("tkq!mwe!d#s")) // secretを "tkq!mwe!d#s"に 設定
	if err != nil {
		return "", err
	}

	return t, nil
}

func generateRefreshToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() //

	t, err := token.SignedString([]byte("wnca%dlod!")) // secret을 "wnca%dlod!"に設定
	if err != nil {
		return "", err
	}

	return t, nil
}

func WelcomeHandler(c *gin.Context) {
	// JWTトークンからユーザー情報を取得
	accessToken := c.Query("access_token")
	refreshToken := c.Query("refresh_token")
	if accessToken == "" || refreshToken == "" {
		// トークンがない場合、エラーメッセージを返す
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Access token or refresh token is missing",
		})
		return
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("tkq!mwe!d#s"), nil
	})
	if err != nil || !token.Valid {
		// トークンが無効な場合、エラーメッセージを返す
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid access token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get user info from access token",
		})
		return
	}

	user := &models.User{
		Name: claims["name"].(string),
	}

	c.JSON(http.StatusOK, user)
}
