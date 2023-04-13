package util

import (
	"os"

	"github.com/joho/godotenv"
)

func SetEnvIfMissing(envName string, value string) {
	if os.Getenv(envName) == "" {
		os.Setenv(envName, value)
	}
}

func SetupDefaultEnv(debug bool) {
	if debug {
		err := godotenv.Load(".env.test.local")
		if err != nil {
			panic("unable to load environment variables")
		}
	}

	// Define ONLY NON-SENSITIVE default environment variables here
	// Use a secret manager service for any sensitive information
	defaultEnv := map[string]string{}

	for k, v := range defaultEnv {
		SetEnvIfMissing(k, v)
	}
}
