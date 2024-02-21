package migration

import (
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitDB() *gorm.DB {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.ClassUser{},
		&models.GroupBoard{},
		&models.GroupCode{},
		&models.GroupSchedule{},
		&models.Role{},
		&models.Attendance{},
	)
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
