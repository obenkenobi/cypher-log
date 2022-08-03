package grpcserveroptions

import "google.golang.org/grpc"

type ServerOptionCreator interface {
	CreateServerOption() grpc.ServerOption
}
