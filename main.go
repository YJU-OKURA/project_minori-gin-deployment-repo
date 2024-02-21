package main

import (
	"net/http"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"

	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

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
	g.JSON(http.StatusOK, "helloworld")
}

func main() {
	//infra.Initialize()
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Google OAuth2 ラウター
	r.GET("/", controllers.GoogleForm)
	r.GET("/auth/google/login", controllers.GoogleLoginHandler)
	r.GET("/auth/google/callback", controllers.GoogleAuthCallback)

	r.Run(":8080")
}
