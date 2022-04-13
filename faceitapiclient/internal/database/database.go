package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(host, username, password, database string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Europe/Berlin",
		host, username, password, database, port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
