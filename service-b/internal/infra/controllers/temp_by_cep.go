package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielsc1998/go-open-telemetry/package/domain"
	viacep_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/viacep"
	weather_api_gateway "github.com/gabrielsc1998/go-open-telemetry/service-b/internal/infra/gateways/weather-api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TempByCepController struct {
	tracer            trace.Tracer
	otelRequestName   string
	viacepGateway     *viacep_gateway.ViaCepGateway
	weatherApiGateway *weather_api_gateway.WeatherAPIGateway
}

type TempByCepControllerDtoOutput struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewTempByCepController(
	viacepGateway *viacep_gateway.ViaCepGateway,
	weatherApiGateway *weather_api_gateway.WeatherAPIGateway,
	tracer trace.Tracer,
	otelRequestName string,
) *TempByCepController {
	return &TempByCepController{
		viacepGateway:     viacepGateway,
		weatherApiGateway: weatherApiGateway,
		tracer:            tracer,
		otelRequestName:   otelRequestName,
	}
}

func (c *TempByCepController) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	_, span := c.tracer.Start(ctx, c.otelRequestName)
	defer span.End()

	receivedCep := r.URL.Query().Get("cep")

	_, err := domain.NewCep(receivedCep)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	// ----- Get Address ----- //

	_, spanGetAddress := c.tracer.Start(ctx, c.otelRequestName+"-get-address")

	address, err := c.viacepGateway.GetAddressByCep(receivedCep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spanGetAddress.End()

	// ----- Get Temperature ----- //

	_, spanGetTemp := c.tracer.Start(ctx, c.otelRequestName+"-get-temp")

	temperature, err := c.weatherApiGateway.GetTemperatureByCity(address.City)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	spanGetTemp.End()

	// ----- Response ----- //

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
