package main

import (
	"log"
	"t3/m/v2/data"
	"t3/m/v2/models"
	"t3/m/v2/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	log.Println("test")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	r.Use(cors.New(config))

	models.ConnectDatabase()

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(
	// 		http.StatusOK,
	// 		gin.H{
	// 			"message": "pong",
	// 		})
	// })

	// r.POST("/api/v1/parameters", func(c *gin.Context) {
	// 	var parameters models.Parameters
	// 	err := c.BindJSON(&parameters)
	// 	if err != nil {
	// 		c.JSON(
	// 			http.StatusBadRequest,
	// 			gin.H{
	// 				"message": "bad request",
	// 			})
	// 		return
	// 	}
	// 	c.JSON(
	// 		http.StatusOK,
	// 		gin.H{
	// 			"message":  "parameters created",
	// 			"Location": parameters.Location,
	// 			"Days":     parameters.Days,
	// 			"Types":    parameters.Types,
	// 		})
	// })

	// err := r.Run()
	// if err != nil {
	// 	panic(err)
	// }

	sites, err := data.LoadSites()

	if err != nil {
		log.Fatal(err)
	}

	cities := make([]*models.Site, len(sites)) //Postcode    string  `json:"postcode"`

	for i, site := range sites {
		cities[i] = &site
	}

	for _, city1 := range cities {
		for _, city2 := range cities {
			if city1 != city2 {
				services.AddNeighbour(city1, city2)
			}
		}
	}

	services.DFS(cities[2]) // Start DFS from the first site.
}
