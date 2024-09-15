module github.com/thiagohmm/fulcycleTemperaturaPorCep

go 1.23.0

require (
	github.com/go-chi/chi/v5 v5.1.0
	go.opentelemetry.io/otel v1.30.0
	go.opentelemetry.io/otel/exporters/zipkin v1.30.0
	go.opentelemetry.io/otel/sdk v1.30.0
)

require github.com/felixge/httpsnoop v1.0.4 // indirect

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.55.0
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
