package geolocation

import "github.com/atulya-pandey/golang-geolocation/models"

type GeoLocation interface {
	Get(ipAddr string) ([]models.GeoLocationModel, error)
	Create(geoLoc models.GeoLocationModel) (map[string]string, error)
	Delete(ipAddr string) (map[string]string, error)
}
