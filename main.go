package main

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/migration"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := migration.InitDB() // DB接続
	if shouldRunMigrations() {
		migration.Migrate(db)
		os.Exit(0)
	}

	groupCodeController := controllers.NewGroupCodeController(db)
	r := setupRouter(groupCodeController)
	if err := r.Run(getServerPort()); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

// shouldRunMigrations マイグレーションだけ実行するかどうか
func shouldRunMigrations() bool {
	return os.Getenv("RUN_MIGRATIONS") == "true"
}

func setupRouter(groupCodeController *controllers.GroupCodeController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		gc := v1.Group("/gc")
		{
			gc.POST("/checkSecretExists", groupCodeController.CheckSecretExists) // グループコードにシークレットが存在するかチェック
			gc.POST("/verifyGroupCode", groupCodeController.VerifyGroupCode)     // グループコードとシークレットを検証
		}

	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}
	return ":" + port
}
