package environment

func ActivateAppServer() bool  { return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateAppServer, true) }
func ActivateGrpcServer() bool { return GetEnvVarAsBoolOrDefault(EnvVarKeyActivateGrpcServer, true) }
func ActivateRabbitMqConsumer() bool {
	return GetEnvVarAsBoolOrDefault(EnvVarActivateRabbitMQConsumer, true)
}
