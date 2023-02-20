package db

import (
	"fmt"
	"go-usermgmt-grpc/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Init() *gorm.DB {

	DB_Host := "192.168.10.21"
	DbUser := os.Getenv("POSTGRES_DEVELOP_DB_USERNAME")
	DbPassword := os.Getenv("POSTGRES_DEVELOP_DB_PASSWORD")
	DbName := os.Getenv("POSTGRES_DEVELOP_DB_NAME")
	DbPort := 5432
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", DB_Host, DbUser, DbPassword, DbName, DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the Database: %v", err)
	}
	db.AutoMigrate(&models.Accounts{})
	return db
}
