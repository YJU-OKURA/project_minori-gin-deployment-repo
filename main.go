package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/middlewares"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"

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

var redisClient *redis.Client

func main() {
	configureGinMode()
	ensureEnvVariables()

	db := initializeDatabase()
	redisClient := initializeRedis()

	jwtService := services.NewJWTService()

	services.NewRoomManager(redisClient)
	migrateDatabaseIfNeeded(db)

	router := setupRouter(db, jwtService)
	startServer(router)

}

// Initialize gRPC server
//func startGRPCServer() {
//	lis, err := net.Listen("tcp", ":50051")
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	grpcServer := grpc.NewServer()
//	//	Register gRPC services here
//	log.Printf("server listening at %v", lis.Addr())
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}

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

// migrateDatabaseIfNeeded データベースを移行する
func migrateDatabaseIfNeeded(db *gorm.DB) {
	if getEnvOrDefault("RUN_MIGRATIONS", "false") == "true" {
		migration.Migrate(db)
		log.Println("データベース移行の実行が完了しました。")
	}
}

// setupRouter ルーターをセットアップする
func setupRouter(db *gorm.DB, jwtService services.JWTService) *gin.Engine {
	router := gin.Default()
	router.Use(globalErrorHandler)
	router.Use(CORS())
	initializeSwagger(router)
	userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController, chatController, liveClassController := initializeControllers(db, redisClient)

	setupRoutes(router, userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController, chatController, liveClassController, jwtService)
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
func initializeControllers(db *gorm.DB, redisClient *redis.Client) (*controllers.UserController, *controllers.ClassBoardController, *controllers.ClassCodeController, *controllers.ClassScheduleController, *controllers.ClassUserController, *controllers.AttendanceController, services.ClassUserService, *controllers.GoogleAuthController, *controllers.ClassController, *controllers.ChatController, *controllers.LiveClassController) {
	userRepo := repositories.NewUserRepository(db)
	classRepo := repositories.NewClassRepository(db)
	classBoardRepo := repositories.NewClassBoardRepository(db)
	classCodeRepo := repositories.NewClassCodeRepository(db)
	classScheduleRepo := repositories.NewClassScheduleRepository(db)
	classUserRepo := repositories.NewClassUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	googleAuthRepo := repositories.NewGoogleAuthRepository(db)

	userService := services.NewCreateUserService(userRepo)
	classBoardService := services.NewClassBoardService(classBoardRepo)
	classCodeService := services.NewClassCodeService(classCodeRepo)
	classUserService := services.NewClassUserService(classUserRepo, roleRepo)
	classScheduleService := services.NewClassScheduleService(classScheduleRepo)
	attendanceService := services.NewAttendanceService(attendanceRepo)
	googleAuthService := services.NewGoogleAuthService(googleAuthRepo)
	jwtService := services.NewJWTService()
	chatManager := services.NewRoomManager(redisClient)
	liveClassService := services.NewLiveClassService(classUserRepo, redisClient)
	go manageChatRooms(db, chatManager)

	createClassService := services.NewCreateClassService(classRepo, classUserRepo, classCodeRepo)

	uploader := utils.NewAwsUploader()
	userController := controllers.NewCreateUserController(userService)
	classBoardController := controllers.NewClassBoardController(classBoardService, uploader)
	classCodeController := controllers.NewClassCodeController(classCodeService, classUserService)
	classScheduleController := controllers.NewClassScheduleController(classScheduleService)
	classUserController := controllers.NewClassUserController(classUserService)
	attendanceController := controllers.NewAttendanceController(attendanceService)
	googleAuthController := controllers.NewGoogleAuthController(googleAuthService, jwtService)
	createClassController := controllers.NewCreateClassController(createClassService, uploader)
	chatController := controllers.NewChatController(chatManager, redisClient)
	liveClassController := controllers.NewLiveClassController(liveClassService)

	return userController, classBoardController, classCodeController, classScheduleController, classUserController, attendanceController, classUserService, googleAuthController, createClassController, chatController, liveClassController
}

// setupRoutes ルートをセットアップする
func setupRoutes(router *gin.Engine, userController *controllers.UserController, classBoardController *controllers.ClassBoardController, classCodeController *controllers.ClassCodeController, classScheduleController *controllers.ClassScheduleController, classUserController *controllers.ClassUserController, attendanceController *controllers.AttendanceController, classUserService services.ClassUserService, googleAuthController *controllers.GoogleAuthController, createClassController *controllers.ClassController, chatController *controllers.ChatController, liveClassController *controllers.LiveClassController, jwtService services.JWTService) {
	setupUserRoutes(router, userController, jwtService)
	setupClassBoardRoutes(router, classBoardController, classUserService, jwtService)
	setupClassCodeRoutes(router, classCodeController, jwtService)
	setupClassScheduleRoutes(router, classScheduleController, classUserService, jwtService)
	setupClassUserRoutes(router, classUserController, classUserService, jwtService)
	setupAttendanceRoutes(router, attendanceController, classUserService, jwtService)
	setupGoogleAuthRoutes(router, googleAuthController)
	setupCreateClassRoutes(router, createClassController, jwtService)
	setupChatRoutes(router, chatController, jwtService)
	setupLiveClassRoutes(router, liveClassController)
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
	}
}

