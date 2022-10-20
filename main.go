package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	routes "github.com/yaduvendra/E-commerce/routes"
)

func main(){

	er := godotenv.Load(".env")

	if er != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == ""{
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	
	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api1",func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for api1"})
	})

	router.GET("/api2",func(c *gin.Context){
		c.JSON(200, gin.H{"success":"Access granted for api2"})
	})

	router.Run(":"+port)
}