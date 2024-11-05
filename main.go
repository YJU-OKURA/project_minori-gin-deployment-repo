package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/middlewares"
	"github.com/go-redis/redis/v8"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/docs"
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

var (
	redisClient *redis.Client
	addr        = flag.String("addr", ":8080", "http service address")
)

func main() {
	configureGinMode()
	ensureEnvVariables()

	db := initializeDatabase()
	redisClient := initializeRedis()

	jwtService := services.NewJWTService()

	services.NewRoomManager(redisClient)

	router := setupRouter(db, jwtService)
	startServer(router)

	// Parse the flags passed to program
	flag.Parse()

	// start HTTP server
	log.Fatal(http.ListenAndServe(*addr, nil))
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

	requiredVars := []string{"POSTGRES_HOST", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DATABASE", "POSTGRES_PORT"}

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

// initializeRedis Redisを初期化する
func initializeRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
		//Password: redisPassword,
		DB: 0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Redisの初期化に失敗しました： %v\nREDIS_HOST: %s\nREDIS_PORT: %s\nREDIS_PASSWORD: %s",
			err, redisHost, redisPort, redisPassword)
	}

	redisClient = client
	return client
}

// setupRouter ルーターをセットアップする
func setupRouter(db *gorm.DB, jwtService services.JWTService) *gin.Engine {
	router := gin.Default()

	allowedOrigins := []string{
		"http://localhost:3000",
		"http://43.203.66.25",
		"http://10.0.9.193",
		"http://minori-next-lb-1326724168.ap-northeast-2.elb.amazonaws.com",
		"https://minoriedu.com",
		"https://minoriedu.com:80",
		"http://43.203.66.25/api/gin/swagger/index.html",
	}

	ignoredPaths := []string{
		"/api/gin/swagger/",
	}

	router.Use(globalErrorHandler)
	router.Use(CORS(allowedOrigins, ignoredPaths))
	initializeSwagger(router)
	userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController, lineAuthController, createClassController, chatController := initializeControllers(db, redisClient)

	setupRoutes(router, userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController, lineAuthController, createClassController, chatController, jwtService)
	return router
}

func globalErrorHandler(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		// You can log the errors here or send them to an external system
		// Optionally process different types of errors differently
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": c.Errors.String(),
		})
	}
}

