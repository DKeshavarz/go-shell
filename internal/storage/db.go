package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbHost = "localhost"
	dbPort = "5432"
	dbName = "go_shell"
	dbUser = "postgres"
	dbPass = "password"
)


var db *gorm.DB

func init() {
	fmt.Println("========================================================================================================================")
	fmt.Println("connecting to database...")

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPass)

	sesssion, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Silent) ,
	})
	if err != nil {
		err = fmt.Errorf("unable to connect to database: %w", err)
		panic(err)
	}

	db = sesssion

	fmt.Printf("connected to database: %s\n", dsn)
	fmt.Println("========================================================================================================================")
}

func GetDB() *gorm.DB {
	return db
}
