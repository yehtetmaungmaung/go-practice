package main

import (
	"example/web-service-gin/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/albums", controllers.GetAlbums)
	router.POST("/albums", controllers.PostAlbums)

	router.Run("localhost:8080")
}