// setupClassBoardRoutes ClassBoardのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassBoardRoutes(router *gin.Engine, controller *controllers.ClassBoardController, classUserService services.ClassUserService, jwtService services.JWTService) {
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
	}
}

// setupClassScheduleRoutes ClassScheduleのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupClassScheduleRoutes(router *gin.Engine, controller *controllers.ClassScheduleController, classUserService services.ClassUserService, jwtService services.JWTService) {
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
func setupClassUserRoutes(router *gin.Engine, controller *controllers.ClassUserController, classUserService services.ClassUserService, jwtService services.JWTService) {
	cu := router.Group("/api/gin/cu")
	cu.Use(middlewares.TokenAuthMiddleware(jwtService))
	{
		// TODO: フロントエンド側の実装が完了したら、削除
		cu.GET("class/:cid/:role/members", controller.GetClassMembers)

		userRoutes := cu.Group(":uid")
		{
			userRoutes.GET(":cid/info", controller.GetUserClassUserInfo)
			userRoutes.GET("classes", controller.GetUserClasses)
			userRoutes.GET("favorite-classes", controller.GetFavoriteClasses)
			userRoutes.GET("classes/:roleID", controller.GetUserClassesByRole)
			userRoutes.PATCH(":cid/:roleID", controller.ChangeUserRole)
			userRoutes.PATCH(":cid/toggle-favorite", controller.ToggleFavorite)
			userRoutes.PUT(":cid/:rename", controller.UpdateUserName)
			userRoutes.DELETE(":cid/remove", controller.RemoveUserFromClass)
		}

		// TODO: フロントエンド側の実装が完了したら、コメントアウトを外す
		//protected := cu.Group("/:uid/:cid")
		//protected.Use(middlewares.AdminMiddleware(classUserService), middlewares.AssistantMiddleware(classUserService))
		//{
		//	protected.PATCH("/:uid/:cid/:role", controller.ChangeUserRole)
		//}
	}
}

// setupAttendanceRoutes Attendanceのルートをセットアップする
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func setupAttendanceRoutes(router *gin.Engine, controller *controllers.AttendanceController, classUserService services.ClassUserService, jwtService services.JWTService) {
	at := router.Group("/api/gin/at")
	at.Use(middlewares.TokenAuthMiddleware(jwtService))
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

func setupLiveClassRoutes(router *gin.Engine, liveClassController *controllers.LiveClassController) {
	live := router.Group("/api/gin/live")
	{
		live.GET("screen_share/:uid/:cid", liveClassController.GetScreenShareInfo)
		live.POST("screen_share/start/:cid", liveClassController.StartScreenShare)
	}
}
