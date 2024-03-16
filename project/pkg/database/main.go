package database

import (
	"fmt"
	"os"
	"time"

	"bitbucket.org/dyfrag-internal/mass-media-core/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetUpDB() {
	dsn := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" port=" + os.Getenv("POSTGRES_PORT")

	retries := 5
	var err error

	for retries > 0 {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Println("Couldn't connect to database")
		retries--
		time.Sleep(5 * time.Second)
	}

	DB.AutoMigrate(
		&models.User{},
	)
}
