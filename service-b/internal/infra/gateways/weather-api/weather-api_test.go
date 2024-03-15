package weather_api_gateway

import (
	"os"
	"strings"
	"testing"
)

func init() {
	// ----- Read .env file ----- //
	data, err := os.ReadFile("../../../../../.env")
	if err != nil {
		panic(err)
	}
	dataSt := string(data)
	key := strings.Split(strings.Split(dataSt, "\n")[2], "=")[1]
	os.Setenv("WEATHER_API_KEY", key)
}

func TestShouldGetTempInfoByCity(t *testing.T) {
	gateway := NewWeatherAPIGateway(os.Getenv("WEATHER_API_KEY"))
	temp, err := gateway.GetTemperatureByCity("São Paulo")
	if err != nil {
		t.Errorf("Expected nil, got '%s'", err)
	}
	if temp.TempC == 0 {
		t.Errorf("Expected 0, got '%f'", temp.TempC)
	}
	if temp.TempF == 0 {
		t.Errorf("Expected 0, got '%f'", temp.TempF)
	}
	if temp.TempK == 0 {
		t.Errorf("Expected 0, got '%f'", temp.TempK)
	}
}
