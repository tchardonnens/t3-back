package services

import (
	"fmt"
	"math"
	"t3/m/v2/models"

	"github.com/dominikbraun/graph"
)

// calculate the distance using lat and lng
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371 // radius of the earth in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return r * c // returns the distance in km
}

// Extract the preparation of graph and sitemap as a separate function
func PrepareGraphAndSiteMap(sites []*models.Site) (*graph.Graph, map[string]*models.Site) {
	g := graph.New(models.SiteHash, graph.Weighted())
	siteMap := make(map[string]*models.Site)
	for i := 0; i < len(sites)-1; i++ {
		_ = g.AddVertex(sites[i])
		siteMap[models.SiteHash(sites[i])] = sites[i]
		distance := Haversine(sites[i].Lat, sites[i].Lng, sites[i+1].Lat, sites[i+1].Lng)
		_ = g.AddVertex(sites[i+1])
		siteMap[models.SiteHash(sites[i+1])] = sites[i+1]
		distanceInt := int(distance)
		_ = g.AddEdge(models.SiteHash(sites[i]), models.SiteHash(sites[i+1]), graph.EdgeWeight(distanceInt))
	}
	return g, siteMap
}

// limit is the result numbers of dfs
func ImplementDFS(limit int) map[string]*models.Site {

	// g := graph.New(models.SiteHash, graph.Weighted())

	// for i := 0; i < len(sites)-1; i++ {
	// 	_ = g.AddVertex(&sites[i])
	// 	distance := data.Haversine(sites[i].Lat, sites[i].Lng, sites[i+1].Lat, sites[i+1].Lng)
	// 	_ = g.AddVertex(&sites[i+1])
	// 	distanceInt := int(distance)
	// 	_ = g.AddEdge(models.SiteHash(&sites[i]), models.SiteHash(&sites[i+1]), graph.EdgeWeight(distanceInt))
	// }
	// log.Println("-----graph-----", g)

	// _ = graph.DFS(g, models.SiteHash(&sites[0]), func(site *models.Site) bool {
	// 	fmt.Println(site) // print the site or any specific field you are interested in
	// 	return false
	// })
	// _ = graph.BFS(g, 1, false)
	sites := models.GetAllSites()
	g, siteMap := PrepareGraphAndSiteMap(sites)

	count := 0
	_ = graph.DFS(g, models.SiteHash(&sites[0]), func(hash string) bool {
		site := siteMap[hash]
		fmt.Println(site.Name)
		count++
		if count >= limit {
			return true
		}
		return false
	})

	return siteMap
}
