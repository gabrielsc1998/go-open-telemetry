package viacep_gateway

import "testing"

func TestShouldGetAddressByCep(t *testing.T) {
	gateway := NewViaCepGateway()
	address, err := gateway.GetAddressByCep("01031970")
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err)
	}
	if address.Code != "01031-970" {
		t.Errorf("Expected 01031-970, got '%s'", address.Code)
	}
	if address.State != "SP" {
		t.Errorf("Expected SP, got '%s'", address.State)
	}
	if address.City != "São Paulo" {
		t.Errorf("Expected São Paulo, got '%s'", address.City)
	}
	if address.District != "Centro" {
		t.Errorf("Expected Centro, got '%s'", address.District)
	}
	if address.Address != "Praça do Correio" {
		t.Errorf("Expected Praça do Correio, got '%s'", address.Address)
	}
}

func TestShouldReturnErrorWhenTheCepIsInvalid(t *testing.T) {
	gateway := NewViaCepGateway()
	_, err := gateway.GetAddressByCep("123")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestShouldReturnErrorWhenTheCepIsNotFound(t *testing.T) {
	gateway := NewViaCepGateway()
	_, err := gateway.GetAddressByCep("88160000")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	print(err.Error())
	if err.Error() != "cep not found" {
		t.Error("Expected cep not found")
	}
}
