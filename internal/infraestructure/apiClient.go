package infraestructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Interface para buscar a temperatura por CEP
type GetTemperatureForCep interface {
	GetTemperatureByCep(ctx context.Context, cep string) (float64, string, error)
}

// Estrutura que implementa a interface
type OpenWeatherClient struct {
	apiKey             string
	getCepDataFunc     func(ctx context.Context, cep string) (*ViaCepResponse, error)
	getGeoDataFunc     func(ctx context.Context, bairro, localidade string) (*GeoResponse, error)
	getWeatherDataFunc func(ctx context.Context, lat, lon float64) (*WeatherResponse, error)
}

// Construtor para o cliente OpenWeather
func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	client := &OpenWeatherClient{
		apiKey: apiKey,
	}
	client.getCepDataFunc = client.getCepData
	client.getGeoDataFunc = client.getGeoData
	client.getWeatherDataFunc = client.getWeatherData
	return client
}

// Resposta da API ViaCEP
type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

// Resposta da API de geolocalização
type GeoResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

// Resposta da API de clima
type WeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

// Método para obter os dados do CEP
func (c *OpenWeatherClient) getCepData(ctx context.Context, cep string) (*ViaCepResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cepData ViaCepResponse
	if err := json.Unmarshal(body, &cepData); err != nil {
		fmt.Println("Erro ao decodificar resposta da API ViaCEP:", string(body))
		return nil, err
	}

	return &cepData, nil
}

// Método para obter os dados de geolocalização
func (c *OpenWeatherClient) getGeoData(ctx context.Context, bairro, localidade string) (*GeoResponse, error) {
	bairro = url.QueryEscape(bairro)
	localidade = url.QueryEscape(localidade)
	fmt.Println(localidade, bairro)
	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,%s&appid=%s", bairro, localidade, c.apiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geoData []GeoResponse
	if err := json.Unmarshal(body, &geoData); err != nil {
		fmt.Println("Erro ao decodificar resposta da API OpenWeatherMap (Geo):", string(body))
		return nil, err
	}

	if len(geoData) == 0 {
		return nil, fmt.Errorf("nenhum dado de geolocalização encontrado")
	}

	return &geoData[0], nil
}

// Método para obter os dados de clima
func (c *OpenWeatherClient) getWeatherData(ctx context.Context, lat, lon float64) (*WeatherResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s", lat, lon, c.apiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		fmt.Println("Erro ao decodificar resposta da API OpenWeatherMap (Weather):", string(body))
		return nil, err
	}

	return &weatherData, nil
}

// Método principal da interface para obter a temperatura por CEP
func (c *OpenWeatherClient) GetTemperatureByCep(ctx context.Context, cep string) (float64, string, error) {
	cepData, err := c.getCepDataFunc(ctx, cep)
	if err != nil {
		return 0, "", err
	}

	geoData, err := c.getGeoDataFunc(ctx, cepData.Bairro, cepData.Localidade)
	if err != nil {
		return 0, "", err
	}

	weatherData, err := c.getWeatherDataFunc(ctx, geoData.Lat, geoData.Lon)
	if err != nil {
		return 0, "", err
	}

	return weatherData.Main.Temp, cepData.Localidade, nil
}
