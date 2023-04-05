package util

import (
	"os"

	"github.com/joho/godotenv"
)

func setEnvIfMissing(envVarName string, val string) {
	if os.Getenv(envVarName) == "" {
		os.Setenv(envVarName, val)
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
		setEnvIfMissing(k, v)
	}
}
