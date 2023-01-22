package environment

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"go.uber.org/atomic"
)

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

var _cachedAppEnv = atomic.NewString("")

func GetAppEnvironment() string {
	cachedEnv := _cachedAppEnv.Load()
	if utils.StringIsNotBlank(cachedEnv) {
		return cachedEnv
	}
	loadedEnv := GetEnvVarOrDefault(EnvVarKeyAppEnvironment, Development)
	var appEnv string
	switch loadedEnv {
	case Development, Staging, Production:
		appEnv = loadedEnv
	default:
		appEnv = Development
	}
	_cachedAppEnv.Store(appEnv)
	return appEnv
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
