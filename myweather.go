package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
)

type weatherResponse struct {
	TimeZoneOffset	int		`json:"offset"`
	Currently		currently	`json:"currently"`
}

type currently struct {
	Summary			string		`json:"summary"`
	Temperature		float64		`json:"temperature"`
}

func main() {
	app := cli.NewApp()
	app.Name = "my-weather"
	app.Usage = "get weather information from console."

	var places string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:"places, p",
			Value: "istanbul",
			Usage: "places to get weather info",
			Destination: &places,
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println(places)
		return handleWeatherAction(places)
	}

	app.Run(os.Args)
}

func handleWeatherAction(places string) error{
	googleGeoLocationKey := "AIzaSyB_Nib5oe_wIKljESFzNLnPFx2SLa2fCHo"

	c,err := maps.NewClient(maps.WithAPIKey(googleGeoLocationKey))
	if err != nil {
		log.Fatalf("fatal error on creating map client. error : %s", err)
	}
	r := &maps.GeocodingRequest{
		Address: places,
	}
	geoResult, err := c.Geocode(context.Background(),r)
	lat := geoResult[0].Geometry.Location.Lat
	lng := geoResult[0].Geometry.Location.Lng
	fmt.Printf("result : %v, %v \n", lat, lng)

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

	fmt.Println("summary :" + weatherResponse.Currently.Summary)

	return nil
}