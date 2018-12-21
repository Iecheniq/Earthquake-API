package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"earthquake/controllers"
)

func TestEathquakesHandler(t *testing.T) {
	type TestBody struct {
		Magnitude float32
		Place     string
		Time      int64
	}
	testCases := []struct {
		name   string
		method string
		Body   TestBody
		code   int
	}{
		{name: "List all reports",
			method: "GET",
			Body:   TestBody{},
			code:   http.StatusOK},
		{name: "Create new report",
			method: "POST",
			Body: TestBody{
				Magnitude: 8,
				Place:     "Mexico",
				Time:      1112920,
			},
			code: http.StatusAccepted},
		{name: "Create new report with missin params",
			method: "POST",
			Body: TestBody{
				Magnitude: 8,
				Place:     "Mexico",
			},
			code: http.StatusBadRequest},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(tc.Body)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest(tc.method, "localhost:8080/earthquakes", bytes.NewBuffer(payload))
			if err != nil {
				t.Fatalf("Could not create request")
			}
			rec := httptest.NewRecorder()
			controllers.EarthquakesHandler(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tc.code {
				t.Errorf("Expected %v, got %v", tc.code, res.StatusCode)
			}
		})
	}
}
