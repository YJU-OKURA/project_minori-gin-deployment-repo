package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
)

func GoogleForm(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		"<html>"+
			"\n<head>\n    "+
			"<title>Go Oauth2.0 Test</title>\n"+
			"</head>\n"+
			"<body>\n<p>"+
			"<a href='./auth/google/login'>Google Login</a>"+
			"</p>\n"+
			"</body>\n"+
			"</html>"))
}

func GoogleLoginHandler(c *gin.Context) {
	state := services.GenerateStateOauthCookie(c.Writer)
	url := services.GoogleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := services.GetGoogleUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	fmt.Fprint(c.Writer, string(data))
}
