package main

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/app"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/environment"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
)

func main() {
	environment.ReadEnvFiles(".env", "keyservice.env") // Load env files
	logger.ConfigureLoggerFromEnv()
	app.Start()

}
