package database

import (
	"dashboard-server/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dashboard.db"
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	sqlDB, err := DB.DB()

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.Dashboard{}, &models.Widget{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	var dashboard models.Dashboard
	result := DB.First(&dashboard)

	if result.Error != nil {
		dashboard = models.Dashboard{Name: "Default Dashboard"}
		DB.Create(&dashboard)
		log.Println("Created default dashboard")
	}

	log.Println("Database connected and migrated successfully")
}
