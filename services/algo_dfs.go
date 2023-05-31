package services

import (
	"log"
	_ "math"
	"t3/m/v2/data"
	"t3/m/v2/models"
)

func newSite(name string, lat string, lng string, Type string, postcode string, region string, dep string, city string, street string, website string, des string, visted bool, neighbors []*models.Site) *models.Site {
	return &models.Site{
		Name:        name,
		Lat:         lat,
		Lng:         lng,
		Type:        Type,
		Postcode:    postcode,
		Region:      region,
		Department:  dep,
		City:        city,
		Street:      street,
		Website:     website,
		Description: des,
		Visited:     visted,
		Neighbours:  neighbors,
	}
}

func AddNeighbour(site1, site2 *models.Site) {
	site1.Neighbours = append(site1.Neighbours, site2)
	site2.Neighbours = append(site2.Neighbours, site1)
}

func DFS(site *models.Site) {
	log.Println("Visiting " + site.Name)
	site.Visited = true
	log.Println("neighbors ", site.Neighbours)
	for _, neighbour := range site.Neighbours {
		log.Println("====start algo==== ")
		if !neighbour.Visited {
			DFS(neighbour)
		}
	}
}

func implementDFS() {

	sites, err := data.LoadSites()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize cities from loaded sites.
	cities := make([]*models.Site, len(sites))
	for i, city := range cities {
		//(name string, lat, lng float64, Type string, postcode string, region string, dep string, city string, street string, website string, des string, visted bool, neighbors []*models.Site) *models.Site {
		cities[i] = newSite(city.Name, city.Lat, city.Lng, city.Type, city.Postcode, city.Region, city.Department, city.City, city.Street, city.Website, city.Description, city.Visited, city.Neighbours)

	}

	// Add neighbours to cities. We assume every city is connected with every other city.
	for _, city1 := range cities {
		for _, city2 := range cities {
			if city1 != city2 {
				AddNeighbour(city1, city2)
			}
		}
	}

	// Start DFS from the first city.
	DFS(cities[0])

}
