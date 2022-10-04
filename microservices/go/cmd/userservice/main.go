package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/app"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
)

func main() {
	environment.ReadEnvFiles(".env", "userservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()
	app.Start()
}
