package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB


func init() {
	InitGormDB()
}

func InitGormDB() {
	dsn := getDSN()
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("connect to database failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("unable to get DB instance %v", err)
	}

	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
}

func getDSN() string {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "admin")
	pass := getEnv("DB_PASSWORD", "1234")
	dbname := getEnv("DB_NAME", "gcs_db")
	timezone := getEnv("DB_TIMEZONE", "Asia/Shanghai")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		host, port, user, pass, dbname, timezone,
	)
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func GetGormDB() *gorm.DB {
	return db
}
