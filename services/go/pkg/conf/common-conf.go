package conf

const EnvVariableEnvironment = "ENVIRONMENT"

const Development = "DEVELOPMENT"
const Staging = "STAGING"
const Production = "PRODUCTION"

type CommonConf interface {
	GetEnvironment() string
	IsDevelopment() bool
	IsStaging() bool
	IsProduction() bool
}

type commonConfImpl struct {
	environment string
}

func (c commonConfImpl) GetEnvironment() string {
	return c.environment
}

func (c commonConfImpl) IsDevelopment() bool {
	return c.environment == Development
}

func (c commonConfImpl) IsStaging() bool {
	return c.environment == Staging
}

func (c commonConfImpl) IsProduction() bool {
	return c.environment == Production
}

func NewCommonConf(envVarReader EnvVarReader, envVarKeyEnvironment string) CommonConf {
	return &commonConfImpl{environment: envVarReader.GetEnvVariableOrDefault(envVarKeyEnvironment, Development)}
}
