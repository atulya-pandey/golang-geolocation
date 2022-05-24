package geolocation

import (
	"encoding/json"
	"fmt"
	"golang-geolocation/datastore"
	"golang-geolocation/models"
	"io/ioutil"
	"log"
	"net/http"
)

type GeoLocationHandler struct {
	datastore datastore.GeoLocationInterface
}

func New(geoLocation datastore.GeoLocationInterface) GeoLocationHandler {
	return GeoLocationHandler{datastore: geoLocation}
}

func (a GeoLocationHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.get(w, r)
	case http.MethodPost:
		a.create(w, r)
	case http.MethodDelete:
		a.delete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a GeoLocationHandler) get(w http.ResponseWriter, r *http.Request) {
	ipAddr := r.URL.Query().Get("ipAddr")

	resp, err := a.datastore.Get(ipAddr)
	if err != nil {
		log.Fatal(err)
		_, _ = w.Write([]byte("Could not retrieve geoLocation"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a GeoLocationHandler) create(w http.ResponseWriter, r *http.Request) {
	var geoLocation models.GeoLoc

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &geoLocation)
	if err != nil {
		fmt.Println(err)
		_, _ = w.Write([]byte("Invalid body"))
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resp, err := a.datastore.Create(geoLocation)
	if err != nil {
		_, _ = w.Write([]byte("Could not add geoLocation record"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ = json.Marshal(resp)
	_, _ = w.Write(body)
}

func (a GeoLocationHandler) delete(w http.ResponseWriter, r *http.Request) {
	ipAddr := r.URL.Query().Get("ipAddr")

	resp, err := a.datastore.Delete(ipAddr)
	if err != nil {
		log.Fatal(err)
		_, _ = w.Write([]byte("Could not delete geoLocation record"))
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	body, _ := json.Marshal(resp)
	_, _ = w.Write(body)
}