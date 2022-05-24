package geolocation

import (
	"database/sql"
	"log"
	"web-app/models"
	"github.com/go-sql-driver/mysql"
)

type geoLocationStore struct {
	db *sql.DB
}

func New(db *sql.DB) geoLocationStore {
	return geoLocationStore{db: db}
}

func (a geoLocationStore) Get(ipAddr string) ([]models.GeoLoc, error) {
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

	var geoLocations []models.GeoLoc

	for rows.Next() {
		var geoLoc models.GeoLoc
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

func (a geoLocationStore) Create(geoLoc models.GeoLoc) (map[string]string, error) {
	res, err := a.db.Exec("INSERT INTO geo_loc VALUES(?, ?, ?, ?, ?, ?, ?);", 
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

	rowsAffected, _ := res.RowsAffected()

	if rowsAffected >= 0 {
	}

	return map[string]string{"createStatus": "SUCCESS"}, nil
}

func (a geoLocationStore) Delete(ipAddr string) (map[string]string, error) {
	res, err := a.db.Exec("DELETE FROM geo_loc WHERE ip_address=?;", string(ipAddr))

	if err != nil {
		x := map[string]string{"deleteStatus": "FAIL"}
		return x, err
	}

	_, _ = res.RowsAffected()

	return map[string]string{"deleteStatus": "SUCCESS"}, nil
}

func (a geoLocationStore) LoadData(csvPath string) (int64, error) {
	mysql.RegisterLocalFile(csvPath)

	log.Println("Truncating table before loading data")

	truncateResult, err := a.db.Exec("TRUNCATE geo_loc;")
	truncateAffectedRows, err := truncateResult.RowsAffected()

	log.Printf("Truncate table. Removed %d rows", truncateAffectedRows)

	log.Printf("Loading data from %s...", csvPath)

	result, err := a.db.Exec("LOAD DATA LOCAL INFILE '" + csvPath + "' INTO TABLE geo_loc FIELDS TERMINATED BY ',' IGNORE 1 ROWS")

	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	log.Println("Data loaded successfully.", rowsAffected)

	return rowsAffected, nil
}