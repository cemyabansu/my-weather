package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type weatherResponse struct {
	TimeZoneOffset int       `json:"offset"`
	Currently      currently `json:"currently"`
	Daily          daily     `json:"daily"`
}

type currently struct {
	Summary     string  `json:"summary"`
	Temperature float64 `json:"temperature"`
}

type daily struct {
	Summary string      `json:"summary"`
	Data    []dailyData `json:"data"`
}

type dailyData struct {
	Time              int     `json:"time"`
	Summary           string  `json:"summary"`
	TemperatureHigh   float64 `json:"temperatureHigh"`
	TemperatureLow    float64 `json:"temperatureLow"`
	PrecipType        string  `json:"precipType"`
	PrecipProbability float64 `json:"precipProbability"`
}

func main() {
	app := cli.NewApp()
	app.Name = "my-weather"
	app.Usage = "get weather information from console."

	var place string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "place, p",
			Value:       "istanbul",
			Usage:       "place to get weather info",
			Destination: &place,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(place)
		return handleWeatherAction(place)
	}

	app.Run(os.Args)
}

func handleWeatherAction(places string) error {
	googleGeoLocationKey := "AIzaSyB_Nib5oe_wIKljESFzNLnPFx2SLa2fCHo"

	c, err := maps.NewClient(maps.WithAPIKey(googleGeoLocationKey))
	if err != nil {
		log.Fatalf("fatal error on creating map client. error : %s", err)
	}
	r := &maps.GeocodingRequest{
		Address: places,
	}
	geoResult, err := c.Geocode(context.Background(), r)
	lat := geoResult[0].Geometry.Location.Lat
	lng := geoResult[0].Geometry.Location.Lng
	fmt.Printf("result : %s --> %v, %v \n", geoResult[0].FormattedAddress, lat, lng)
	fmt.Printf("Here is the weather for %s \n", geoResult[0].FormattedAddress)

	weatherApiUrlTemplate := "https://api.darksky.net/forecast/%s/%f,%f?exclude=hourly&lang=%s&units=%s"
	apiKey := "8783cc8f39ddbbd37ecc97ea0c958e0d"
	lang := "tr"
	units := "si"

	weatherApiUrl := fmt.Sprintf(weatherApiUrlTemplate, apiKey, lat, lng, lang, units)
	fmt.Println(weatherApiUrl)

	response, err := http.Get(weatherApiUrl)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return err
	}

	data, _ := ioutil.ReadAll(response.Body)

	weatherResponse := weatherResponse{}
	err = json.Unmarshal(data, &weatherResponse)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Currently \nSummary\t\t: %s\nTemperature\t: %f°C\n",
		weatherResponse.Currently.Summary,
		weatherResponse.Currently.Temperature)

	fmt.Println("Here is weekly result.")
	fmt.Printf("Summary: %s\n\n", weatherResponse.Daily.Summary)

	counter := 0
	for _, dailyData := range weatherResponse.Daily.Data {
		fmt.Println(counter)

		fmt.Printf("Summary: %s\nTemperature Low \t: %f°C\nTemperature High \t: %f°C\nPrecipType \t\t: %s\nPrecipProbability \t: %f\n",
			dailyData.Summary,
			dailyData.TemperatureLow,
			dailyData.TemperatureHigh,
			dailyData.PrecipType,
			dailyData.PrecipProbability*100)
		fmt.Println("==================")

		counter++
	}

	return nil
}
