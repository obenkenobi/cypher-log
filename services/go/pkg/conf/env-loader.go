package conf

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

// EnvVarReader Reads environment variables
type EnvVarReader interface {
	GetEnvVariable(key string) string
	GetEnvVariableOrDefault(key string, defaultVal string) string
}

type EnvVarReaderImpl struct {
}

func (e EnvVarReaderImpl) GetEnvVariable(key string) string {
	return os.Getenv(key)
}

func (e EnvVarReaderImpl) GetEnvVariableOrDefault(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func NewEnvVariableAccessor(envFilesNames []string) EnvVarReader {
	for _, envFileName := range envFilesNames {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Info("Failed to load env file %v", envFileName)
		}
	}
	return EnvVarReaderImpl{}
}
