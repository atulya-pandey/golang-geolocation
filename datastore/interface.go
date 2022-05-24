package datastore

import "web-app/models"

type GeoLocationInterface interface {
	LoadData(csvPath string) (int64, error)
	Get(ipAddr string) ([]models.GeoLoc, error)
	Create(geoLoc models.GeoLoc) (map[string]string, error)
	Delete(ipAddr string) (map[string]string, error)
}
