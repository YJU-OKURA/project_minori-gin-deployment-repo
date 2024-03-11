package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuthServiceはGoogle認証サービスのインターフェース
type GoogleAuthService interface {
	GenerateStateOauthCookie(w http.ResponseWriter) string
	GenerateToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	GetGoogleUserInfo(code string) ([]byte, error)
	OauthConfig() *oauth2.Config
	UpdateOrCreateUser(userInput dto.UserInput) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}

// GoogleAuthServiceImplはGoogle認証サービスの実装
type GoogleAuthServiceImpl struct {
	oauthConfig *oauth2.Config
	UrlAPI      string
	repo        repositories.GoogleAuthRepository
}

func (s *GoogleAuthServiceImpl) GetUserByID(id uint) (models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *GoogleAuthServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("tkq!mwe!d#s"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *GoogleAuthServiceImpl) UpdateOrCreateUser(userInput dto.UserInput) (models.User, error) {
	return s.repo.UpdateOrCreateUser(userInput)
}

// OauthConfigはOAuth設定を返す
func (s *GoogleAuthServiceImpl) OauthConfig() *oauth2.Config {
	return s.oauthConfig
}

func (s *GoogleAuthServiceImpl) GenerateToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *GoogleAuthServiceImpl) GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["type"] = "refresh"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// NewGoogleAuthServiceはGoogle認証サービスの新しいインスタンスを作成
func NewGoogleAuthService(repo repositories.GoogleAuthRepository) GoogleAuthService {
	return &GoogleAuthServiceImpl{
		oauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		UrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
		repo:   repo,
	}
}

// GenerateStateOauthCookieはOAuthのstateパラメータを生成し、それをクッキーに設定
func (s *GoogleAuthServiceImpl) GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

// GetGoogleUserInfoはGoogleのユーザー情報を取得
func (s *GoogleAuthServiceImpl) GetGoogleUserInfo(code string) ([]byte, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(s.UrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}

func (s *GoogleAuthServiceImpl) RefreshAccessToken(refreshToken string) (*oauth2.Token, error) {
	tokenSource := s.oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Failed to refresh access token: %v", err)
	}
	return newToken, nil
}
