package environment

import "github.com/obenkenobi/cypher-log/services/go/pkg/utils"

const DefaultEnvVarKeyAppEnv = "ENVIRONMENT"

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

var environmentVarKeyAppEnvironment = DefaultEnvVarKeyAppEnv
var CachedEnvironment = ""

func GetEnvironment() string {
	if utils.StringIsBlank(CachedEnvironment) {
		return GetEnvVariableOrDefault(environmentVarKeyAppEnvironment, Development)
	}
	return CachedEnvironment
}

func IsStaging() bool {
	return GetEnvironment() == Staging
}

func IsDevelopment() bool {
	return GetEnvironment() == Development
}

func IsProduction() bool {
	return GetEnvironment() == Production
}

func SetEnvVarKeyForAppEnvironment(key string) {
	environmentVarKeyAppEnvironment = key
}
