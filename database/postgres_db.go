package datababase

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kwamekyeimonies/Go-OTP/model"
)

var DB *gorm.DB

func Database_Connection() {

	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbName := os.Getenv("DATABASE_NAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")

	//Credentials
	ConnectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", dbHost, dbPort, dbUser, dbName, dbPass)
	fmt.Println(ConnectionString)
	db, err := gorm.Open("postgres", ConnectionString)

	if err != nil {
		fmt.Println("Error Connecting to Database:", err)
	} else {
		fmt.Println("Database Connected Succesfully....")
	}

	// //Table Creation in PostgreSQL

	db.AutoMigrate(
		&model.User{},
		&model.LoginUserInput{},
		&model.OTPInput{},
		&model.RegisterUserInput{},
	)

	DB = db
}
