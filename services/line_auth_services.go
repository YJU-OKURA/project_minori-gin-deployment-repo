package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/dto"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"golang.org/x/oauth2"
)

type LINEAuthService interface {
	GenerateStateOauthCookie(w http.ResponseWriter) string
	GetLINEUserInfo(code string) ([]byte, error)
	OauthConfig() *oauth2.Config
	UpdateOrCreateUser(userInput dto.LineUserInput) (models.User, error)
	GetUserByID(userID uint) (models.User, error)
}

type LINEAuthServiceImpl struct {
	oauthConfig *oauth2.Config
	UrlAPI      string
	repo        repositories.LINEAuthRepository
}

func NewLINEAuthService(repo repositories.LINEAuthRepository) LINEAuthService {
	return &LINEAuthServiceImpl{
		oauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("LINE_REDIRECT_URL"),
			ClientID:     os.Getenv("LINE_CLIENT_ID"),
			ClientSecret: os.Getenv("LINE_CLIENT_SECRET"),
			Scopes:       []string{"profile", "openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
				TokenURL: "https://api.line.me/oauth2/v2.1/token",
			},
		},
		UrlAPI: "https://api.line.me/v2/profile",
		repo:   repo,
	}
}

func (s *LINEAuthServiceImpl) OauthConfig() *oauth2.Config {
	return s.oauthConfig
}

func (s *LINEAuthServiceImpl) GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func (s *LINEAuthServiceImpl) GetLINEUserInfo(code string) ([]byte, error) {
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange: %s\n", err.Error())
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get(s.UrlAPI)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo: %s\n", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %s\n", err.Error())
	}

	log.Println(string(body))

	return body, nil
}

func (s *LINEAuthServiceImpl) UpdateOrCreateUser(userInput dto.LineUserInput) (models.User, error) {
	return s.repo.UpdateOrCreateUser(userInput)
}

func (s *LINEAuthServiceImpl) GetUserByID(userID uint) (models.User, error) {
	return s.repo.GetUserByID(userID)
}
