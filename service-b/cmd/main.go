package main

import (
	"github.com/gabrielsc1998/go-open-telemetry/configs"
	"github.com/gabrielsc1998/go-open-telemetry/package/infra/server"
	"github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/controllers"
	viacep_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/viacep"
	weather_api_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/weather-api"
)

func main() {
	config, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	server := server.NewServer(config.WebServerPortServiceB)

	viacepGateway := viacep_gateway.NewViaCepGateway()
	weatherApiGateway := weather_api_gateway.NewWeatherAPIGateway(config.WeatherApiKey)

	tempByCepController := controllers.NewTempByCepController(viacepGateway, weatherApiGateway)
	server.AddRoute("GET", "/temp-by-cep", tempByCepController.Handle)

	server.Start()
}
