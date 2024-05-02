package migration

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	// ポート番号が数値でない場合はデフォルトの3306を使用
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Invalid port number. Using default port 5432. Error: %v", err)
		portInt = 5432
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Tokyo", host, user, pass, dbName, portInt)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get gennric database: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.ClassUser{},
		&models.ClassBoard{},
		&models.ClassCode{},
		&models.ClassSchedule{},
		&models.Attendance{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
