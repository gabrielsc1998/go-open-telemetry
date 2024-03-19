package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gabrielsc1998/go-open-telemetry/package/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type TempByCepController struct {
	tracer          trace.Tracer
	otelRequestName string
}

type TempByCepDtoInput struct {
	Cep string `json:"cep"`
}

func NewTempByCepController(tracer trace.Tracer, otelRequestName string) *TempByCepController {
	return &TempByCepController{
		tracer:          tracer,
		otelRequestName: otelRequestName,
	}
}

func (c *TempByCepController) Handle(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := c.tracer.Start(ctx, c.otelRequestName)
	defer span.End()

	dto := TempByCepDtoInput{}
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dto.Cep == "" {
		http.Error(w, "zipcode is required", http.StatusBadRequest)
		return
	}

	cep, err := domain.NewCep(dto.Cep)
	if err != nil {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://app:8081/temp-by-cep?cep="+cep.Value,
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var body bytes.Buffer
	body.ReadFrom(resp.Body)

	if resp.StatusCode != http.StatusOK {
		errorMsg := strings.Trim(body.String(), "\n")
		if errorMsg == "cep not found" {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			return
		}
		if errorMsg == "invalid zipcode" {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	w.Write(body.Bytes())
}
