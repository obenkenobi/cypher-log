# Cypher Log
A passcode based encrypted note taking app.

## Getting Started

### Technology Requirements
- Docker
- Go 1.18
- Node version 16.14.0
- Yarn
- Any Bash / Shell / ZSH terminal
- [Protocol Buffer compiler](https://developers.google.com/protocol-buffers), `protoc`, [version 3](https://developers.google.com/protocol-buffers/docs/proto3)
- VS Code or Goland

### Install Go Plugins
In order to support GRPC code generation, install the following plugins:
```shell
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
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
exists in `dev/`. Go to `dev/` and run `docker-compose up` to run your dependencies. Add the `-d` argument to run in 
detach mode. Alternatively you can use your IDE to run `docker-compose`.

### Todo: Add migrations and Go run instructions