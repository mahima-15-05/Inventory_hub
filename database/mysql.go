package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(){
	err:= godotenv.Load()
	if err!=nil{
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")


	// create db string from .env file's variables 
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host,port, dbname)

	// connection to database 
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err!=nil{
		log.Fatal("Failed to connect with database")
	}
	DB = db
}