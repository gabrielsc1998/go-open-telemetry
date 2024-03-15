package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielsc1998/go-open-telemetry/package/domain"
	viacep_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/viacep"
	weather_api_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/weather-api"
)

type TempByCepController struct {
	viacepGateway     *viacep_gateway.ViaCepGateway
	weatherApiGateway *weather_api_gateway.WeatherAPIGateway
}

type TempByCepControllerDtoOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewTempByCepController(viacepGateway *viacep_gateway.ViaCepGateway, weatherApiGateway *weather_api_gateway.WeatherAPIGateway) *TempByCepController {
	return &TempByCepController{
		viacepGateway:     viacepGateway,
		weatherApiGateway: weatherApiGateway,
	}
}

func (c *TempByCepController) Handle(w http.ResponseWriter, r *http.Request) {
	receivedCep := r.URL.Query().Get("cep")

	_, err := domain.NewCep(receivedCep)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	address, err := c.viacepGateway.GetAddressByCep(receivedCep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temperature, err := c.weatherApiGateway.GetTemperatureByCity(address.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temperatureDto := TempByCepControllerDtoOutput{
		City:  address.City,
		TempC: temperature.TempC,
		TempF: temperature.TempF,
		TempK: temperature.TempK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temperatureDto)
}
