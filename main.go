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
func initializeControllers(db *gorm.DB) (*controllers.ClassCodeController, *controllers.ClassBoardController, *controllers.ClassScheduleController, *controllers.AttendanceController) {
	classBoardRepo := repositories.NewClassBoardRepository(db)
	classCodeRepo := repositories.NewClassCodeRepository(db)
	classScheduleRepo := repositories.NewClassScheduleRepository(db)
	classUserRepo := repositories.NewClassUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)

	classCodeService := services.NewClassCodeService(classCodeRepo)
	classUserService := services.NewClassUserService(classUserRepo, roleRepo)
	classScheduleService := services.NewClassScheduleService(classScheduleRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)

	uploader := utils.NewAwsUploader()
	classBoardController := controllers.NewClassBoardController(services.NewClassBoardService(classBoardRepo), uploader)
	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	classScheduleController := controllers.NewClassScheduleController(classScheduleService)
	attendanceController := controllers.NewAttendanceController(attendanceService)

	return classCodeController, classBoardController, classScheduleController, attendanceController
}

// setupRouter ルーターをセットアップする
func setupRouter(classCodeController *controllers.ClassCodeController, classBoardController *controllers.ClassBoardController, classScheduleController *controllers.ClassScheduleController, controller *controllers.AttendanceController) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/gin"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	setupClassBoardRoutes(r, classBoardController)
	setupClassCodeRoutes(r, classCodeController)
	setupClassScheduleRoutes(r, classScheduleController)
	setupAttendanceRoutes(r, controller)

	return r
}

// setupClassCodeRoutes ClassCodeのルートをセットアップする
func setupClassCodeRoutes(r *gin.Engine, controller *controllers.ClassCodeController) {
	cc := r.Group("/api/gin/cc")
	{
		cc.GET("/checkSecretExists", controller.CheckSecretExists) // シークレットが存在するかどうかを確認する
		cc.GET("/verifyClassCode", controller.VerifyClassCode)     // グループコードを検証する
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
func setupClassBoardRoutes(r *gin.Engine, controller *controllers.ClassBoardController) {
	cb := r.Group("/api/gin/cb")
	{
		cb.GET("/", controller.GetAllClassBoards)                // 全てのグループ掲示板を取得する
		cb.GET("/:id", controller.GetClassBoardByID)             // グループ掲示板を取得する
		cb.POST("/", controller.CreateClassBoard)                // グループ掲示板を作成する
		cb.GET("/announced", controller.GetAnnouncedClassBoards) // 公開されたグループ掲示板を取得する
		cb.PATCH("/:id", controller.UpdateClassBoard)            // グループ掲示板を更新する
		cb.DELETE("/:id", controller.DeleteClassBoard)           // グループ掲示板を削除する
	}
}

// setupClassScheduleRoutes ClassScheduleのルートをセットアップする
func setupClassScheduleRoutes(r *gin.Engine, controller *controllers.ClassScheduleController) {
	cs := r.Group("/api/gin/cs")
	{
		cs.GET("/", controller.GetAllClassSchedules)        // 全てのクラススケジュールを取得する
		cs.GET("/:id", controller.GetClassScheduleByID)     // クラススケジュールを取得する
		cs.POST("/", controller.CreateClassSchedule)        // 新しいクラススケジュールを作成する
		cs.PATCH("/:id", controller.UpdateClassSchedule)    // クラススケジュールを更新する
		cs.DELETE("/:id", controller.DeleteClassSchedule)   // クラススケジュールを削除する
		cs.GET("/live", controller.GetLiveClassSchedules)   // ライブ中のクラススケジュールを取得する
		cs.GET("/date", controller.GetClassSchedulesByDate) // 日付でクラススケジュールを取得する
	}
}

// setupAttendanceRoutes Attendanceのルートをセットアップする
func setupAttendanceRoutes(r *gin.Engine, controller *controllers.AttendanceController) {
	at := r.Group("/api/gin/at")
	{
		at.POST("/:cid/:uid/:csid", controller.CreateOrUpdateAttendance) // 全ての出席を取得する
		at.GET("/:cid", controller.GetAllAttendances)                    // グループの全ての出席を取得する
		at.GET("/attendance/:id", controller.GetAttendance)              // 出席を取得する
		at.DELETE("/attendance/:id", controller.DeleteAttendance)        // 出席を削除する
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