// Swaggerのセキュリティ定義
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func initializeSwagger(router *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/gin"
	docs.SwaggerInfo.Title = "API Documentation"
	docs.SwaggerInfo.Description = "This is minori gin server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router.GET("/api/gin/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

//// initializeSwagger Swaggerを初期化する
//func initializeSwagger(router *gin.Engine) {
//	docs.SwaggerInfo.BasePath = "/api/gin"
//	router.GET("/api/gin/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
//}

func CORS(allowedOrigins []string, ignoredPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPath := c.Request.URL.Path

		for _, path := range ignoredPaths {
			if strings.HasPrefix(requestPath, path) {
				c.Next()
				return
			}
		}

		origin := c.Request.Header.Get("Origin")
		var isOriginAllowed bool
		for _, o := range allowedOrigins {
			if origin == o {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				isOriginAllowed = true
				break
			}
		}

		// 모든 API 경로 허용
		if strings.HasPrefix(requestPath, "/api/gin/") {
			isOriginAllowed = true
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// 특정 경로 무시
		for _, path := range ignoredPaths {
			if strings.HasPrefix(requestPath, path) {
				c.Next()
				return
			}
		}

		if isOriginAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PATCH, GET, PUT, DELETE, OPTIONS")
		}

		if c.Request.Method == "OPTIONS" {
			if isOriginAllowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		if !isOriginAllowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

// startServer サーバーを起動する
func startServer(router *gin.Engine) {
	srv := &http.Server{
		Addr:    ":" + getEnvOrDefault("PORT", "8080"),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// initializeControllers コントローラーを初期化する
func initializeControllers(db *gorm.DB, redisClient *redis.Client) (*controllers.UserController, *controllers.ClassBoardController, *controllers.ClassCodeController, *controllers.ClassScheduleController, *controllers.ClassUserController, *controllers.AttendanceController, *controllers.GoogleAuthController, *controllers.LINEAuthController, *controllers.ClassController, *controllers.ChatController) {
	userRepo := repositories.NewUserRepository(db)
	classRepo := repositories.NewClassRepository(db)
	classBoardRepo := repositories.NewClassBoardRepository(db)
	classCodeRepo := repositories.NewClassCodeRepository(db)
	classScheduleRepo := repositories.NewClassScheduleRepository(db)
	classUserRepo := repositories.NewClassUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	googleAuthRepo := repositories.NewGoogleAuthRepository(db)
	lineAuthRepo := repositories.NewLINEAuthRepository(db)

	userService := services.NewCreateUserService(userRepo)
	classBoardService := services.NewClassBoardService(classBoardRepo)
	classCodeService := services.NewClassCodeService(classCodeRepo)
	classUserService := services.NewClassUserService(classUserRepo, roleRepo)
	classScheduleService := services.NewClassScheduleService(classScheduleRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	googleAuthService := services.NewGoogleAuthService(googleAuthRepo)
	lineAuthService := services.NewLINEAuthService(lineAuthRepo)
	jwtService := services.NewJWTService()
	chatManager := services.NewRoomManager(redisClient)
	go manageChatRooms(db, chatManager)

	createClassService := services.NewCreateClassService(classRepo, classUserRepo, classCodeRepo, userRepo)

	uploader := utils.NewAwsUploader()
	userController := controllers.NewCreateUserController(userService)
	classBoardController := controllers.NewClassBoardController(classBoardService, uploader)
	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	classScheduleController := controllers.NewClassScheduleController(classScheduleService)
	classUserController := controllers.NewClassUserController(classUserService)
	attendanceController := controllers.NewAttendanceController(attendanceService)
	googleAuthController := controllers.NewGoogleAuthController(googleAuthService, jwtService)
	lineAuthController := controllers.NewLINEAuthController(lineAuthService, jwtService)
	createClassController := controllers.NewCreateClassController(createClassService, uploader)
	chatController := controllers.NewChatController(chatManager, redisClient)

	return userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, googleAuthController, lineAuthController, createClassController, chatController
}

// setupRoutes ルートをセットアップする
func setupRoutes(router *gin.Engine, userController *controllers.UserController, classBoardController *controllers.ClassBoardController, classCodeController *controllers.ClassCodeController, classScheduleController *controllers.ClassScheduleController, classUserController *controllers.ClassUserController, attendanceController *controllers.AttendanceController, googleAuthController *controllers.GoogleAuthController, lineAuthController *controllers.LINEAuthController, createClassController *controllers.ClassController, chatController *controllers.ChatController, jwtService services.JWTService) {
	setupUserRoutes(router, userController, jwtService)
	setupClassBoardRoutes(router, classBoardController, jwtService)
	setupClassCodeRoutes(router, classCodeController, jwtService)
	setupClassScheduleRoutes(router, classScheduleController, jwtService)
	setupClassUserRoutes(router, classUserController, jwtService)
	setupAttendanceRoutes(router, attendanceController, jwtService)
	setupGoogleAuthRoutes(router, googleAuthController)
	setupLINEAuthRoutes(router, lineAuthController)
	setupCreateClassRoutes(router, createClassController, jwtService)
	setupChatRoutes(router, chatController, jwtService)
}

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupUserRoutes(router *gin.Engine, controller *controllers.UserController, jwtService services.JWTService) {
	u := router.Group("/api/gin/u")
	u.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		u.GET(":userID/applying-classes", controller.GetApplyingClasses)
		u.GET("search", controller.SearchByName)
		u.DELETE(":userID/delete", controller.RemoveUserFromService)
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassBoardRoutes(router *gin.Engine, controller *controllers.ClassBoardController, jwtService services.JWTService) {
	cb := router.Group("/api/gin/cb")
	cb.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		cb.GET("", controller.GetAllClassBoards)
		cb.GET(":id", controller.GetClassBoardByID)
		cb.GET("announced", controller.GetAnnouncedClassBoards)

		// TODO: フロントエンド側の実装が完了したら、削除
		cb.POST("", controller.CreateClassBoard)
		cb.PATCH(":id/:cid/:uid", controller.UpdateClassBoard)
		cb.DELETE(":id", controller.DeleteClassBoard)

		cb.GET("subscribe", controller.SubscribeClassBoardUpdates)
		cb.GET("search", controller.SearchClassBoards)
	}
}

// setupClassCodeRoutes ClassCodeのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassCodeRoutes(router *gin.Engine, controller *controllers.ClassCodeController, jwtService services.JWTService) {
	cc := router.Group("/api/gin/cc")
	cc.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		cc.GET("checkSecretExists", controller.CheckSecretExists)
		cc.GET("verifyClassCode", controller.VerifyClassCode)
		cc.GET("verifyAndRequestAccess", controller.VerifyAndRequestAccess)
	}
}

// setupClassScheduleRoutes ClassScheduleのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassScheduleRoutes(router *gin.Engine, controller *controllers.ClassScheduleController, jwtService services.JWTService) {
	cs := router.Group("/api/gin/cs")
	cs.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		cs.GET("", controller.GetAllClassSchedules)
		cs.GET(":id", controller.GetClassScheduleByID)

		// TODO: フロントエンド側の実装が完了したら、削除
		cs.POST("", controller.CreateClassSchedule)
		cs.PATCH(":id", controller.UpdateClassSchedule)
		cs.DELETE(":id", controller.DeleteClassSchedule)
		cs.GET("live", controller.GetLiveClassSchedules)
		cs.GET("date", controller.GetClassSchedulesByDate)
	}
}

// setupGoogleAuthRoutes GoogleLoginのルートをセットアップする
func setupGoogleAuthRoutes(router *gin.Engine, controller *controllers.GoogleAuthController) {
	g := router.Group("/api/gin/auth/google")
	{
		g.GET("login", controller.GoogleLoginHandler)
		g.POST("process", controller.ProcessAuthCode)
		g.POST("refresh-token", controller.RefreshAccessTokenHandler)
	}
}

func setupLINEAuthRoutes(router *gin.Engine, controller *controllers.LINEAuthController) {
	g := router.Group("/api/gin/auth/line")
	{
		g.GET("login", controller.LINELoginHandler)
		g.POST("process", controller.ProcessAuthCode)
		g.POST("refresh-token", controller.RefreshAccessTokenHandler)
	}
}

// setupCreateClassRoutes CreateClassのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupCreateClassRoutes(router *gin.Engine, controller *controllers.ClassController, jwtService services.JWTService) {
	cl := router.Group("/api/gin/cl")
	cl.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		cl.GET(":cid", controller.GetClass)
		cl.POST("create", controller.CreateClass)
		cl.PATCH(":uid/:cid", controller.UpdateClass)
		cl.DELETE(":uid/:cid", controller.DeleteClass)
	}
}

