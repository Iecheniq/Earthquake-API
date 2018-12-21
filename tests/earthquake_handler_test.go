package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"earthquake/controllers"
)

func TestEarthquakeHandler(t *testing.T) {
	testCases := []struct {
		name   string
		method string
		url    string
		code   int
	}{
		{name: "Get a report", method: "GET", url: "localhost:8080/earthquakes/1", code: http.StatusOK},
		{name: "No report found ", method: "GET", url: "localhost:8080/earthquakes/15000", code: http.StatusNotFound},
		{name: "No report found ", method: "DELETE", url: "localhost:8080/earthquakes/15000", code: http.StatusNotFound},
		{name: "No report found ", method: "UPDATE", url: "localhost:8080/earthquakes/15000", code: http.StatusNotFound},
		{name: "Update a report", method: "UPDATE", url: "localhost:8080/earthquakes?magnitude=5&place=Mexico&time=1545414341160", code: http.StatusAccepted},
		{name: "Update a report with missin params", method: "UPDATE", url: "localhost:8080/earthquakes?magnitude=5time=1545414341160", code: http.StatusBadRequest},
		{name: "Delete a report ", method: "DELETE", url: "localhost:8080/earthquakes/5", code: http.StatusOK},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, nil)
			if err != nil {
				t.Fatalf("Could not create request")
			}
			rec := httptest.NewRecorder()
			controllers.EarthquakeHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.code {
				t.Errorf("Expected %v, got %v", tc.code, res.StatusCode)
			}
		})
	}
}
