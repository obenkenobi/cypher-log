package environment

import (
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func GetEnvVariable(key string) string {
	return os.Getenv(key)
}

func GetEnvVariableOrDefault(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func ReadEnvFiles(envFilesNames ...string) {
	for _, envFileName := range envFilesNames {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Infof("Failed to load environment file %v", envFileName)
		}
	}
}
