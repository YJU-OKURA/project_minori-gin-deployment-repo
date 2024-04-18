package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
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

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	GenerateRefreshToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	RefreshAccessToken(refreshToken string) (*jwt.Token, error)
}

type JWTServiceImpl struct {
	secretKey []byte
}

func NewJWTService() *JWTServiceImpl {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT secret is not set")
	}
	return &JWTServiceImpl{
		secretKey: []byte(secret),
	}
}

func (s *JWTServiceImpl) GenerateToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(3 * time.Hour).Unix(),
	})
	return token.SignedString(s.secretKey)
}

func (s *JWTServiceImpl) GenerateRefreshToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour * 7).Unix(),
		"type":    "refresh",
	})
	return token.SignedString(s.secretKey)
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
}

func (s *JWTServiceImpl) RefreshAccessToken(refreshToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] == "refresh" {
			newAccessToken, err := s.GenerateToken(uint(claims["user_id"].(float64)))
			if err != nil {
				return nil, err
			}
			return jwt.Parse(newAccessToken, func(token *jwt.Token) (interface{}, error) {
				return s.secretKey, nil
			})
		}
	}
	return nil, fmt.Errorf("invalid refresh token")
}

// GoogleAuthServiceはGoogle認証サービスのインターフェース
type GoogleAuthService interface {
	GenerateStateOauthCookie(w http.ResponseWriter) string
	GetGoogleUserInfo(code string) ([]byte, error)
	OauthConfig() *oauth2.Config
	UpdateOrCreateUser(userInput dto.UserInput) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
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

func (s *GoogleAuthServiceImpl) UpdateOrCreateUser(userInput dto.UserInput) (models.User, error) {
	return s.repo.UpdateOrCreateUser(userInput)
}

// OauthConfigはOAuth設定を返す
func (s *GoogleAuthServiceImpl) OauthConfig() *oauth2.Config {
	return s.oauthConfig
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
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %s\n", err.Error())
	}

	return body, nil
}
