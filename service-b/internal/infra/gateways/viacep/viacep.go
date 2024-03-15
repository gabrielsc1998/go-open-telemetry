package viacep_gateway

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ViaCepGateway struct{}

type ViaCepGatewayDtoOutput struct {
	Code     string `json:"cep"`
	State    string `json:"uf"`
	City     string `json:"localidade"`
	District string `json:"bairro"`
	Address  string `json:"logradouro"`
}

func NewViaCepGateway() *ViaCepGateway {
	return &ViaCepGateway{}
}

func (v *ViaCepGateway) GetAddressByCep(cep string) (*ViaCepGatewayDtoOutput, error) {
	resp, err := http.Get("https://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data["erro"] != nil && data["erro"].(bool) {
		return nil, errors.New("cep not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error")
	}

	dto := &ViaCepGatewayDtoOutput{
		Code:     data["cep"].(string),
		State:    data["uf"].(string),
		City:     data["localidade"].(string),
		District: data["bairro"].(string),
		Address:  data["logradouro"].(string),
	}
	return dto, nil
}
