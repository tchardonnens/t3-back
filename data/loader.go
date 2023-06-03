package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"t3/m/v2/models"

	_ "github.com/mattn/go-sqlite3"
)

func correctFormat(val string) float64 {
	if val == "" {
		return 0.0 // Or return a default value, e.g., "0"
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Println("Error converting string to float:", err)
		return floatVal
	}

	// Adjusting for latitudes >90 and longitudes >180
	// Note: The division number (1000000, 1000000000 etc) depends on how your data is structured.
	// You might need to adjust these values to correctly process your data.
	if floatVal > 90000000 { // Greater than 90 million, likely in degree*1e6 format
		floatVal = floatVal / 1000000
	} else if floatVal > 900000000000 { // Greater than 90 billion, likely in degree*1e12 format
		floatVal = floatVal / 1000000000000
	} else if floatVal > 18000000 && floatVal < 90000000 { // Greater than 18 million and less than 90 million, likely in degree*1e6 format
		floatVal = floatVal / 1000000
	} else if floatVal > 180000000000 && floatVal < 900000000000 { // Greater than 18 billion and less than 90 billion, likely in degree*1e12 format
		floatVal = floatVal / 1000000000000
	}

	return floatVal
}

func createSite(record []string) (site models.Site, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error occurred while creating site: %v", r)
		}
	}()

	correctedLat := correctFormat(record[2])
	correctedLng := correctFormat(record[3])

	site = models.Site{
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

	return site, nil
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

		site, err := createSite(record)
		if err != nil {
			log.Printf("error model site %v", err)
			continue
		}

		sites = append(sites, site)

		_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
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
