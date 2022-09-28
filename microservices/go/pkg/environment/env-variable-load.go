package environment

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

func GetEnvVariableAsListSplitByComma(key string) []string {
	rawVal := GetEnvVariable(key)
	return strings.Split(rawVal, ",")
}

func GetEnvVarAsBoolOrDefault(key string, defaultValue bool) bool {
	rawVal := GetEnvVariable(key)
	switch strings.ToLower(rawVal) {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
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

var _envFilesRead = false

func ReadEnvFiles(envFilesNames ...string) {
	if _envFilesRead {
		return
	}
	for _, envFileName := range envFilesNames {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Debugf("Did not load environment file %v", envFileName)
		}
	}
	_envFilesRead = true
}
