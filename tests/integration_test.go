//go:build integration
// +build integration

package tests

import (
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/controllers"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/migration"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/repositories"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var db *gorm.DB
var router *gin.Engine

func TestMain(m *testing.M) {
	// Set up the environment
	if err := os.Setenv("GIN_MODE", gin.TestMode); err != nil {
		log.Fatalf("failed to set GIN_MODE: %v", err)
	}
	var err error
	db, err = migration.InitDB()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	jwtService := services.NewJWTService()

	// Set up the router
	router = setupRouter(db, jwtService)

	// Run the tests
	code := m.Run()

	// Clean up
	if err := db.Exec("DROP DATABASE test_db").Error; err != nil {
		log.Fatalf("failed to drop test database: %v", err)
	}
	os.Exit(code)
}

func setupRouter(db *gorm.DB, jwtService services.JWTService) *gin.Engine {
	router := gin.Default()
	// Initialize controllers and routes here
	userController := controllers.NewCreateUserController(services.NewCreateUserService(repositories.NewUserRepository(db)))
	router.GET("/api/gin/u/:userID/applying-classes", userController.GetApplyingClasses)
	return router
}

func TestGetApplyingClasses(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/gin/u/1/applying-classes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