// setupClassUserRoutes ClassUserのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassUserRoutes(router *gin.Engine, controller *controllers.ClassUserController, jwtService services.JWTService) {
	cu := router.Group("/api/gin/cu")
	cu.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		// TODO: フロントエンド側の実装が完了したら、削除
		cu.GET("class/:cid/members", controller.GetClassMembers)

		userRoutes := cu.Group(":uid")
		{
			userRoutes.GET(":cid/info", controller.GetUserClassUserInfo)
			userRoutes.GET("classes", controller.GetUserClasses)
			userRoutes.GET("favorite-classes", controller.GetFavoriteClasses)
			userRoutes.GET("classes/by-role", controller.GetUserClassesByRole)
			userRoutes.PATCH(":cid/role/:roleName", controller.ChangeUserRole)
			userRoutes.PATCH(":cid/toggle-favorite", controller.ToggleFavorite)
			userRoutes.PUT(":cid/:rename", controller.UpdateUserName)
			userRoutes.DELETE(":cid/remove", controller.RemoveUserFromClass)
			userRoutes.GET("classes/search", controller.SearchUserClassesByName)
		}
	}
}

// setupAttendanceRoutes Attendanceのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupAttendanceRoutes(router *gin.Engine, controller *controllers.AttendanceController, jwtService services.JWTService) {
	at := router.Group("/api/gin/at")
	at.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		at.POST("", controller.CreateOrUpdateAttendance)
		at.GET(":cid", controller.GetAllAttendances)
		at.GET("attendance/:id", controller.GetAttendance)
		at.DELETE("attendance/:id", controller.DeleteAttendance)
	}
}

// setupChatRoutes Chatのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupChatRoutes(router *gin.Engine, chatController *controllers.ChatController, jwtService services.JWTService) {
	chat := router.Group("/api/gin/chat")
	chat.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		chat.POST("create-room/:scheduleId", chatController.CreateChatRoom)
		chat.GET("room/:scheduleId/:userId", chatController.HandleChatRoom)
		chat.POST("room/:scheduleId", chatController.PostToChatRoom)
		chat.DELETE("room/:scheduleId", chatController.DeleteChatRoom)
		chat.GET("stream/:scheduleId", chatController.StreamChat)
		chat.GET("messages/:roomid", chatController.GetChatMessages)
		chat.POST("dm/:senderId/:receiverId", chatController.SendDirectMessage)
		chat.GET("dm/:senderId/:receiverId", chatController.GetDirectMessages)
		chat.DELETE("dm/:senderId/:receiverId", chatController.DeleteDirectMessages)
	}
}

func manageChatRooms(db *gorm.DB, chatManager *services.Manager) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		now := time.Now()
		var schedules []models.ClassSchedule

		// 수업 시작 5분 전과 수업 종료 10분 후에 채팅방 상태를 확인
		db.Where("started_at <= ? AND started_at >= ?", now.Add(5*time.Minute), now).
			Or("ended_at <= ? AND ended_at >= ?", now, now.Add(-10*time.Minute)).Find(&schedules)

		for _, schedule := range schedules {
			roomID := fmt.Sprintf("class_%d", schedule.ID)
			// 종료 10분 후 검사를 위해 ended_at에 10분을 더해 현재 시간과 비교
			if now.After(schedule.EndedAt.Add(10 * time.Minute)) {
				chatManager.DeleteBroadcast(roomID)
			}
		}
	}
}
