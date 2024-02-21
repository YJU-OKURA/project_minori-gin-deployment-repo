package main

import (
	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
)

func main() {
	if shouldRunMigrations() {
		db := migration.InitDB() // DB initialization
		migration.Migrate(db)
		os.Exit(0)
	}

	r := setupRouter()
	if err := r.Run(getServerPort()); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

func shouldRunMigrations() bool {
	return os.Getenv("RUN_MIGRATIONS") == "true"
}

func setupRouter() *gin.Engine {
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
	return r
}

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}
	return ":" + port
}

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
