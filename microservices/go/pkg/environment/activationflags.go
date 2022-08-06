package environment

func ActivateAppServer() bool  { return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateAppServer, true) }
func ActivateGrpcServer() bool { return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateGrpcServer, true) }
func ActivateRabbitMqConsumer() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateRabbitMQConsumer, true)
}

func ActivateGRPCAuth() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateGRPCAuth, true)
}

func ActivateAppServerTLS() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateAppServerTLS, true)
}
