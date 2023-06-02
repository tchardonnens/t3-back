package data

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"t3/m/v2/models"

	_ "github.com/mattn/go-sqlite3"
)

func correctFormat(val string) float64 {

	if val == "" {
		return 0.0000 // Or return a default value, e.g., "0"
	}
	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Println("Error converting string to float:", err)
		return floatVal
	}

	// Adjusting for latitudes >90 and longitudes >180
	if floatVal > 90.0 || floatVal > 180.0 {
		floatVal = floatVal / 1000000
	}

	return floatVal
}

func LoadSites() ([]models.Site, error) {
	log.Println("====== Creating sites database... ======")
	csvFile, err := os.Open("t3-v7.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	_, _ = reader.Read()
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(records) == 0 {
		log.Fatal("No records found at: " + csvFile.Name())
	}

	db, err := sql.Open("sqlite3", "./t3.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	query := `CREATE TABLE IF NOT EXISTS sites (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		lat FLOAT NOT NULL,
		lng FLOAT NOT NULL,
		street VARCHAR(255) NOT NULL,
		city VARCHAR(255) NOT NULL,
		postcode VARCHAR(255) NOT NULL,
		department VARCHAR(255) NOT NULL,
		region VARCHAR(255) NOT NULL,
		website VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Sites table created")
	log.Default().Println("Loading sites...")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO sites (name, type, lat, lng, street, city, postcode, department, region, website, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	var sites []models.Site
	log.Default().Println("===== Map model.Sites in the local variable =====")
	for _, record := range records {
		correctedLat := correctFormat(record[2])
		correctedLng := correctFormat(record[3])

		site := models.Site{
			Name:        record[0],
			Type:        record[1],
			Lat:         correctedLat,
			Lng:         correctedLng,
			Street:      record[4],
			City:        record[5],
			Postcode:    record[6],
			Department:  record[7],
			Region:      record[8],
			Website:     record[9],
			Description: record[10],
			Visited:     false,
			Neighbours:  make([]*models.Site, 0),
		}

		sites = append(sites, site)

		_, err = stmt.Exec(record[0], record[1], correctedLat, correctedLng, record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		if err != nil {
			log.Println("Error inserting record:", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				log.Println("Rollback error:", rollbackErr)
			}
			log.Fatal("Transaction rolled back due to error")
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Sites loaded")
	log.Println("====== Sites database created ======")
	return sites, nil
}
