package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	API_BASE_URL       = "https://api.openweathermap.org/data/2.5/weather?units=imperial&appid="
	OPEN_CAGE_BASE_URL = "https://api.opencagedata.com/geocode/v1/json?"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getWeatherInfo(apiKey, openCageAPIKey, location string) (map[string]string, error) {
	// Geocode the location to get latitude and longitude
	geocodeURL := fmt.Sprintf("%sq=%s&key=%s", OPEN_CAGE_BASE_URL, url.QueryEscape(location), openCageAPIKey)
	resp, err := http.Get(geocodeURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geocodeResult map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&geocodeResult)
	if err != nil {
		return nil, err
	}

	if geocodeResult["results"] == nil || len(geocodeResult["results"].([]interface{})) == 0 {
		return nil, fmt.Errorf("location not found")
	}

	// Get latitude and longitude from the geocode response
	destination := geocodeResult["results"].([]interface{})[0].(map[string]interface{})
	geometry := destination["geometry"].(map[string]interface{})
	latitude := geometry["lat"].(float64)
	longitude := geometry["lng"].(float64)

	// Get the weather information based on latitude and longitude
	fullAPIURL := fmt.Sprintf("%s%s&lat=%s&lon=%s", API_BASE_URL, apiKey, strconv.FormatFloat(latitude, 'f', -1, 64), strconv.FormatFloat(longitude, 'f', -1, 64))
	resp, err = http.Get(fullAPIURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	weatherInfo := map[string]string{
		"address":     destination["formatted"].(string),
		"coordinates": fmt.Sprintf("(%.4f, %.4f)", latitude, longitude),
		"description": result["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
		"temperature": fmt.Sprintf("%.1f\u00B0", result["main"].(map[string]interface{})["temp"].(float64)),
	}

	return weatherInfo, nil
}

func main() {
	// Retrieve API keys from environment variables
	apiKey := os.Getenv("OPEN_WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPEN_WEATHER_API_KEY is not set")
	}

	openCageAPIKey := os.Getenv("OPEN_CAGE_API_KEY")
	if openCageAPIKey == "" {
		log.Fatal("OPEN_CAGE_API_KEY is not set")
	}

	// Initialize Gin
	r := gin.Default()

	// Serve static files from the "static" folder
	r.Static("/static", "./static")

	// Load HTML templates from the "templates" folder
	r.LoadHTMLGlob("templates/*")

	// GET request to render the index page
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// POST request to fetch weather info based on user input location
	r.POST("/", func(c *gin.Context) {
		location := c.PostForm("location")
		weatherInfo, err := getWeatherInfo(apiKey, openCageAPIKey, location)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"weather_info": weatherInfo})
	})

	// Get the port from environment variables or default to "8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	r.Run(":" + port)
}
