package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	// General configs
	WeatherApiKey string `mapstructure:"WEATHER_API_KEY"`

	// Service A configs
	ServiceAWebServerPort            string `mapstructure:"SERVICE_A_WEB_SERVER_PORT"`
	ServiceAOtelRequestName          string `mapstructure:"OTEL_REQUEST_NAME_SERVICE_A"`
	ServiceAOtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME_SERVICE_A"`
	ServiceAOtelExporterOTLPEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT_SERVICE_A"`

	// Service B configs
	ServiceBWebServerPort            string `mapstructure:"SERVICE_B_WEB_SERVER_PORT"`
	ServiceBOtelRequestName          string `mapstructure:"OTEL_REQUEST_NAME_SERVICE_B"`
	ServiceBOtelServiceName          string `mapstructure:"OTEL_SERVICE_NAME_SERVICE_B"`
	ServiceBOtelExporterOTLPEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT_SERVICE_B"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigType("env")
	viper.AddConfigPath("../../../../..")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
