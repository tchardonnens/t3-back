package data

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func LoadSites() {
	log.Println("====== Creating sites database... ======")
	csvFile, _ := os.Open("t3-v7.csv")
	reader := csv.NewReader(csvFile)
	reader.Comma = ';'
	_, _ = reader.Read()
	records, _ := reader.ReadAll()
	if records == nil {
		log.Fatal("No records found at: " + csvFile.Name())
	}

	db, _ := sql.Open("sqlite3", "./t3.db")
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

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("Sites table created")
	log.Default().Println("Loading sites...")
	tx, _ := db.Begin()

	stmt, _ := tx.Prepare("INSERT INTO sites (name, type, lat, lng, street, city, postcode, department, region, website, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	for _, record := range records {
		_, err = stmt.Exec(record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], record[8], record[9], record[10])
		if err != nil {
			log.Fatal(err)
		}
	}

	tx.Commit()
	log.Println("Sites loaded")
	log.Println("====== Sites database created ======")
}
