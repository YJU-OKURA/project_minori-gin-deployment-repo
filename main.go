package main

import (
	"log"
	"os"

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
)

func main() {
	configureGinMode()
	ensureEnvVariables()

	db := initializeDatabase()
	migrateDatabaseIfNeeded(db)

	router := setupRouter(db)
	startServer(router)
}

// configureGinMode Ginのモードを設定する
func configureGinMode() {
	ginMode := getEnvOrDefault("GIN_MODE", gin.ReleaseMode)
	gin.SetMode(ginMode)
}

// getEnvOrDefault 環境変数が設定されていない場合はデフォルト値を返す
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ensureEnvVariables 環境変数が設定されているか確認する
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

// initializeDatabase データベースを初期化する
func initializeDatabase() *gorm.DB {
	db, err := migration.InitDB()
	if err != nil {
		log.Fatalf("データベースの初期化に失敗しました: %v", err)
	}
	return db
}

// migrateDatabaseIfNeeded データベースを移行する
func migrateDatabaseIfNeeded(db *gorm.DB) {
	if getEnvOrDefault("RUN_MIGRATIONS", "false") == "true" {
		migration.Migrate(db)
		log.Println("データベース移行の実行が完了しました。")
	}
}

// setupRouter ルーターをセットアップする
func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(CORS())

	initializeSwagger(router)

	userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController := initializeControllers(db)

	setupRoutes(router, userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController)

	return router
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// initializeSwagger Swaggerを初期化する
func initializeSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/gin"
	router.GET("/api/gin/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

// startServer サーバーを起動する
func startServer(router *gin.Engine) {
	port := getEnvOrDefault("PORT", "8080")
	log.Fatal(router.Run(":" + port))
}

// initializeControllers コントローラーを初期化する
func initializeControllers(db *gorm.DB) (*controllers.UserController, *controllers.ClassBoardController, *controllers.ClassCodeController, *controllers.ClassScheduleController, *controllers.ClassUserController, *controllers.AttendanceController, services.ClassUserService, *controllers.GoogleAuthController, *controllers.ClassController) {
	userRepo := repositories.NewUserRepository(db)
	createClassRepo := repositories.NewCreateClassRepository(db)
	classBoardRepo := repositories.NewClassBoardRepository(db)
	classCodeRepo := repositories.NewClassCodeRepository(db)
	classScheduleRepo := repositories.NewClassScheduleRepository(db)
	classUserRepo := repositories.NewClassUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	googleAuthRepo := repositories.NewGoogleAuthRepository(db)

	userService := services.NewCreateUserService(userRepo)
	createClassService := services.NewCreateClassService(createClassRepo)
	classBoardService := services.NewClassBoardService(classBoardRepo)
	classCodeService := services.NewClassCodeService(classCodeRepo)
	classUserService := services.NewClassUserService(classUserRepo, roleRepo)
	classScheduleService := services.NewClassScheduleService(classScheduleRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	googleAuthService := services.NewGoogleAuthService(googleAuthRepo)

	uploader := utils.NewAwsUploader()
	userController := controllers.NewCreateUserController(userService)
	createClassController := controllers.NewCreateClassController(createClassService, uploader)
	classBoardController := controllers.NewClassBoardController(classBoardService, uploader)
	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	classScheduleController := controllers.NewClassScheduleController(classScheduleService)
	classUserController := controllers.NewClassUserController(classUserService)
	attendanceController := controllers.NewAttendanceController(attendanceService)
	googleAuthController := controllers.NewGoogleAuthController(googleAuthService)

	return userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController
}

// setupRoutes ルートをセットアップする
func setupRoutes(router *gin.Engine, userController *controllers.UserController, classBoardController *controllers.ClassBoardController, classCodeController *controllers.ClassCodeController, classScheduleController *controllers.ClassScheduleController, classUserController *controllers.ClassUserController, attendanceController *controllers.AttendanceController, classUserService services.ClassUserService, googleAuthController *controllers.GoogleAuthController, createClassController *controllers.ClassController) {
	setupUserRoutes(router, userController)
	setupClassBoardRoutes(router, classBoardController, classUserService)
	setupClassCodeRoutes(router, classCodeController)
	setupClassScheduleRoutes(router, classScheduleController, classUserService)
	setupClassUserRoutes(router, classUserController, classUserService)
	setupAttendanceRoutes(router, attendanceController, classUserService)
	setupGoogleAuthRoutes(router, googleAuthController)
	setupCreateClassRoutes(router, createClassController)
}

func setupUserRoutes(router *gin.Engine, controller *controllers.UserController) {
	u := router.Group("/api/gin/u")
	{
		u.GET(":userID/applying-classes", controller.GetApplyingClasses)
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
func setupClassBoardRoutes(router *gin.Engine, controller *controllers.ClassBoardController, classUserService services.ClassUserService) {
	cb := router.Group("/api/gin/cb")
	{
		cb.GET("", controller.GetAllClassBoards)
		cb.GET(":id", controller.GetClassBoardByID)
		cb.GET("announced", controller.GetAnnouncedClassBoards)

		// TODO: フロントエンド側の実装が完了したら、削除
		cb.POST("", controller.CreateClassBoard)
		cb.PATCH(":id", controller.UpdateClassBoard)
		cb.DELETE(":id", controller.DeleteClassBoard)

		// TODO: フロントエンド側の実装が完了したら、コメントアウトを外す
		//protected := cb.Group("/:uid/:cid")
		//protected.Use(middlewares.AdminMiddleware(classUserService), middlewares.AssistantMiddleware(classUserService))
		//{
		//	protected.POST("/", controller.CreateClassBoard)
		//	protected.PATCH("/:id", controller.UpdateClassBoard)
		//	protected.DELETE("/:id", controller.DeleteClassBoard)
		//}
	}
}

// setupClassCodeRoutes ClassCodeのルートをセットアップする
func setupClassCodeRoutes(router *gin.Engine, controller *controllers.ClassCodeController) {
	cc := router.Group("/api/gin/cc")
	{
		cc.GET("checkSecretExists", controller.CheckSecretExists)
		cc.GET("verifyClassCode", controller.VerifyClassCode)
	}
}

// setupClassScheduleRoutes ClassScheduleのルートをセットアップする
func setupClassScheduleRoutes(router *gin.Engine, controller *controllers.ClassScheduleController, classUserService services.ClassUserService) {
	cs := router.Group("/api/gin/cs")
	{
		cs.GET("", controller.GetAllClassSchedules)
		cs.GET(":id", controller.GetClassScheduleByID)

		// TODO: フロントエンド側の実装が完了したら、削除
		cs.POST("", controller.CreateClassSchedule)
		cs.PATCH(":id", controller.UpdateClassSchedule)
		cs.DELETE(":id", controller.DeleteClassSchedule)
		cs.GET("live", controller.GetLiveClassSchedules)
		cs.GET("date", controller.GetClassSchedulesByDate)

		// TODO: フロントエンド側の実装が完了したら、コメントアウトを外す
		//protected := cs.Group("/:uid/:cid")
		//protected.Use(middlewares.AdminMiddleware(classUserService), middlewares.AssistantMiddleware(classUserService))
		//{
		//	protected.POST("/", controller.CreateClassSchedule)
		//	protected.PATCH("/:id", controller.UpdateClassSchedule)
		//	protected.DELETE("/:id", controller.DeleteClassSchedule)
		//	protected.GET("/live", controller.GetLiveClassSchedules)
		//	protected.GET("/date", controller.GetClassSchedulesByDate)
		//}
	}
}

// setupGoogleAuthRoutes GoogleLoginのルートをセットアップする
func setupGoogleAuthRoutes(router *gin.Engine, controller *controllers.GoogleAuthController) {
	g := router.Group("/api/gin/auth/google")
	{
		g.GET("login", controller.GoogleLoginHandler)
		g.GET("callback", controller.GoogleAuthCallback)
	}
}

//setupCreateClassRoutes CreateClassのルートをセットアップする

func setupCreateClassRoutes(router *gin.Engine, controller *controllers.ClassController) {
	cs := router.Group("/api/gin/cs")
	{
		cs.POST("create", controller.CreateClass)

	}
}

// setupClassUserRoutes ClassUserのルートをセットアップする
func setupClassUserRoutes(router *gin.Engine, controller *controllers.ClassUserController, classUserService services.ClassUserService) {
	cu := router.Group("/api/gin/cu")
	{
		// TODO: フロントエンド側の実装が完了したら、削除
		cu.GET(":uid/classes", controller.GetUserClasses)

		cu.PATCH(":uid/:cid/:role", controller.ChangeUserRole)

		cu.PUT(":uid/:cid/:rename", controller.UpdateUserName)

		// TODO: フロントエンド側の実装が完了したら、コメントアウトを外す
		//protected := cu.Group("/:uid/:cid")
		//protected.Use(middlewares.AdminMiddleware(classUserService), middlewares.AssistantMiddleware(classUserService))
		//{
		//	protected.PATCH("/:uid/:cid/:role", controller.ChangeUserRole)
		//}
	}
	cm := router.Group("/api/gin/cm")
	{
		cm.GET(":cid/members", controller.GetClassMembers)
	}
}

// setupAttendanceRoutes Attendanceのルートをセットアップする
func setupAttendanceRoutes(router *gin.Engine, controller *controllers.AttendanceController, classUserService services.ClassUserService) {
	at := router.Group("/api/gin/at")
	{
		// TODO: フロントエンド側の実装が完了したら、削除
		at.POST(":cid/:uid/:csid", controller.CreateOrUpdateAttendance)
		at.GET(":cid", controller.GetAllAttendances)
		at.GET("attendance/:id", controller.GetAttendance)
		at.DELETE("attendance/:id", controller.DeleteAttendance)

		// TODO: フロントエンド側の実装が完了したら、コメントアウトを外す
		//protected := at.Group("/:uid/:cid")
		//protected.Use(middlewares.AdminMiddleware(classUserService))
		//{
		//	protected.POST("/:cid/:uid/:csid", controller.CreateOrUpdateAttendance)
		//	protected.GET("/:cid", controller.GetAllAttendances)
		//	protected.GET("/attendance/:id", controller.GetAttendance)
		//	protected.DELETE("/attendance/:id", controller.DeleteAttendance)
		//}
	}
}
