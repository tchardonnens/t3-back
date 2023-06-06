package models

type Point struct {
	Id      int
	Name    string
	X       float64
	Y       float64
	Visited bool
}

func SiteToPoint(site Site) Point {
	return Point{
		Id:      site.Id,
		Name:    site.Name,
		X:       site.Lng, // Usually, longitude is considered as 'x'
		Y:       site.Lat, // and latitude is considered as 'y'
		Visited: false,
	}
}
