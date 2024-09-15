package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/infraestructure"
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/usecase" // Add this line to import the usecase package
	"github.com/thiagohmm/fulcycleTemperaturaPorCep/internal/web"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() func() {

	// Configurar o exportador Zipkin
	endpoint := os.Getenv("ZIPKIN_ENDPOINT")
	exporter, err := zipkin.New(endpoint)
	if err != nil {
		log.Fatalf("Falha ao criar exportador Zipkin: %v", err)
	}

	// Criar um tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Opcional: Sempre traçar para teste
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("weather-service"),
		)),
	)

	// Registrar o Tracer Provider como o global
	otel.SetTracerProvider(tp)

	return func() {
		// Shutdown
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Erro ao encerrar tracer provider: %v", err)
		}
	}

}

func main() {

	// Inicializar o tracer
	shutdown := initTracer()
	defer shutdown()
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
	r.Post("/weather", otelhttp.NewHandler(http.HandlerFunc(handler.GetWeather), "/weather").ServeHTTP)

	// Iniciar o servidor HTTP
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
