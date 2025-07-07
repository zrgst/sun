package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#34567c")).
	Background(lipgloss.Color("#2e2a31")).
	AlignHorizontal(lipgloss.Center).
	Width(50)

var weatherLocationStyle = lipgloss.NewStyle().
	Bold(false).
	Foreground(lipgloss.Color("#080715")).
	Background(lipgloss.Color("#34567c")).
	AlignHorizontal(lipgloss.Center).
	AlignVertical(lipgloss.Center).
	Width(50)

var rainyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#8b1a1e")).
	AlignHorizontal(lipgloss.Left).
	Width(60)

var notRainyStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#d5e9ea")).
	Bold(true)

var infoStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#34567c")).
	Italic(true)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"current`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

// API Key: d454af69a7df4f6f9ca113352250507

func main() {
	q := "Stavanger" // Replace with your desired location

	if len(os.Args) >= 2 {
		q = os.Args[1] // Use the first command line argument as the query
	}

	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=d454af69a7df4f6f9ca113352250507&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not aviliable")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	fmt.Println(headerStyle.Render("Zrgst@sun.app"))

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	locationMessage := fmt.Sprintf(
		"%s, %s: %.0fC, %s",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	fmt.Print(weatherLocationStyle.Render(locationMessage))
	fmt.Println()
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if date.Before(time.Now()) {
			continue // Skip past hours
		}

		message := fmt.Sprintf(
			"%s - %.0fC, %.0f%%, %s",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)

		if hour.ChanceOfRain < 40 {
			fmt.Print(notRainyStyle.Render(message))
			fmt.Println()
		} else {
			fmt.Print(rainyStyle.Render(message)) // Print in red if chance of rain is 40% or more
			fmt.Println()
		}

	}
	fmt.Println(infoStyle.Render("- For custom location add location\n after command:\n$ sun <location>"))
}
