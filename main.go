package main

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/migration"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	if err := loadEnvironment(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := migration.InitDB()
	performMigrations(db)

	r := setupRouter(initializeControllers(db))
	startServer(r)
}

// shouldRunMigrations マイグレーションだけ実行するかどうか
func shouldRunMigrations() bool {
	return os.Getenv("RUN_MIGRATIONS") == "true"
}

// loadEnvironment .envファイルを読み込む
func loadEnvironment() error {
	return godotenv.Load()
}

// performMigrations マイグレーションを実行する
func performMigrations(db *gorm.DB) {
	if shouldRunMigrations() {
		migration.Migrate(db)
		os.Exit(0)
	}
}

// initializeControllers コントローラーを初期化する
func initializeControllers(db *gorm.DB) (*controllers.GroupCodeController, *controllers.GroupBoardController) {
	groupCodeRepo := repositories.NewGroupCodeRepository(db)
	groupBoardRepo := repositories.NewGroupBoardRepository(db)
	uploader := utils.NewAwsUploader()

	groupCodeController := controllers.NewGroupCodeController(services.NewGroupCodeService(groupCodeRepo))
	groupBoardController := controllers.NewGroupBoardController(services.NewGroupBoardService(groupBoardRepo), uploader)

	return groupCodeController, groupBoardController
}

// setupRouter ルーターをセットアップする
func setupRouter(groupCodeController *controllers.GroupCodeController, groupBoardController *controllers.GroupBoardController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	setupGroupCodeRoutes(r, groupCodeController)
	setupGroupBoardRoutes(r, groupBoardController)

	return r
}

// setupGroupCodeRoutes GroupCodeのルートをセットアップする
func setupGroupCodeRoutes(r *gin.Engine, controller *controllers.GroupCodeController) {
	gc := r.Group("/api/v1/gc")
	{
		gc.GET("/checkSecretExists", controller.CheckSecretExists) // シークレットが存在するかどうかを確認する
		gc.GET("/verifyGroupCode", controller.VerifyGroupCode)     // グループコードを検証する
	}
}

// setupGroupBoardRoutes GroupBoardのルートをセットアップする
func setupGroupBoardRoutes(r *gin.Engine, controller *controllers.GroupBoardController) {
	gb := r.Group("/api/v1/gb")
	{
		gb.POST("/", controller.CreateGroupBoard)                // グループ掲示板を作成する
		gb.GET("/:id", controller.GetGroupBoardByID)             // グループ掲示板を取得する
		gb.GET("/", controller.GetAllGroupBoards)                // 全てのグループ掲示板を取得する
		gb.GET("/announced", controller.GetAnnouncedGroupBoards) // 公開されたグループ掲示板を取得する
		gb.PATCH("/:id", controller.UpdateGroupBoard)            // グループ掲示板を更新する
		gb.DELETE("/:id", controller.DeleteGroupBoard)           // グループ掲示板を削除する
	}
}

// startServer サーバーを起動する
func startServer(r *gin.Engine) {
	if err := r.Run(getServerPort()); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

// getServerPort サーバーポートを取得する
func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080" // デフォルトポート
	}
	return ":" + port
}
