package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Debug        bool
	MetircsPort  string
	InternalHost string
	InternalIP   string
	ExternalHost string
	ExternalIP   string
	Delay        int
	Timeout      int
}

var CFG Config

func LoadConfigFromEnv() Config {
	config := Config{
		Debug:        parseEnvBool("DEBUG"),
		MetircsPort:  getEnvOrDefault("PORT", "8080"),
		InternalHost: getEnvOrDefault("INTERNAL_HOST", "kubernetes.default.svc.cluster.local"),
		InternalIP:   getEnvOrDefault("INTERNAL_IP", "10.43.0.1"),
		ExternalHost: getEnvOrDefault("EXTERNAL_HOST", "a.root-servers.net"),
		ExternalIP:   getEnvOrDefault("EXTERNAL_IP", "198.41.0.4"),
		Delay:        parseEnvInt("DELAY", 5),
		Timeout:      parseEnvInt("TIMEOUT", 5),
	}

	CFG = config

	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	if err != nil {
		log.Printf("Failed to parse environment variable %s: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func parseEnvBool(key string) bool {
	value := os.Getenv(key)
	boolValue := false
	if value == "true" {
		boolValue = true
	}
	return boolValue
}
