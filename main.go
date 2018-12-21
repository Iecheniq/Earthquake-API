package main

import (
	"log"
	"net/http"

	"earthquake/controllers"

	_ "github.com/go-sql-driver/mysql"
	mux "github.com/gorilla/mux"
)

func main() {
	setRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", controllers.HomeHandler).
		Methods("GET").
		Name("home")
	r.HandleFunc("/earthquakes", controllers.EarthquakesHandler).
		Methods("GET", "POST").
		Name("earthquakes")
	r.HandleFunc("/earthquakes/{Id:[0-9]+}", controllers.EarthquakeHandler).
		Methods("GET", "DELETE", "PUT").
		Name("erthquake")

	http.Handle("/", r)
}
