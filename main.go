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
func initializeControllers(db *gorm.DB) (*controllers.ClassCodeController, *controllers.ClassBoardController) {
	classCodeRepo := repositories.NewClassCodeRepository(db)
	classBoardRepo := repositories.NewClassBoardRepository(db)
	classUserRepo := repositories.NewClassUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	classCodeService := services.NewClassCodeService(classCodeRepo)
	classUserService := services.NewClassUserService(classUserRepo, roleRepo)

	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	uploader := utils.NewAwsUploader()
	classBoardController := controllers.NewClassBoardController(services.NewClassBoardService(classBoardRepo), uploader)

	return classCodeController, classBoardController
}

// setupRouter ルーターをセットアップする
func setupRouter(classCodeController *controllers.ClassCodeController, classBoardController *controllers.ClassBoardController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	setupClassCodeRoutes(r, classCodeController)
	setupClassBoardRoutes(r, classBoardController)

	return r
}

// setupClassCodeRoutes ClassCodeのルートをセットアップする
func setupClassCodeRoutes(r *gin.Engine, controller *controllers.ClassCodeController) {
	gc := r.Group("/api/v1/cc")
	{
		gc.GET("/checkSecretExists", controller.CheckSecretExists) // シークレットが存在するかどうかを確認する
		gc.GET("/verifyClassCode", controller.VerifyClassCode)     // グループコードを検証する
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
func setupClassBoardRoutes(r *gin.Engine, controller *controllers.ClassBoardController) {
	gb := r.Group("/api/v1/cb")
	{
		gb.POST("/", controller.CreateClassBoard)                // グループ掲示板を作成する
		gb.GET("/:id", controller.GetClassBoardByID)             // グループ掲示板を取得する
		gb.GET("/", controller.GetAllClassBoards)                // 全てのグループ掲示板を取得する
		gb.GET("/announced", controller.GetAnnouncedClassBoards) // 公開されたグループ掲示板を取得する
		gb.PATCH("/:id", controller.UpdateClassBoard)            // グループ掲示板を更新する
		gb.DELETE("/:id", controller.DeleteClassBoard)           // グループ掲示板を削除する
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
