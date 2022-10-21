package environment

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

var _cachedAppEnvironment = ""

func GetAppEnvironment() string {
	if utils.StringIsBlank(_cachedAppEnvironment) {
		env := GetEnvVariableOrDefault(EnvVarKeyAppEnvironment, Development)
		switch env {
		case Development, Staging, Production:
			_cachedAppEnvironment = env
		default:
			_cachedAppEnvironment = Development
		}
	}
	return _cachedAppEnvironment
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
