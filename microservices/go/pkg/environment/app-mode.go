package environment

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

var cachedAppEnvironment = ""

func GetAppEnvironment() string {
	if utils.StringIsBlank(cachedAppEnvironment) {
		return GetEnvVariableOrDefault(EnvVarKeyAppEnvironment, Development)
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
