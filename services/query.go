package services

import (
	"database/sql"
	"os"
	"t3/m/v2/models"

	_ "github.com/mattn/go-sqlite3"
)

func QueryLocations(queryString string) ([]string, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT DISTINCT city FROM sites WHERE city LIKE ? UNION SELECT DISTINCT region FROM sites WHERE region LIKE ? UNION SELECT DISTINCT department FROM sites WHERE department LIKE ? LIMIT 10;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	pattern := "%" + queryString + "%"
	rows, err := stmt.Query(pattern, pattern, pattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []string

	for rows.Next() {
		var location string

		err := rows.Scan(&location)
		if err != nil {
			return nil, err
		}

		locations = append(locations, location)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locations, nil
}

func CheckIfLocationExistsInDB(location string) (bool, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return false, err
	}
	defer db.Close()

	query := `SELECT EXISTS(SELECT 1 FROM sites WHERE city = ? OR region = ? OR department = ? OR postcode = ?);`

	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var exists bool

	err = stmt.QueryRow(location, location, location, location).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func QuerySitesFromDB(parameters models.Parameters) ([]models.Site, error) {
	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT * FROM sites WHERE (city = ? OR region = ? OR department = ? OR postcode = ?) AND type = ?;`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(parameters.Location, parameters.Location, parameters.Location, parameters.Location, parameters.Types)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []models.Site

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
			return nil, err
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
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}
