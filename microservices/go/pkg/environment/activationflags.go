package environment

func ActivateAppServer() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateAppServer, true)
}

func ActivateGrpcServer() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateGrpcServer, true)
}

func ActivateKafkaListener() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateKafkaListener, true)
}

func ActivateGRPCAuth() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateGRPCAuth, true)
}

func ActivateAppServerTLS() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateAppServerTLS, true)
}

func ActivateCronRunner() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateCronRunner, true)
}
