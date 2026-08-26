package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Main struct {
	App          *APIConfig          `json:"app"`
	Auth         *AuthConfig         `json:"auth"`
	DB           *DBConfig           `json:"db"`
	ImageStorage *ImageStorageConfig `json:"imageStorage"`
}

type ImageStorageConfig struct {
	Endpoint       string `json:"endpoint"`
	Region         string `json:"region"`
	Bucket         string `json:"bucket"`
	AccessKey      string `json:"accessKey"`
	SecretKey      string `json:"secretKey"`
	PublicBaseURL  string `json:"publicBaseUrl"`
	ForcePathStyle bool   `json:"forcePathStyle"`
}

type APIConfig struct {
	Port               string   `json:"port"`
	Debug              bool     `json:"debug"`
	Environment        string   `json:"environment"`
	AllowedOrigins     []string `json:"allowedOrigins"`
	AllowCredentials   bool     `json:"allowCredentials"`
	ReadTimeout        int      `json:"readTimeout"`
	WriteTimeout       int      `json:"writeTimeout"`
	ImageUploadTimeout int      `json:"imageUploadTimeout"`
	MigrateOnStart     bool     `json:"migrateOnStart"`
}

type AuthConfig struct {
	JWTSecret             string `json:"jwtSecret"`
	JWTTokenLifetimeHours int64  `json:"jwtTokenLifetimeHours"`
	BotToken              string `json:"botToken"`
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func InitConfig() *Main {

	config := new(Main)

	envConfigPath := os.Getenv("CONFIG_PATH")
	if envConfigPath == "" {
		envConfigPath = "./configs"
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
