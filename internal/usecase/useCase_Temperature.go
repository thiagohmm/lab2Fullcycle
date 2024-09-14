package usecase

import (
	"context"

	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/entity"
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/infraestructure"
)

type TemperatureInputDTO struct {
	Cep string `json:"cep"`
}

type TemperatureOutputDTO struct {
	City      string  `json:"city"`
	Celsius   float64 `json:"temp_C"`
	Farenheit float64 `json:"temp_F"`
	Kelvin    float64 `json:"temp_K"`
}

type TemperatureUseCase struct {
	Apiclient   infraestructure.GetTemperatureForCep
	Temperature entity.Temperature
}

func NewTemperatureUseCase(apiClient infraestructure.GetTemperatureForCep) *TemperatureUseCase {
	return &TemperatureUseCase{
		Apiclient:   apiClient,
		Temperature: entity.Temperature{}, // Inicializa a estrutura vazia
	}
}

func (t *TemperatureUseCase) Execute(ctx context.Context, input TemperatureInputDTO) (*TemperatureOutputDTO, error) {
	wheatherK, city, err := t.Apiclient.GetTemperatureByCep(ctx, input.Cep)
	if err != nil {
		return nil, err
	}

	temperature, err := entity.NewTemperature(city, wheatherK)
	if err != nil {
		return nil, err
	}

	return &TemperatureOutputDTO{
		City:      temperature.City,
		Celsius:   temperature.Celsius,
		Farenheit: temperature.Farenheit,
		Kelvin:    temperature.Kelvin,
	}, nil
}
