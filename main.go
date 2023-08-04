package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type City struct {
	Result [1]struct {
		Name string `json:"name"`
		Latitude float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
		Country string `json:"country"`
	} `json:"results"`
}

type Forecast struct {
	CurrentWeather struct {
		Temperature float32 `json:"temperature"`
		WindSpeed float32 `json:"windspeed"`
	} `json:"current_weather"`
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Insert a city name: ")
	city, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("Error while getting the input")
		reader.ReadLine()
		return
	}
	// fmt.Println(len(city))
	city = city[:len(city) - 2] // To remove the \n and the blank space
	
	c := http.Client{}

	greq, _ := http.NewRequest("GET", fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(city)), nil)
	
	gresp, err := c.Do(greq)
	
	if err != nil || gresp.StatusCode != 200 {
		fmt.Println("City not found")
		reader.ReadLine()
		return
	}

	defer gresp.Body.Close()
	
	gbytes, err := io.ReadAll(gresp.Body)
	
	if err != nil {
		fmt.Println("Cant read the response")
		reader.ReadLine()
		return
    }
	
	var ct City
	
	json.Unmarshal(gbytes, &ct)


	if ct.Result[0].Country == "" {
		fmt.Println("City not found")
		reader.ReadLine()
		return
	}
	
	wreq, _ := http.NewRequest("GET", fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true&timezone=auto", ct.Result[0].Latitude, ct.Result[0].Longitude), nil)

	wresp, err := c.Do(wreq)
	
	if err != nil || wresp.StatusCode != 200 {
		fmt.Println("Error while retrieving weather data")
		reader.ReadLine()
		return
	}
	defer wresp.Body.Close()
	
	wbytes, err := io.ReadAll(wresp.Body)
	
	if err != nil {
		fmt.Println("Cant read the response")
		reader.ReadLine()
		return
    }

	var f Forecast

	json.Unmarshal(wbytes, &f)

	fmt.Println(fmt.Sprintf("In %s (%s) there are %.1f°C with a wind speed of %.2f km/h", ct.Result[0].Name, ct.Result[0].Country, f.CurrentWeather.Temperature, f.CurrentWeather.WindSpeed))
	reader.ReadLine()

}