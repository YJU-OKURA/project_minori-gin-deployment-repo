package migration

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

func InitDB() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")

	// ポート番号が数値でない場合はデフォルトの3306を使用
	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Invalid port number. Using default port 3306. Error: %v", err)
		portInt = 3306
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, portInt, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
		&models.Role{},
		&models.Attendance{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
