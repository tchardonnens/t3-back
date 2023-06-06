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
	csvFile, err := os.Open("t3-v8.csv")
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

	db, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))
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

	for _, record := range records {
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
}
