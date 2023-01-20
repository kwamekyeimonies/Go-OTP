package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	datababase "github.com/kwamekyeimonies/Go-OTP/database"
)

var (
	server *gin.Engine
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	datababase.Database_Connection()

	server = gin.Default()

}

func main() {

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang Two Factor Authentication"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})

	})

	log.Fatal(
		server.Run(":9080"),
	)
}
