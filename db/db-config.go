package db

import "os"

type Config struct {
	DbHost     string
	DbName     string
	DbUser     string
	DbPassword string
	DbPort     string
}

func NewConfig() Config {
	return Config{
		DbHost:     getEnv("DB_HOST", "localhost"),
		DbPort:     getEnv("DB_PORT", "5432"),
		DbName:     getEnv("DB_NAME", "test"),
		DbUser:     getEnv("DB_USER", "test"),
		DbPassword: getEnv("DB_PASSWORD", "test"),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue

}
