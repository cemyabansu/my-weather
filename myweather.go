package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"

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
	Icon              string  `json:"icon"`
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

	weatherApiUrlTemplate := "https://api.darksky.net/forecast/%s/%f,%f?exclude=hourly,currently&lang=%s&units=%s"
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

	fmt.Println("Here is weekly result.")
	fmt.Printf("Summary: %s\n\n", weatherResponse.Daily.Summary)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"day", "icon", "precip-type", "precip-probability", "tempature", "summary"})
	table.SetRowLine(true)

	for _, dailyData := range weatherResponse.Daily.Data {
		tableItem := []string{"", dailyData.Icon, dailyData.PrecipType, formatPercentage(dailyData.PrecipProbability), formatTempature(dailyData.TemperatureLow), dailyData.Summary}
		table.Append(tableItem)
	}

	table.Render()

	return nil
}

func formatTempature(f float64) string {
	return strconv.FormatFloat(f, 'f', 0, 64) + " °C"
}

func formatPercentage(f float64) string {
	return strconv.FormatFloat(f*100, 'f', 1, 64) + " %"
}
