package main

import (
	"context"
	"log"

	"github.com/gabrielsc1998/go-open-telemetry/configs"
	opentelemetry "github.com/gabrielsc1998/go-open-telemetry/package/infra/open-telemetry"
	"github.com/gabrielsc1998/go-open-telemetry/package/infra/server"
	"github.com/gabrielsc1998/go-open-telemetry/service-a/internal/infra/controllers"
)

func main() {
	config, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	provider := opentelemetry.NewProvider(
		config.ServiceAOtelServiceName,
		config.ServiceAOtelExporterOTLPEndpoint,
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

	tracer := opentelemetry.NewTracer(config.ServiceAOtelServiceName)

	webserver := server.NewServer(config.ServiceAWebServerPort)
	webserver.AddMiddleware(server.RequestIdMiddleware)

	tempByCepController := controllers.NewTempByCepController(tracer, config.ServiceAOtelRequestName)
	webserver.AddRoute("POST", "/temp-by-cep", tempByCepController.Handle)

	webserver.Start()
}
