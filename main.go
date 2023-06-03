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
	sites := models.GetAllSites()
	// for _, site := range sites {
	// 	fmt.Printf("%+v\n", site)
	// }

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
		c.JSON(
			http.StatusOK,
			gin.H{
				"message":  "parameters created",
				"Location": parameters.Location,
				"Days":     parameters.Days,
				"Types":    parameters.Types,
			})
	})

	r.POST("/api/v1/dfs", func(c *gin.Context) {
		var parameters models.Parameters
		err := c.BindJSON(&parameters)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}

		results := services.QueryFromDB(parameters)

		// create graph and siteMap(siteMap to store all the edges, vertex and weighted from the distance)
		services.PrepareGraphAndSiteMap(results)

		output := services.ImplementDFS(10)

		c.JSON(http.StatusOK, gin.H{
			"message":  "parameters created",
			"Location": parameters.Location,
			"Days":     parameters.Days,
			"Types":    parameters.Types,
			"Results":  output,
		})
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
