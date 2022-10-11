package environment

// General

const EnvVarKeyAppEnvironment = "ENVIRONMENT"

// Boolean Activation Flags

const EnvVarKeyActivateAppServer = "ACTIVATE_APP_SERVER"
const EnvVarKeyActivateGrpcServer = "ACTIVATE_GRPC_SERVER"
const EnvVarActivateRabbitMQListener = "ACTIVATE_RABBITMQ_LISTENER"
const EnvVarActivateAppServerTLS = "ACTIVATE_APP_SERVER_TLS"
const EnvVarActivateGRPCAuth = "ACTIVATE_GRPC_AUTH"
const EnvVarActivateCronRunner = "ACTIVATE_CRON_RUNNER"

// SSL/TLS

const EnvVarKeyServerCertPath = "SERVER_CERT_PATH"
const EnvVarKeyServerKeyPath = "SERVER_KEY_PATH"
const EnvVarKeyCACertPath = "CA_CERT_PATH"
const EnvVarLoadCACert = "LOAD_CA_CERT"

// GRPC Client

const EnvVarKeyGrpcUserServiceAddress = "GRPC_USER_SERVICE_ADDRESS"
const EnvVarKeyGrpcKeyServiceAddress = "GRPC_KEY_SERVICE_ADDRESS"

// Server

const EnvVarKeyAppServerPort = "APP_SERVER_PORT"
const EnvVarKeyGrpcServerPort = "GRPC_SERVER_PORT"

// Auth0

const EnvVarKeyAuth0ApiAudience = "AUTH0_API_AUDIENCE"
const EnvVarKeyAuth0GrpcAudience = "AUTH0_GRPC_AUDIENCE"
const EnvVarKeyAuth0Domain = "AUTH0_DOMAIN"
const EnvVarKeyAuth0ClientCredentialsId = "AUTH0_CLIENT_CREDENTIALS_ID"
const EnvVarKeyAuth0ClientCredentialsSecret = "AUTH0_CLIENT_CREDENTIALS_SECRET"

// MongoDB

const EnvVarKeyMongoUri = "MONGO_URI"
const EnvVarMongoDBName = "MONGO_DB_NAME"
const EnvVarMongoConnTimeoutMS = "MONGO_CONNECTION_TIMEOUT_MS"

// RabbitMQ

const EnvVarRabbitMQUri = "RABBITMQ_URI"

// redis

const EnvVarKeyRedisAddr = "REDIS_ADDRESS"
const EnvVarKeyRedisPassword = "REDIS_PASSWORD"
const EnvVarKeyRedisDB = "REDIS_DB"
