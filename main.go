package main

import (
	"io/ioutil"
	"log"
	"os"
	"github.com/atulya-pandey/golang-geolocation/datastore"
	"github.com/atulya-pandey/golang-geolocation/geolocation"
	"github.com/joho/godotenv"
)

func main() {
	db, err := datastore.ConnectToMySQL()

	if err != nil {
		log.Println("Could not connect to sql, err:", err)
		return
	}
	
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// Create table
	log.Printf("Creating Table %s", os.Getenv("TABLENAME"))

	createSqlScriptPath := "sql/geolocation/create_geolocation.sql"

	c, ioErr := ioutil.ReadFile(createSqlScriptPath)
	if ioErr != nil {
		log.Fatalf("Failed to load %s", createSqlScriptPath)
	}
	createTableSql := string(c)
	db.Exec(createTableSql)

	datastore := datastore.New(db)
	importer := geolocation.New(datastore)

	parsedCsvPath, _ := importer.Parse("data/data_dump.csv")

	importStatus, err := importer.Import(parsedCsvPath)
	log.Printf("Import Status: %s", importStatus)

	importer.Statistics()
	
}
