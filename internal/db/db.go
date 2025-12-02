package db

import (
	"AppointmentAPI/internal/models"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDb() *gorm.DB {

	dsn := "root:root@tcp(127.0.0.1:3306)/appointmentdb?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Coach{},
		&models.Availability{},
		&models.Booking{},
	); err != nil {
		log.Fatalf("auto-migrate failed: %v", err)
	}

	log.Println("database connected and migrated")
	return db
}
