package environment

import "github.com/obenkenobi/cypher-log/services/go/pkg/utils"

const DefaultEnvVarKeyAppEnv = "ENVIRONMENT"

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

var environmentVarKeyAppEnvironment = DefaultEnvVarKeyAppEnv
var cachedAppEnvironment = ""

func GetAppEnvironment() string {
	if utils.StringIsBlank(cachedAppEnvironment) {
		return GetEnvVariableOrDefault(environmentVarKeyAppEnvironment, Development)
	}
	return cachedAppEnvironment
}

func IsStaging() bool {
	return GetAppEnvironment() == Staging
}

func IsDevelopment() bool {
	return GetAppEnvironment() == Development
}

func IsProduction() bool {
	return GetAppEnvironment() == Production
}

func SetEnvVarKeyForAppEnvironment(key string) {
	environmentVarKeyAppEnvironment = key
}
