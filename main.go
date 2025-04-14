package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

const apiKey = "918ea9a31d4edec1143618a5b4457053"

var (
	apiBaseURL = "https://api.openweathermap.org/data/2.5/weather"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
}

func getWeather(city string) (*WeatherResponse, error) {
	url := fmt.Sprintf(
		"%s?q=%s&appid=%s&units=metric",
		apiBaseURL, city, apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get weather: %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var data WeatherResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Weather App</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; }
				.weather-card { background: #f0f8ff; padding: 20px; border-radius: 10px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
				form { margin: 20px 0; }
				input { padding: 8px; width: 70%; }
				button { padding: 8px 15px; background: #4CAF50; color: white; border: none; border-radius: 4px; }
			</style>
		</head>
		<body>
			<h1>Weather App</h1>
			<form method="POST">
				<input type="text" name="city" placeholder="Enter city name" required>
				<button type="submit">Get Weather</button>
			</form>
			{{if .}}
			<div class="weather-card">
				<h2>Weather in {{.Name}}</h2>
				<p><strong>Temperature:</strong> {{.Main.Temp}}°C</p>
				<p><strong>Condition:</strong> {{(index .Weather 0).Description}}</p>
				<p><strong>Humidity:</strong> {{.Main.Humidity}}%</p>
				<img src="https://openweathermap.org/img/wn/{{(index .Weather 0).Icon}}@2x.png" alt="Weather icon">
			</div>
			{{end}}
		</body>
		</html>
	`))

	if r.Method == "GET" {
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		city := r.FormValue("city")
		weather, err := getWeather(city)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, weather)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)
	fmt.Println("Server started on http://localhost:8081")
	http.ListenAndServe(":8081", nil)
}
