package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	APIToken            string
	IBGEBaseURL         string
	IBGETimeoutSecond   int
	IBGECacheTTLSeconds int
}

func Load() Config {
	cfg := Config{
		Port:                getEnv("PORT", "8080"),
		APIToken:            os.Getenv("API_TOKEN"),
		IBGEBaseURL:         getEnv("IBGE_BASE_URL", "https://servicodados.ibge.gov.br/api"),
		IBGECacheTTLSeconds: getEnvInt("IBGE_CACHE_TTL_SECONDS", 21600), // 6h
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
