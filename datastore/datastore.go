package datastore

import (
	"database/sql"
	"log"
	"github.com/atulya-pandey/golang-geolocation/models"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type GeoLocationStore struct {
	db *sql.DB
}

func New(db *sql.DB) GeoLocationStore {
	return GeoLocationStore{db: db}
}

func (a GeoLocationStore) Get(ipAddr string) ([]models.GeoLocationModel, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if len(ipAddr) > 0 {
		rows, err = a.db.Query("SELECT * FROM geo_loc WHERE ip_address = ?;", string(ipAddr))
	} else {
		rows, err = a.db.Query("SELECT * FROM geo_loc LIMIT 10")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var geoLocations []models.GeoLocationModel

	for rows.Next() {
		var geoLoc models.GeoLocationModel
		_ = rows.Scan(&geoLoc.IpAddress,
			&geoLoc.CountryCode,
			&geoLoc.Country,
			&geoLoc.City,
			&geoLoc.Latitude,
			&geoLoc.Longitude,
			&geoLoc.MysteryValue)
		geoLocations = append(geoLocations, geoLoc)
	}

	return geoLocations, nil
}

func (a GeoLocationStore) Create(geoLoc models.GeoLocationModel) (map[string]string, error) {
	_, err := a.db.Exec("INSERT INTO geo_loc VALUES(?, ?, ?, ?, ?, ?, ?);",
		geoLoc.IpAddress,
		geoLoc.CountryCode,
		geoLoc.Country,
		geoLoc.City,
		geoLoc.Latitude,
		geoLoc.Longitude,
		geoLoc.MysteryValue,
	)

	if err != nil {
		x := map[string]string{"createStatus": "FAIL"}
		return x, err
	}

	return map[string]string{"createStatus": "SUCCESS"}, nil
}

func (a GeoLocationStore) Delete(ipAddr string) (map[string]string, error) {
	_, err := a.db.Exec("DELETE FROM geo_loc WHERE ip_address=?;", string(ipAddr))

	if err != nil {
		x := map[string]string{"deleteStatus": "FAIL"}
		return x, err
	}

	return map[string]string{"deleteStatus": "SUCCESS"}, nil
}

func (a GeoLocationStore) LoadData(csvPath string) (int64, error) {
	mysql.RegisterLocalFile(csvPath)

	log.Println("Truncating table before loading data")

	_, err := a.db.Exec("TRUNCATE geo_loc;")

	log.Printf("Truncate table successful")

	log.Printf("Loading data from %s...", csvPath)

	result, err := a.db.Exec("LOAD DATA LOCAL INFILE '" + csvPath + "' INTO TABLE geo_loc FIELDS TERMINATED BY ',' IGNORE 1 ROWS")

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	log.Println("Data loaded successfully.")

	return rowsAffected, nil
}
