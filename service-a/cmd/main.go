package main

import (
	"github.com/gabrielsc1998/go-open-telemetry/configs"
	"github.com/gabrielsc1998/go-open-telemetry/package/infra/server"
	"github.com/gabrielsc1998/go-open-telemetry/service-a/internal/infra/controllers"
)

func main() {
	config, err := configs.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	server := server.NewServer(config.WebServerPortServiceA)

	tempByCepController := controllers.NewTempByCepController()
	server.AddRoute("POST", "/temp-by-cep", tempByCepController.Handle)

	server.Start()
}
