echo 'Helper script to generate to generate go grpc files.
Pass an argument which will be the name of the proto file without the .proto line ending'
PROTO_DIR="$1pb"
protoc --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative "$PROTO_DIR/$1.proto"
