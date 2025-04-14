package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Mock server for testing getWeather
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		response := `{
			"name": "London",
			"main": {
				"temp": 15.5,
				"humidity": 75
			},
			"weather": [{
				"description": "clear sky",
				"icon": "01d"
			}]
		}`
		w.Write([]byte(response))
	}))
}

func TestGetWeather(t *testing.T) {
	// Setup mock server
	server := setupMockServer()
	defer server.Close()

	originalURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalURL }()

	// Test with a valid city
	weather, err := getWeather("London")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if weather.Name != "London" {
		t.Errorf("Expected city name London, got %s", weather.Name)
	}

	if weather.Main.Temp != 15.5 {
		t.Errorf("Expected temp 15.5, got %f", weather.Main.Temp)
	}

	if weather.Weather[0].Description != "clear sky" {
		t.Errorf("Expected description 'clear sky', got '%s'", weather.Weather[0].Description)
	}
}

func TestHomeHandlerGET(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the form is present
	if !strings.Contains(rr.Body.String(), "<form method=\"POST\">") {
		t.Error("handler response does not contain the form")
	}
}

func TestHomeHandlerPOST(t *testing.T) {
	// Setup form data
	form := url.Values{}
	form.Add("city", "London")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Setup mock server for the API call
	server := setupMockServer()
	defer server.Close()
	originalURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalURL }()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Weather in London") {
		t.Error("handler response does not contain city name")
	}

	if !strings.Contains(body, "15.5°C") {
		t.Error("handler response does not contain expected temperature")
	}

	if !strings.Contains(body, "clear sky") {
		t.Error("handler response does not contain expected weather description")
	}
}

func TestHomeHandlerPOSTError(t *testing.T) {
	// Setup form data
	form := url.Values{}
	form.Add("city", "InvalidCity")

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Setup mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	originalURL := apiBaseURL
	apiBaseURL = server.URL
	defer func() { apiBaseURL = originalURL }()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(homeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}