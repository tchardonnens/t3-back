package services

import (
	"database/sql"
	"log"
	"os"
	"t3/m/v2/models"

	_ "github.com/mattn/go-sqlite3"
)

func QueryFromDB(parameters models.Parameters) (sites []models.Site) {
	log.Println(os.Getenv("DB_PATH"))
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `SELECT * FROM sites WHERE (city = ? OR region = ? or department = ? or postcode = ?)AND type = ?;`
	//query := `SELECT * FROM sites WHERE spellfix1_city MATCH ? AND type = ?;`

	rows, err := db.Query(query, parameters.Location, parameters.Location, parameters.Location, parameters.Location, parameters.Types)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	sites = make([]models.Site, 0)

	for rows.Next() {
		var (
			id          int
			name        string
			lat         float64
			lng         float64
			_type       string
			postcode    string
			region      string
			department  string
			city        string
			street      string
			website     string
			description string
		)

		err := rows.Scan(&id, &name, &lat, &lng, &_type, &postcode, &region, &department, &city, &street, &website, &description)
		if err != nil {
			log.Fatal(err)
		}

		site := models.Site{
			Id:          id,
			Name:        name,
			Lat:         lat,
			Lng:         lng,
			Type:        _type,
			Postcode:    postcode,
			Region:      region,
			Department:  department,
			City:        city,
			Street:      street,
			Website:     website,
			Description: description,
		}

		sites = append(sites, site)

		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return sites
}
