package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var db map[int]EarthquakeData
var idCounter int

type EarthquakeService interface {
	GetData() error
	ParseData() []EarthquakeData
}

type USGSService struct {
	Data []struct {
		Properties struct {
			Magnitude float64 `json:"mag"`
			Place     string  `json:"place"`
			Time      int64   `json:"time"`
		} `json:"properties"`
	} `json:"features"`
}

type EarthquakeData struct {
	Status    int
	Id        int
	Magnitude float64
	Place     string
	Time      time.Time
}

func init() {
	if err := createDb(); err != nil {
		log.Fatal(err)
	}
}

func createDb() error {
	db = make(map[int]EarthquakeData)
	idCounter = 1
	service := &USGSService{}
	if err := service.GetData(); err != nil {
		return err
	}
	earthquakes := service.ParseData()
	for _, e := range earthquakes {
		e.Id = idCounter
		db[idCounter] = e
		idCounter += 1
	}
	return nil
}

func (s *USGSService) GetData() error {
	url := "https://earthquake.usgs.gov/fdsnws/event/1/query?format=geojson"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(body, s); err != nil {
		return err
	}
	return nil
}

func (s USGSService) ParseData() []EarthquakeData {
	earthquakes := make([]EarthquakeData, 0)
	for _, d := range s.Data {
		e := EarthquakeData{}
		e.Magnitude = d.Properties.Magnitude
		e.Place = d.Properties.Place
		e.Time = time.Unix(d.Properties.Time/1000, 0)
		e.Status = 200
		earthquakes = append(earthquakes, e)
	}
	return earthquakes
}

func GetEarthquakes() []EarthquakeData {
	earthquakes := make([]EarthquakeData, 0)
	const limit float64 = 48
	for k, v := range db {
		if time.Since(db[k].Time).Hours() <= limit {
			earthquakes = append(earthquakes, v)
		}
	}
	return earthquakes
}
func GetEarthquake(id int) (*EarthquakeData, error) {
	if _, ok := db[id]; !ok {
		return nil, fmt.Errorf("No event found with ID %v", id)
	}
	e := db[id]
	return &e, nil

}
func (e *EarthquakeData) AddEarthquake() {
	e.Id = idCounter
	e.Status = http.StatusCreated
	db[idCounter] = *e
	idCounter += 1
}

func DeleteEarthquake(id int) error {
	if _, ok := db[id]; !ok {
		return fmt.Errorf("No event found with ID %v", id)
	}
	delete(db, id)
	return nil
}

func (e *EarthquakeData) UpdateEarthquake() error {
	if _, ok := db[e.Id]; !ok {
		return fmt.Errorf("No event found with ID %v", e.Id)
	}
	e.Status = http.StatusAccepted
	db[e.Id] = *e
	return nil
}
