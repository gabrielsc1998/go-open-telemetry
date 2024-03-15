package weather_api_gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type WeatherAPIGateway struct {
	key string
}

type WeatherAPIGatewayDtoOutput struct {
	TempC float64
	TempF float64
	TempK float64
}

type WeatherAPICurrent struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
}

type WeatherAPIResponse struct {
	Current WeatherAPICurrent `json:"current"`
}

func NewWeatherAPIGateway(key string) *WeatherAPIGateway {
	return &WeatherAPIGateway{key: key}
}

func (w *WeatherAPIGateway) GetTemperatureByCity(city string) (*WeatherAPIGatewayDtoOutput, error) {
	adjustedCityName, err := w.adjustCityNameForQuery(city)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(w.url(adjustedCityName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error")
	}

	data := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	weatherAPIResponse := WeatherAPIResponse{}
	err = json.Unmarshal(jsonData, &weatherAPIResponse)
	if err != nil {
		return nil, err
	}

	return &WeatherAPIGatewayDtoOutput{
		TempC: weatherAPIResponse.Current.TempC,
		TempF: weatherAPIResponse.Current.TempF,
		TempK: weatherAPIResponse.Current.TempC + 273.15,
	}, nil
}

func (w *WeatherAPIGateway) url(city string) string {
	return fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", w.key, city)
}

func (w *WeatherAPIGateway) adjustCityNameForQuery(city string) (string, error) {
	// Replaces accents and special characters
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	city, _, err := transform.String(t, city)
	if err != nil {
		return "", err
	}

	// Replaces spaces with dashes
	city = strings.ToLower(city)
	splittedCity := strings.Split(city, " ")
	return strings.Join(splittedCity, "-"), nil
}
