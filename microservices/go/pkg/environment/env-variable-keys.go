package environment

// General

const EnvVarKeyAppEnvironment = "ENVIRONMENT"

// Server

const EnvVarKeyAppServerPort = "APP_SERVER_PORT"
const EnvVarKeyGrpcServerPort = "GRPC_SERVER_PORT"

// Auth0

const EnvVarKeyAuth0IssuerUrl = "AUTH0_ISSUER_URL"
const EnvVarKeyAuth0Audience = "AUTH0_AUDIENCE"
const EnvVarKeyAuth0Domain = "AUTH0_DOMAIN"
const EnvVarKeyAuth0ClientId = "AUTH0_CLIENT_ID"
const EnvVarKeyAuth0ClientSecret = "AUTH0_CLIENT_SECRET"

// MongoDB

const EnvVarKeyMongoUri = "MONGO_URI"
const EnvVarMongoDBName = "MONGO_DB_NAME"
const EnvVarMongoConnTimeoutMS = "MONGO_CONNECTION_TIMEOUT_MS"

// redis

const EnvVarKeyRedisAddr = "REDIS_ADDRESS"
const EnvVarKeyRedisPassword = "REDIS_PASSWORD"
const EnvVarKeyRedisDB = "REDIS_DB"
