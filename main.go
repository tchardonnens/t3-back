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
	config.AllowAllOrigins = false
	config.AllowOrigins = []string{"https://map.verycurious.xyz/"}

	r.Use(cors.New(config))

	models.ConnectDatabase()
	//data.LoadSites()

	r.GET("/api/v1/locations", func(c *gin.Context) {
		queryString := c.DefaultQuery("q", "")
		locations, err := services.QueryLocations(queryString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	})

	r.POST("/api/v1/tsp", func(c *gin.Context) {
		var parameters models.Parameters
		err := c.BindJSON(&parameters)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "Bad request",
				})
			return
		}
		if len(parameters.Types) == 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "No site type",
				})
			return
		}
		if parameters.Days <= 0 {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "Invalid days",
				})
			return
		}

		locationExists, err := services.CheckIfLocationExistsInDB(parameters.Location)
		if parameters.Location == "" || !locationExists || err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"message": "Invalid location",
				})
			return
		}

		siteResults, err := services.QuerySitesFromDB(parameters)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"message": "Internal Server Error",
				})
			return
		}

		pointToSite := make(map[int]models.Site)

		points := make([]models.Point, len(siteResults))
		for i, site := range siteResults {
			point := models.SiteToPoint(site)
			points[i] = point
			pointToSite[point.Id] = site // Save the site for later
		}

		freeTimeHours := services.GetFreeTime(parameters.Days)
		totalNumberOfActivities := services.GetTotalNumberOfActivities(freeTimeHours)
		maxActivityPerDay := services.GetAmountOfActivitiesPerDay(float64(totalNumberOfActivities), parameters.Days)

		tsp := services.TSP(points, maxActivityPerDay, parameters.Days)

		tspSites := make([][]models.Site, len(tsp))
		for i, day := range tsp {
			tspSites[i] = make([]models.Site, len(day))
			for j, point := range day {
				tspSites[i][j] = pointToSite[point.Id]
			}
		}
		c.JSON(
			http.StatusOK,
			gin.H{
				"Location": parameters.Location,
				"Days":     parameters.Days,
				"Types":    parameters.Types,
				"TSP":      tspSites,
			})
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
