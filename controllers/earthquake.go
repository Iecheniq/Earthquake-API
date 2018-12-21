package controllers

import (
	"earthquake/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	mux "github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(fmt.Sprintf("Welcome to the Earthquakes API"))); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func EarthquakesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		earthquakes := models.GetEarthquakes()
		earthquakesJ, err := json.Marshal(earthquakes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(earthquakesJ); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			fmt.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bla, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s", bla)

		m := r.Form.Get("magnitude")
		p := r.Form.Get("place")
		t := r.Form.Get("time")

		if m == "" || p == "" || t == "" {
			http.Error(w, "One or more fileds are empty, you must enter magnitude, place, time", http.StatusBadRequest)
			return
		}
		magnitude, err := strconv.ParseFloat(m, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		timedate, err := strconv.Atoi(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		earthquake := models.EarthquakeData{Magnitude: magnitude,
			Place: p,
			Time:  time.Unix(int64(timedate), 0),
		}
		earthquake.AddEarthquake()
		earthquakeJ, err := json.Marshal(earthquake)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		if _, err := w.Write(earthquakeJ); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func EarthquakeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id, err := strconv.Atoi(vars["Id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "GET" {
		earthquake, err := models.GetEarthquake(Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		earthquakeJ, err := json.Marshal(earthquake)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(earthquakeJ); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		if err := models.DeleteEarthquake(Id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if _, err := w.Write([]byte(fmt.Sprintf("Event with ID %v deleted", Id))); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "PUT" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.Form.Get("magnitude")
		p := r.Form.Get("place")
		t := r.Form.Get("time")
		if m == "" || p == "" || t == "" {
			http.Error(w, "One or more fileds are empty, you must enter magnitude, place, time", http.StatusBadRequest)
			return
		}
		magnitude, err := strconv.ParseFloat(m, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
		timedate, err := strconv.Atoi(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		earthquake := models.EarthquakeData{
			Id:        Id,
			Magnitude: magnitude,
			Place:     p,
			Time:      time.Unix(int64(timedate), 0),
		}
		if err := earthquake.UpdateEarthquake(); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		earthquakeJ, err := json.Marshal(earthquake)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		if _, err := w.Write(earthquakeJ); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

}
