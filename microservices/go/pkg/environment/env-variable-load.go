package environment

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/utils"
	"os"
	"strconv"
	"time"

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

func GetEnvVarAsIntOrDefault(key string, defaultValue int) int {
	value := defaultValue
	if str := GetEnvVariable(key); utils.StringIsNotBlank(str) {
		if parsedInt, err := strconv.Atoi(str); err == nil {
			value = parsedInt
		}
	}
	return value
}

func GetEnvVarAsTimeDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	value := defaultValue
	if str := GetEnvVariable(key); utils.StringIsNotBlank(str) {
		if durationInt, err := strconv.ParseInt(str, 10, 64); err == nil {
			value = time.Duration(durationInt)
		}
	}
	return value
}

func ReadEnvFiles(envFilesNames ...string) {
	for _, envFileName := range envFilesNames {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Infof("Did not load environment file %v", envFileName)
		}
	}
}
