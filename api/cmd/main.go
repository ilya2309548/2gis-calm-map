package main

import (
	"log"

	_ "2gis-calm-map/api/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	"2gis-calm-map/api/internal/handler"

	swaggerFiles "github.com/swaggo/files"
)

// @title 2gis-calm-map API
// @version 1.0
// @description This is a sample server.
// @host localhost:8080
// @BasePath /

func main() {
	r := gin.Default()

	r.GET("users", handler.GetUsers)
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("start at :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
