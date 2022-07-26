package grpcserveropts

import (
	"crypto/tls"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type CredentialsOptionCreator interface {
	CreateCredentialsOption() grpc.ServerOption
}

type CredentialsOptionCreatorImpl struct {
	tlsConf conf.TLSConf
}

func (c CredentialsOptionCreatorImpl) CreateCredentialsOption() grpc.ServerOption {
	return grpc.Creds(c.loadTLSCredentials())
}

func (c CredentialsOptionCreatorImpl) loadTLSCredentials() credentials.TransportCredentials {
	serverCert, err := tls.LoadX509KeyPair(c.tlsConf.ServerCertPath(), c.tlsConf.ServerKeyPath())
	if err != nil {
		logger.Log.Fatal("cannot load TLS credentials: ", err)
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config)
}

func NewCredentialsOptionCreatorImpl(tlsConf conf.TLSConf) *CredentialsOptionCreatorImpl {
	return &CredentialsOptionCreatorImpl{tlsConf: tlsConf}
}
