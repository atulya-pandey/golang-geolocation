package main

import (
	"fmt"
	"golang-geolocation/datastore/geolocation"
	handlerGeoLocation "golang-geolocation/delivery/geolocation"
	"golang-geolocation/driver"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	conf := driver.MySQLConfig{
		Host:     os.Getenv("HOST"),
		User:     os.Getenv("DBUSER"),
		Password: os.Getenv("DBPASS"),
		Port:     os.Getenv("PORT"),
		Db:       os.Getenv("DB"),
		Net:      "tcp",
	}

	var err error

	db, err := driver.ConnectToMySQL(conf)

	if err != nil {
		log.Println("Could not connect to sql, err:", err)
		return
	}

	datastore := geolocation.New(db)

	datastore.LoadData("data/data_dump.csv")
	handler := handlerGeoLocation.New(datastore)

	http.HandleFunc("/geolocation", handler.Handler)
	fmt.Println(http.ListenAndServe(":9000", nil))
}
