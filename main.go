package main

import (
	"net/http"
	"t3/m/v2/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			})
	})

	r.Run()
}
