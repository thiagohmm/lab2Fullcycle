package entity

import (
	"fmt"
	"strconv"
)

type Temperature struct {
	City      string
	Farenheit float64
	Celsius   float64
	Kelvin    float64
}

func NewTemperature(city string, kelvin float64) (*Temperature, error) {
	return &Temperature{
		City:      city,
		Farenheit: kelvinToFahrenheit(kelvin),
		Celsius:   kelvinToCelsius(kelvin),
		Kelvin:    kelvin,
	}, nil

}

func kelvinToCelsius(kelvin float64) float64 {
	celsius := kelvin - 273.15
	formatCelsius, _ := strconv.ParseFloat(fmt.Sprintf("%.0f", celsius), 64)
	return formatCelsius
}

func kelvinToFahrenheit(kelvin float64) float64 {
	fahrenheit := (kelvin-273.15)*9/5 + 32
	formattedFahrenheit, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", fahrenheit), 64)
	return formattedFahrenheit
}
