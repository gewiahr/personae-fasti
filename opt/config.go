package opt

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Conf struct {
	App struct {
		Port  string `json:"port"`
		Debug bool   `json:"debug"`
	} `json:"app"`
	DB struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"db"`
}

func InitConfig() *Conf {

	config := new(Conf)

	envConfigPath := os.Getenv("CONFIG_PATH")
	if envConfigPath == "" {
		envConfigPath = "./opt"
	}

	envConfigName := os.Getenv("CONFIG_NAME")
	if envConfigName == "" {
		envConfigName = "config"
	}

	viper.SetConfigName(envConfigName)
	viper.SetConfigType("json")
	viper.AddConfigPath(envConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return config

}
