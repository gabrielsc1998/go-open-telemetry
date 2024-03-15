package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gabrielsc1998/go-open-telemetry/package/domain"
)

type TempByCepController struct {
}

type TempByCepDtoInput struct {
	Cep string `json:"cep"`
}

func NewTempByCepController() *TempByCepController {
	return &TempByCepController{}
}

func (c *TempByCepController) Handle(w http.ResponseWriter, r *http.Request) {
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

	resp, err := http.Get("http://service-b:8081/temp-by-cep?cep=" + cep.Value)
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
