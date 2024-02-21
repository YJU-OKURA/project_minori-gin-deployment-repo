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

var GoogleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     "500381013046-96fkfveskcksbd77a833qobg48887jcd.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

var OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func GetGoogleUserInfo(code string) ([]byte, error) { // 17

	token, err := GoogleOauthConfig.Exchange(context.Background(), code) // 18
	if err != nil {                                                      // 19
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(OauthGoogleUrlAPI + token.AccessToken) // 20
	if err != nil {                                              // 21
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body) // 23
}
