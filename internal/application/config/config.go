package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     int    `mapstructure:"DBPORT"`
	DBUser     string `mapstructure:"DBPWD"`
	DBPassword string `mapstructure:"DBPWD"`
	DBName     string `mapstructure:"DBNAME"`
	AppPort    int    `mapstructure:"PORT"`
}

var app Config

// func Load() (*Config, error) {
// 	port, err := strconv.Atoi(os.Getenv("PORT"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	DBPORT, _ := strconv.Atoi(os.Getenv("DBPORT"))
// 	var c = Config{
// 		DBHost:     getEnv("DB_HOST", "localhost"),
// 		DBPort:     DBPORT,
// 		DBUser:     getEnv("DB_USER", "postgres"),
// 		DBPassword: getEnv("DB_PASSWORD", ""),
// 		DBName:     getEnv("DB_NAME", "bookstore"),
// 		AppPort:    port,
// 	}
// 	return &c, nil

// }

func Load() (*Config, error) {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return nil, errors.New("PORT environment variable is not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PORT: %v", err)
	}

	DBPORTStr := os.Getenv("DBPORT")
	if DBPORTStr == "" {
		return nil, errors.New("DBPORT environment variable is not set")
	}

	DBPORT, err := strconv.Atoi(DBPORTStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DBPORT: %v", err)
	}

	var c = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     DBPORT,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "bookstore"),
		AppPort:    port,
	}
	return &c, nil
}

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func (c *Config) DBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
