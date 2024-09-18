package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
)

func LoadEnv(app *pocketbase.PocketBase) error {
	logger := app.Logger().WithGroup("LoadEnv")

	// TODO: make sure to add env in Dockerfile
	env := os.Getenv("BRFACTORY_ENV")
	if env == "" {
		env = "development"
	}

	// godotenv.Load(".env." + env + ".local")
	// if "test" != env {
	//   godotenv.Load(".env.local")
	// }

	if env == "development" {
		envPath := ".env." + env

		logger.Info("Loading env file: " + envPath)

		return godotenv.Load(envPath)
	}

	return nil
}

type EnvVariables struct {
	IsDevelopment   bool
	IsProduction    bool
	FileURLPrefix   string
	IGServiceURL    string
	IGServiceSecret string
}

func ParseEnv() (*EnvVariables, error) {
	environment, err := getEnv("ENVIRONMENT")
	if err != nil {
		return nil, err
	}

	fileURLPrefix, err := getEnv("FILE_URL_PREFIX")
	if err != nil {
		return nil, err
	}

	igServiceURL, err := getEnv("IG_SERVICE_URL")
	if err != nil {
		return nil, err
	}

	igServiceSecret, err := getEnv("IG_SERVICE_SECRET")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:   environment == "development",
		IsProduction:    environment == "production",
		FileURLPrefix:   fileURLPrefix,
		IGServiceURL:    igServiceURL,
		IGServiceSecret: igServiceSecret,
	}, nil
}

// func getEnvInt(key string) (int64, error) {
// 	val, err := getEnv(key)
// 	if err != nil {
// 		return 0, err
// 	}

// 	valInt, err := strconv.ParseInt(val, 10, 64)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return valInt, nil
// }

func getEnv(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		log.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
