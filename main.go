package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli"
)

type weatherResponse struct {
	Weather []weather `json:"weather"`
}

type weather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
}

func main() {
	app := cli.NewApp()
	app.Name = "my-weather"
	app.Usage = "get weather information from console."
	app.Action = func(c *cli.Context) error {
		return handleWeatherAction()
	}

	app.Run(os.Args)
}

func handleWeatherAction() error{
	openWeatherUriTemplate := "http://api.openweathermap.org/data/2.5/weather?q=istanbul&units=metric&appid=e07b609b458eecd57dc0cad0fdb9aa9b"

	weatherClient := http.Client{
		Timeout: time.Second * 10, // maximum 10 seconds to timeout
	}

	req, err := http.NewRequest(http.MethodGet, openWeatherUriTemplate, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := weatherClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("res code " + res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	weatherResponse := weatherResponse{}
	err = json.Unmarshal(body, &weatherResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(weatherResponse.Weather[0].Main + " - " + weatherResponse.Weather[0].Description)
	return nil
}