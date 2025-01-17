package database

import (
	"fmt"
	"inventoryapp/models"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func StartDB() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("PGUSER")
	host := os.Getenv("PGHOST")
	password := os.Getenv("PGPASSWORD")
	dbPort := os.Getenv("PGPORT")
	dbName := os.Getenv("PGDATABASE")
	sslMode := os.Getenv("PGSSLMODE")

	config := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s", host, user, password, dbPort, dbName, sslMode)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database", err)
	}

	fmt.Println("successfully connecting to database")
	db.Debug().AutoMigrate(
		&models.Users{},
		&models.Products{},
		&models.IncomingItems{},
		&models.OutgoingItems{},
	)
}

func GetDB() *gorm.DB {
	return db
}
