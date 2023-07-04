package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"mehmetfd.dev/chessu-backend/models"
)

var (
	DB *gorm.DB
)

func InitDB() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// Assign the database instance to the global variable
	DB = db

	// Migrate the schema
	db.AutoMigrate(&models.AppUser{}, &models.Membership{})

}
