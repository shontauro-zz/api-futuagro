package config

import (
	"log"
	"os"
	"strconv"
)

// DatabaseConf for modeling the configuration attributes for the database connection
type DatabaseConf struct {
	URI      string
	PoolSize uint16
	Name     string
}

// Config for modeling a global object with the global app configurations
type Config struct {
	Database DatabaseConf
	Port     string
}

// NewDefaultConfig return a config object with all application environment variables loaded
func NewDefaultConfig() *Config {
	return &Config{
		Database: DatabaseConf{
			URI:      getEnv("DB_URI", ""),
			PoolSize: uint16(getEnvAsUInt("DB_POOL_SIZE", 10)),
			Name:     getEnv("DB_NAME", ""),
		},
		Port: getEnv("APP_PORT", "3000"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(Key string, defaultVal string) string {
	if value, exists := os.LookupEnv(Key); exists {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err != nil {
		log.Printf("Error reading an environment variable %s %v\n", name, err)
	} else {
		return value
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an unsigned integer of 64 bytes or return a default value
func getEnvAsUInt(name string, defaultVal uint64) uint64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseUint(valueStr, 10, 64); err != nil {
		log.Printf("Error reading an environment variable %s %v\n", name, err)
	} else {
		return value
	}
	return defaultVal
}

// Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valStr); err != nil {
		log.Printf("Error reading an environment variable %s %v\n", name, err)
	} else {
		return value
	}
	return defaultVal
}
