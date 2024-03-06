package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuthServiceはGoogle認証サービスのインターフェース
type GoogleAuthService interface {
	GenerateStateOauthCookie(w http.ResponseWriter) string
	GetGoogleUserInfo(code string) ([]byte, error)
	OauthConfig() *oauth2.Config
}

// GoogleAuthServiceImplはGoogle認証サービスの実装
type GoogleAuthServiceImpl struct {
	oauthConfig *oauth2.Config
	UrlAPI      string
}

// OauthConfigはOAuth設定を返す
func (s *GoogleAuthServiceImpl) OauthConfig() *oauth2.Config {
	return s.oauthConfig
}

// NewGoogleAuthServiceはGoogle認証サービスの新しいインスタンスを作成
func NewGoogleAuthService() GoogleAuthService {
	return &GoogleAuthServiceImpl{
		oauthConfig: &oauth2.Config{
			RedirectURL:  "http://localhost:8080/auth/google/callback",
			ClientID:     "500381013046-96fkfveskcksbd77a833qobg48887jcd.apps.googleusercontent.com",
			ClientSecret: "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",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
			Endpoint:     google.Endpoint,
		},
		UrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
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
