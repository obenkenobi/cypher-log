# Cypher Log
A passcode based encrypted note-taking app.

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
Here is a list of all the environment variables used in the mongo migrator:
- `MONGO_URI` Is the connection string to your mongodb instance/cluster.
- `MIGRATION_DIRECTORY` Is the migration directory under the `src` directory.
- `MONGO_DB_NAME` Is the database name, it is recommended that it's value is the same as `MIGRATION_DIRECTORY`.
- `IGNORE_ENV_FILE` If set to `true`, your .env file is ignored. This will not work if you set this in the `.env` file.

## Todo: add go build and run instructions.