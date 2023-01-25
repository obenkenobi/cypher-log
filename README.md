# Cypher Log
A passcode based encrypted narchitecture 
built on a microservice aarchitecture.

## Getting Started

### Technology Requirements
- Docker
- Go 1.18
- Node version 16.14.0
- Yarn
- Any Bash / Shell / ZSH terminal
- [Protocol Buffer compiler](https://developers.google.com/protocol-buffers), `protoc`, [version 3](https://developers.google.com/protocol-buffers/docs/proto3)
- VS Code or Goland
- Make

### Install Go Plugins
In order to support code generation and other features for developing in Go, install the following plugins:
```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
$ go install github.com/smartystreets/goconvey
$ go install github.com/google/wire/cmd/wire@latest
```
### Generate Development TLS Certificates and Keys
To generate self-signed certificates for your backend,
go to `/dev/certs` in your terminal and run `make certs`. 

### Environment Variables
Go to the directory `microservices/go` and `migrations/mongo`. Take each env file in the format `sample.*.env` and copy 
it with the same name but without the `sample.` prefix. For example `sample.env` becomes `.env` and 
`sample.userservice.env` becomes `userservice.env`. Fill in any necessary values intentionally left out. For Auth0 
credentials, either set up your own Auth0 realm or use an existing one.

### Docker Dev Dependencies
This project requires dependencies such as MongoDB and Redis. To quickly install such files, a docker compose file 
exists in `dev/docker`. Go to `dev/docker` and run `docker-compose up` to run your dependencies. 
Add the `-d` argument to run in detach mode. Alternatively you can use your IDE to run `docker-compose`.

#### MongoDB hosts
MongoDB replica sets in docker require your hosts to be updated.

If you use windows, open the file `C:\Windows\System32\drivers\etc\hosts`
and add `127.0.0.1 mongo1 mongo2 mongo3` to the file.

If on a linux server, use the hostname provided by the docker compose file <br>
e.g. `HOSTNAME = mongo1, mongo2, mongo3`

If on MacOS add the following to your `/etc/hosts` file and use localhost as the HOSTNAME.
```
127.0.0.1  mongo1
127.0.0.1  mongo2
127.0.0.1  mongo3
```

### Database Migrations
#### Mongo Migrations
Move into the directory `microservices/nodejs/mongo-migrator`. Here you will find a node.js project that migrates your 
mongodb database. You will also see a `sample.env` file and a `migrate-up-all-local.sh` file. You can either migrate 
via the terminal or an IDE, which in this case is Goland. We will show how to do it via both options. 
It is recommended regardless to read the terminal section even if you are using your IDE as it will inform 
how to set up your run configuration.
###### 1. Plain terminal
In this case you will need to copy your `sample.env` file into a `.env` file. You will notice `MONGO_URI` is 
already filled with a connection string to point your locally deployed mongoDB cluster. Of course in staging/production 
the value will be different. The variables `MONGO_DB_NAME` and `MIGRATION_DIRECTORY` are empty. 
`MONGO_DB_NAME` refers to the name of the database you will be migrating too.
`MIGRATION_DIRECTORY` refers to any subdirectory under `src/` which will contain your migration scripts.
It is recommended `MONGO_DB_NAME` will be the same as `MIGRATION_DIRECTORY` as this will 
make deploying your microservices easier. 
Then to migrate, the command is simple: `yarn run migrate:up`.
There are other yarn commands to be aware of:
- `migrate:new` Creates a new migration file under `src/${MIGRATION_DIRECTORY}`.
- `migrate:down:last` Undoes your last migration.
- `migrate:down:all` Undoes all of your migrations.
- `migrate:status` Show the status of the migrations.
You can look into `package.json` to see more commands but these are the essentials.

Here is the contents of an example `.env` file for reference:
```shell
MONGO_URI=mongodb://localhost:27017,localhost:27018,localhost:27019/?replicaSet=rs0
MONGO_DB_NAME=keys
MIGRATION_DIRECTORY=keys
```
###### 2. Shell Script
To automatically run migrations for all of your databases, just run the command `sh migrate-up-all-local.sh`.
`migrate-up-all-local.sh` will run all of your migrations for all databases. Very useful if you do not want to use
the yarn commands.
###### 2. IDE (Goland)
Since you have much of the prerequisites to now run your migrations in Goland, it is now important to explain one more
environment variable not mentioned in your sample.env file. `IGNORE_ENV_FILE`, which have the .env file ignored 
if set to `true`. The previous shell script method used the same variable in practice.

To create a run configuration, just open `package.json`. You will see a bunch of green arrows next to each 
migration command. Just right-click the green arrow and click modify run configuration. It is recommended to start with
`migrate:up` but any other command works as well. In your configuration,
go to where it says *Environment* and then add the environment variables `MONGO_URI`, 
`MONGO_DB_NAME`, and `MIGRATION_DIRECTORY` just like how they were set in the terminal method. Then add the variable
`IGNORE_ENV_FILE` and set it to equal `true`. Make sure you are using yarn to run your migration command. 
Click *apply*. You now know how to use the IDE to run any of the yarn commands.
###### Environment variables
Here is a table of all the environment variables used in the mongo migrator:

| Variable            | Description                                                                                   | Default value |
|---------------------|-----------------------------------------------------------------------------------------------|---------------|
| MONGO_URI           | The connection string to your mongodb instance/cluster.                                       |               |
| MONGO_DB_NAME       | Is the database name, it is recommended that it's value is the same as `MIGRATION_DIRECTORY`. |               |
| MIGRATION_DIRECTORY | Is a migration directory under the `src` directory.                                           |               |
| IGNORE_ENV_FILE     | If set to `true`, your .env file is ignored. This is meant to be set outside the `.env` file. |               |

### Go Microservices
Move into the directory `microservices/go`. You will notice that there is a single Go project. 
Now this single Go project is actually a collection of microservices that share a lot of common code. This shared code 
is in the `pkg` directory, and it helps facilitate the communication between each microservice and provide common
logic and abstractions that helps avoid writing boilerplate code. It acts as a shared library/package. Underneath the
`cmd` directory, you will see more subdirectories each representing a different microservice. For example the
subdirectory *userservice* represents a microservice.
#### Environment Variables
Before you even run your Go project, you need to set your environment variables. You
will see an env file with the name `sample.env` and env files in the format 
`sample.${subdirectory of cmd directoryu}.env`. Copy each of the file, removing the
prefix `sample.`
For example:
```shell
cp sample.env .env
cp sample.userservice.env userservice.env
```
In your env files, you can see what each of the environment variables are already filled 
and which ones to fill out.

Here is a table of all the environment variables:

| Variable                        | Description                                                                                                                                 | Default value |
|---------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|---------------|
| ENVIRONMENT                     | Represents the lifecycle environment the app is supposed to run on. Could be `DEVELOPMENT`, `STAGING`, or `PRODUCTION`.                     | DEVELOPMENT   |
| ACTIVATE_APP_SERVER             | Boolean flag that activates your HTTP app server (used for REST, static web pages, etc)                                                     | true          |
| ACTIVATE_GRPC_SERVER            | Boolean flag that activates your GRPC server                                                                                                | true          |
| ACTIVATE_RABBITMQ_LISTENER      | Boolean flag that activates your RabbitMQ listener                                                                                          | true          |
| ACTIVATE_APP_SERVER_TLS         | Boolean flag that activates TLS protection is added to your HTTP app server                                                                 | true          |
| ACTIVATE_GRPC_AUTH              | Boolean flag that activates authentication for your GRPC server and client                                                                  | true          |
| ACTIVATE_CRON_RUNNER            | Boolean flag that activates a background task designed to run cron jobs                                                                     | true          |
| SERVER_CERT_PATH                | Path to your TLS certificate file for your servers                                                                                          |               |
| SERVER_KEY_PATH                 | Path to your TLS private key file for your servers                                                                                          |               |
| CA_CERT_PATH                    | Path to your certificate authority file. This is used for self signed certificates.                                                         |               |
| LOAD_CA_CERT                    | Boolean flag where if true, a certificate authority will be loaded from the `CA_CERT_PATH` flag. This is used for self signed certificates. | false         |
| GRPC_USER_SERVICE_ADDRESS       | URI to the GRPC server running within the user service                                                                                      |               |
| GRPC_KEY_SERVICE_ADDRESS        | URI to the GRPC server running within the key service                                                                                       |               |
| APP_SERVER_PORT                 | The port your HTTP app server is running on                                                                                                 | 8080          |
| GRPC_SERVER_PORT                | The port your GRPC app server is running on                                                                                                 | 50051         |
| AUTH0_API_AUDIENCE              | The OATH2.0 audience provided by Auth0 for REST API endpoints that use client credentials                                                   |               |
| AUTH0_GRPC_AUDIENCE             | The OATH2.0 audience provided by Auth0 for GRPC endpoints that use client credentials                                                       |               |
| AUTH0_DOMAIN                    | Your Auth0 domain                                                                                                                           |               |
| AUTH0_CLIENT_CREDENTIALS_ID     | Your Auth0 client credentials id                                                                                                            |               |
| AUTH0_CLIENT_CREDENTIALS_SECRET | Your Auth0 client credentials secret                                                                                                        |               |
| MONGO_URI                       | The connection string to your mongodb instance/cluster.                                                                                     |               |
| MONGO_DB_NAME                   | MongoDB database                                                                                                                            |               |
| MONGO_CONNECTION_TIMEOUT_MS     | The duration in milliseconds for a MongoDB connection to time out                                                                           |               |
| RABBITMQ_URI                    | The URI to connect to RabbitMQ                                                                                                              |               |
| REDIS_ADDRESS                   | Your Redis `host:port` address                                                                                                              |               |
| REDIS_PASSWORD                  | Your Redis password                                                                                                                         |               |
| REDIS_DB                        | Your Redis database (use a number)                                                                                                          | 0             |

#### Building and Running your Go app
We have 2 ways of building a Go app, Makefile and the IDE Goland. Go does offer commands to build and run your app 
alongside additional plugins but for now, you can refer to the contents of the Makefile to see what the commands are.
In addition, we use [Google Wire](https://github.com/google/wire) for dependency injection, 
so we need to generate some code as a precompile step before we even compile and/or run. 

###### 1. Building and running with a Makefile
The `Makefile` you see in the `go` directory refers to the build automation commands available.

Run `make help` to see what targets are available.

Before you can even compile or run your Go code, you need to do a precompile step, generate code using 
[Google Wire](https://github.com/google/wire) for dependency injection. 

You can run `make wire` to generate the coded needed.

Now if you wish to run your app then without compiling, just run a command in the format
`make run service=${a subdirectory of cmd}`

For example: `make run service=userservice`

You can also just compile your app (for your platform). `make build service=${a subdirectory of cmd}`

To clean your build file, run `make clean`.

To run unit tests via the makefile, run `make test [optionally insert relative path from the go subdirectory]`

###### 2. Building and running with your IDE Goland
Running with your IDE is simpler and the preferred way of running your Go project. Just go into the `main.go` file
of in a microservice directory (e.g. `cmd/userservice`). When opened, a green arrow is placed on the `main` function.
Right click that arrow and select modify configuration. Then in the before launch tab, run the command in the following
format:
```shell
`run github.com/google/wire/cmd/wire gen github.com/obenkenobi/cypher-log/microservices/go/cmd/${subdirectory}/app`
```
${subdirectory} is just a stand in for any subdirectory of `cmd`.

For example when making a before launch plan for the user service, run the following:
```shell
run github.com/google/wire/cmd/wire gen github.com/obenkenobi/cypher-log/microservices/go/cmd/userservice/app
```

With the run configuration ready, you can now run your app.

To run unit tests on a directory or file, right-click the desired item, go to the run selection and pick the go test 
approach. Then you have a run configuration ready to run unit tests
