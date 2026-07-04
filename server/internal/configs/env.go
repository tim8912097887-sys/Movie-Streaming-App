package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configs struct{
	Addr string
	RefreshTokenSecret string
	AccessTokenSecret string
}

func InitConfigs() Configs {
	// Fail silently for production
	_ = godotenv.Load()
	
	return Configs{
		Addr: getEnv("ADDR", ":8080"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET", "secret"),
		AccessTokenSecret: getEnv("ACCESS_TOKEN_SECRET", "secret"),
	}
}

func getEnv(key string, defaultValue string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

func getEnvFromInt(key string, defaultValue int) int {

	if value, ok := os.LookupEnv(key); ok {
		num, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}

		return num
	}

	return defaultValue
}