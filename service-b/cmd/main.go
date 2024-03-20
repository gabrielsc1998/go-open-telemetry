package main

import (
	"context"
	"log"

	"github.com/gabrielsc1998/go-open-telemetry/configs"
	opentelemetry "github.com/gabrielsc1998/go-open-telemetry/package/infra/open-telemetry"
	"github.com/gabrielsc1998/go-open-telemetry/package/infra/server"
	"github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/controllers"
	viacep_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/viacep"
	weather_api_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/weather-api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	provider := opentelemetry.NewProvider(
		config.ServiceBOtelServiceName,
		config.ServiceBOtelExporterOTLPEndpoint,
	)
	shutdown, err := provider.Start()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatal("failed to shutdown TraceProvider: %w", err)
		}
	}()

	tracer := opentelemetry.NewTracer(config.ServiceBOtelServiceName)

	webserver := server.NewServer(config.ServiceBWebServerPort)

	viacepGateway := viacep_gateway.NewViaCepGateway()
	weatherApiGateway := weather_api_gateway.NewWeatherAPIGateway(config.WeatherApiKey)

	tempByCepController := controllers.NewTempByCepController(
		viacepGateway,
		weatherApiGateway,
		tracer,
		config.ServiceBOtelRequestName,
	)
	webserver.AddRoute("GET", "/temp-by-cep", tempByCepController.Handle)
	webserver.AddHandler("GET", "/metrics", promhttp.Handler())

	webserver.Start()
}
