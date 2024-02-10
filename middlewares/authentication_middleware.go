package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"github.com/swaggo/gin-swagger"
	//"github.com/swaggo/files"
	//docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/backend/docs"
)

// 사용자의 인증 상태를 확인하고 관리하는 미들웨어를 포함

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "Hello World")
}
