package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/infraestructure"
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/usecase" // Add this line to import the usecase package
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/web"
)

func main() {
	// Criar instâncias de serviços e use cases
	apiClient := infraestructure.NewOpenWeatherClient("8f2dfe379acdba84dfa143c5648000e3") // Inicializar seu caso de uso
	// Inicializa o caso de uso, passando o cliente da API
	temperatureUseCase := usecase.NewTemperatureUseCase(apiClient)

	// Inicializa o handler HTTP, passando o caso de uso
	handler := &web.WeatherHandler{
		UseCase: temperatureUseCase,
	}
	// Configurar o roteador chi
	r := chi.NewRouter()

	// Adicionar middlewares (opcional)
	r.Use(middleware.Logger)    // Logger para debugar requisições
	r.Use(middleware.Recoverer) // Recuperar de panics

	// Definir rota para o endpoint
	r.Post("/weather", handler.GetWeather)

	// Iniciar o servidor HTTP
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
