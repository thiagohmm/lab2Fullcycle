package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/usecase"
	"go.opentelemetry.io/otel"
)

type WeatherHandler struct {
	UseCase *usecase.TemperatureUseCase
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	var req usecase.TemperatureInputDTO

	// Obter o tracer global
	tracer := otel.Tracer("weather-handler")

	// Iniciar um span para a requisição HTTP
	ctx, span := tracer.Start(r.Context(), "GetWeatherHandler")
	defer span.End()

	// Logando o corpo da requisição
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Request Body: ", string(bodyBytes))

	// Reconstituir o corpo para o decoder
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Reutiliza o corpo da requisição

	// Decodificar o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	cep := req.Cep

	// Validação do CEP: deve conter exatamente 8 dígitos numéricos
	validCep := regexp.MustCompile(`^\d{8}$`)
	if !validCep.MatchString(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// Iniciar um novo span para a execução do caso de uso
	ctx, spanUseCase := tracer.Start(ctx, "ExecuteUseCase")
	defer spanUseCase.End()

	// Chamar o caso de uso para obter o clima
	dto := usecase.TemperatureInputDTO{Cep: cep}
	weather, err := h.UseCase.Execute(ctx, dto) // Passar o contexto traçado para o caso de uso
	if err != nil {
		if err.Error() == "CEP not found" {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get weather data: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Definir o cabeçalho da resposta e enviar o JSON com os dados do clima
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
