package configs

import (
	"github.com/spf13/viper"
)

type Conf struct {
	WebServerPortServiceA string `mapstructure:"WEB_SERVER_PORT_SERVICE_A"`
	WebServerPortServiceB string `mapstructure:"WEB_SERVER_PORT_SERVICE_B"`
	WeatherApiKey         string `mapstructure:"WEATHER_API_KEY"`
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
