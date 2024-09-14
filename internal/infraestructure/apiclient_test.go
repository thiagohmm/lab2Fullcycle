package infraestructure

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock das respostas da API ViaCEP
var mockViaCepResponse = ViaCepResponse{
	Cep:         "12345-678",
	Logradouro:  "Rua Teste",
	Complemento: "Apto 1",
	Bairro:      "Vila Prudente",
	Localidade:  "São Paulo",
	Uf:          "SP",
}

// Mock das respostas da API OpenWeather Geo
var mockGeoResponse = []GeoResponse{
	{
		Name:    "Vila Prudente",
		Lat:     -23.5898,
		Lon:     -46.5646,
		Country: "BR",
		State:   "SP",
	},
}

// Mock das respostas da API OpenWeather Weather
var mockWeatherResponse = WeatherResponse{
	Main: struct {
		Temp float64 `json:"temp"`
	}{
		Temp: 301.9,
	},
}

func TestGetTemperatureByCep(t *testing.T) {
	// Mock server for ViaCEP API
	viaCepServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockViaCepResponse)
	}))
	defer viaCepServer.Close()

	// Mock server for OpenWeatherMap Geo API
	geoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockGeoResponse)
	}))
	defer geoServer.Close()

	// Mock server for OpenWeatherMap Weather API
	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockWeatherResponse)
	}))
	defer weatherServer.Close()

	// Override the URLs in the OpenWeatherClient
	client := &OpenWeatherClient{
		apiKey: "teste-api-key",
		getCepDataFunc: func(ctx context.Context, cep string) (*ViaCepResponse, error) {
			resp, err := http.Get(viaCepServer.URL)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			var cepData ViaCepResponse
			if err := json.NewDecoder(resp.Body).Decode(&cepData); err != nil {
				return nil, err
			}

			return &cepData, nil
		},
		getGeoDataFunc: func(ctx context.Context, bairro, localidade string) (*GeoResponse, error) {
			resp, err := http.Get(geoServer.URL)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			var geoData []GeoResponse
			if err := json.NewDecoder(resp.Body).Decode(&geoData); err != nil {
				return nil, err
			}

			if len(geoData) == 0 {
				return nil, fmt.Errorf("nenhum dado de geolocalização encontrado")
			}

			return &geoData[0], nil
		},
		getWeatherDataFunc: func(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
			resp, err := http.Get(weatherServer.URL)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			var weatherData WeatherResponse
			if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
				return nil, err
			}

			return &weatherData, nil
		},
	}

	ctx := context.Background()
	temp, _, err := client.GetTemperatureByCep(ctx, "12345-678")
	if err != nil {
		t.Fatalf("GetTemperatureByCep() error = %v", err)
	}

	expectedTemp := 301.9
	if temp != expectedTemp {
		t.Errorf("GetTemperatureByCep() = %v, expected %v", temp, expectedTemp)
	}
}
