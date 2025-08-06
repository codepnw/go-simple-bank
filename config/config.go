package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	APP *app
	DB  *db
}

type db struct {
	User string
	Pass string
	DB   string
	Host string
	Port string
	SSL  string
}

type app struct {
	Version string
	Port    string
}

func LoadEnvConfig(configFile string) (*EnvConfig, error) {
	if err := godotenv.Load(configFile); err != nil {
		return nil, err
	}

	env := &EnvConfig{
		APP: &app{
			Version: getEnvString("APP_VERSION", "v1"),
			Port:    getEnvString("APP_PORT", ":8080"),
		},
		DB: &db{
			User: getEnvString("POSTGRES_USER", "postgres"),
			Pass: getEnvString("POSTGRES_PASS", ""),
			DB:   getEnvString("POSTGRES_DB", "simple_bank"),
			Host: getEnvString("POSTGRES_HOST", "localhost"),
			Port: getEnvString("POSTGRES_PORT", "5432"),
			SSL:  getEnvString("POSTGRES_SSL", "disable"),
		},
	}

	return env, nil
}

func getEnvString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
