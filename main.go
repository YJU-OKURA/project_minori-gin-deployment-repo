package main

import (
	"log"
	"os"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/migration"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func main() {
	configureGinMode()
	ensureEnvVariables()

	db := initializeDatabase()
	migrateDatabaseIfNeeded(db)

	router := setupRouter(db)
	startServer(router)
}

func configureGinMode() {
	ginMode := getEnvOrDefault("GIN_MODE", gin.ReleaseMode)
	gin.SetMode(ginMode)
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ensureEnvVariables() {
	if err := godotenv.Load(); err != nil {
		log.Println("環境変数ファイルが読み込めませんでした。")
	}

	requiredVars := []string{"MYSQL_HOST", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE", "MYSQL_PORT"}

	for _, varName := range requiredVars {
		if value := os.Getenv(varName); value == "" {
			log.Fatalf("環境変数 %s が設定されていません。", varName)
		}
	}
}

func initializeDatabase() *gorm.DB {
	db, err := migration.InitDB()
	if err != nil {
		log.Fatalf("データベースの初期化に失敗しました: %v", err)
	}
	return db
}

func migrateDatabaseIfNeeded(db *gorm.DB) {
	if getEnvOrDefault("RUN_MIGRATIONS", "false") == "true" {
		migration.Migrate(db)
		log.Println("データベース移行の実行が完了しました。")
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                         // 許可するオリジン
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE"},                     // リクエストメソッド
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"}, // リクエストヘッダに含めるヘッダ
		ExposeHeaders:    []string{"Content-Length"},                                            // レスポンスヘッダに含めるヘッダ
		AllowCredentials: true,                                                                  // クッキーを許可
		MaxAge:           12 * time.Hour,                                                        // 12時間
	}))

	initializeSwagger(router)

	classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController := initializeControllers(db)

	setupRoutes(router, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController)

	return router
}

func initializeSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/gin"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func startServer(router *gin.Engine) {
	port := getEnvOrDefault("PORT", "8080")
	log.Fatal(router.Run(":" + port))
}

func initializeControllers(db *gorm.DB) (*controllers.ClassBoardController, *controllers.ClassCodeController, *controllers.ClassScheduleController, *controllers.ClassUserController, *controllers.AttendanceController, *controllers.GoogleAuthController) {
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
	googleAuthService := services.NewGoogleAuthService()

	uploader := utils.NewAwsUploader()
	classBoardController := controllers.NewClassBoardController(services.NewClassBoardService(classBoardRepo), uploader)
	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	classScheduleController := controllers.NewClassScheduleController(classScheduleService)
	classUserController := controllers.NewClassUserController(classUserService)
	attendanceController := controllers.NewAttendanceController(attendanceService)
	googleAuthController := controllers.NewGoogleAuthController(googleAuthService)

	return classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController
}

func setupRoutes(router *gin.Engine, classBoardController *controllers.ClassBoardController, classCodeController *controllers.ClassCodeController, classScheduleController *controllers.ClassScheduleController, classUserController *controllers.ClassUserController, attendanceController *controllers.AttendanceController, googleAuthController *controllers.GoogleAuthController) {
	setupClassBoardRoutes(router, classBoardController)
	setupClassCodeRoutes(router, classCodeController)
	setupClassScheduleRoutes(router, classScheduleController)
	setupClassUserRoutes(router, classUserController)
	setupAttendanceRoutes(router, attendanceController)
	setupGoogleAuthRoutes(router, googleAuthController)

}

// setupGoogleAuthRoutes GoogleLoginのルートをセットアップする
func setupGoogleAuthRoutes(router *gin.Engine, controller *controllers.GoogleAuthController) {
	g := router.Group("/api/gin/auth/google")
	{
		g.GET("/login", controller.GoogleLoginHandler)
		g.GET("/callback", controller.GoogleAuthCallback)
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
func setupClassBoardRoutes(router *gin.Engine, controller *controllers.ClassBoardController) {
	cb := router.Group("/api/gin/cb")
	{
		cb.GET("/", controller.GetAllClassBoards)
		cb.GET("/:id", controller.GetClassBoardByID)
		cb.POST("/", controller.CreateClassBoard)
		cb.GET("/announced", controller.GetAnnouncedClassBoards)
		cb.PATCH("/:id", controller.UpdateClassBoard)
		cb.DELETE("/:id", controller.DeleteClassBoard)
	}
}

// setupClassCodeRoutes ClassCodeのルートをセットアップする
func setupClassCodeRoutes(router *gin.Engine, controller *controllers.ClassCodeController) {
	cc := router.Group("/api/gin/cc")
	{
		cc.GET("/checkSecretExists", controller.CheckSecretExists)
		cc.GET("/verifyClassCode", controller.VerifyClassCode)
	}
}

// setupClassScheduleRoutes ClassScheduleのルートをセットアップする
func setupClassScheduleRoutes(router *gin.Engine, controller *controllers.ClassScheduleController) {
	cs := router.Group("/api/gin/cs")
	{
		cs.GET("/", controller.GetAllClassSchedules)
		cs.GET("/:id", controller.GetClassScheduleByID)
		cs.POST("/", controller.CreateClassSchedule)
		cs.PATCH("/:id", controller.UpdateClassSchedule)
		cs.DELETE("/:id", controller.DeleteClassSchedule)
		cs.GET("/live", controller.GetLiveClassSchedules)
		cs.GET("/date", controller.GetClassSchedulesByDate)
	}
}

// setupClassUserRoutes ClassUserのルートをセットアップする
func setupClassUserRoutes(router *gin.Engine, controller *controllers.ClassUserController) {
	cu := router.Group("/api/gin/cu")
	{
		cu.PATCH("/:uid/:cid/:role", controller.ChangeUserRole)
	}
}

// setupAttendanceRoutes Attendanceのルートをセットアップする
func setupAttendanceRoutes(router *gin.Engine, controller *controllers.AttendanceController) {
	at := router.Group("/api/gin/at")
	{
		at.POST("/:cid/:uid/:csid", controller.CreateOrUpdateAttendance)
		at.GET("/:cid", controller.GetAllAttendances)
		at.GET("/attendance/:id", controller.GetAttendance)
		at.DELETE("/attendance/:id", controller.DeleteAttendance)
	}
}
