package main

import (
	"net/http"
	"t3/m/v2/models"
	"t3/m/v2/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	models.ConnectDatabase()
	//data.LoadSites()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"message": "pong",
			})
	})

	r.POST("/api/v1/parameters", func(c *gin.Context) {
		var parameters models.Parameters
		err := c.BindJSON(&parameters)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "bad request",
				})
			return
		}
		results := services.QueryFromDB(parameters)
		c.JSON(
			http.StatusOK,
			gin.H{
				"Location": parameters.Location,
				"Days":     parameters.Days,
				"Types":    parameters.Types,
				"Results":  results,
			})
	})

	r.POST("/api/v1/tsp", func(c *gin.Context) {
		var parameters models.Parameters
		err := c.BindJSON(&parameters)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "bad request",
				})
			return
		}
		results := services.QueryFromDB(parameters)

		points := make([]models.Point, len(results))
		for i, site := range results {
			points[i] = models.SiteToPoint(site)
		}

		tsp := services.TSP(points)
		c.JSON(
			http.StatusOK,
			gin.H{
				"Location": parameters.Location,
				"Days":     parameters.Days,
				"Types":    parameters.Types,
				"Results":  results,
				"Points":   points,
				"TSP":      tsp,
			})
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
