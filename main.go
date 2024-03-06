package main

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	docs "github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/migration"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"log"
	"os"
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
	requiredVars := []string{"MYSQL_HOST", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE", "MYSQL_PORT"}
	for _, varName := range requiredVars {
		if value := os.Getenv(varName); value == "" {
			log.Fatalf("Environment variable %s not set", varName)
		}
	}
}

func initializeDatabase() *gorm.DB {
	db, err := migration.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	return db
}

func migrateDatabaseIfNeeded(db *gorm.DB) {
	if getEnvOrDefault("RUN_MIGRATIONS", "false") == "true" {
		migration.Migrate(db)
		log.Println("Database migrations executed")
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	initializeSwagger(router)

	classBoardController, classCodeController, classScheduleController, attendanceController := initializeControllers(db)

	setupRoutes(router, classBoardController, classCodeController, classScheduleController, attendanceController)

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

func initializeControllers(db *gorm.DB) (*controllers.ClassBoardController, *controllers.ClassCodeController, *controllers.ClassScheduleController, *controllers.AttendanceController) {
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

	return classBoardController, classCodeController, classScheduleController, attendanceController
}

func setupRoutes(router *gin.Engine, classBoardController *controllers.ClassBoardController, classCodeController *controllers.ClassCodeController, classScheduleController *controllers.ClassScheduleController, attendanceController *controllers.AttendanceController) {
	setupClassBoardRoutes(router, classBoardController)
	setupClassCodeRoutes(router, classCodeController)
	setupClassScheduleRoutes(router, classScheduleController)
	setupAttendanceRoutes(router, attendanceController)
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
