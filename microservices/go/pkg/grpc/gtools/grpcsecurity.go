package gtools

import (
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func LoadTLSCredentialsOption(certPath string, overrideServerName bool) (grpc.DialOption, error) {
	// Load certificate of the CA who signed server's certificate
	serverNameOverride := ""
	if overrideServerName {
		serverNameOverride = "x.cypherlog.com"
	}
	creds, err := credentials.NewClientTLSFromFile(certPath, serverNameOverride)
	if err != nil {
		return nil, err
	}
	return grpc.WithTransportCredentials(creds), nil
}

func OathAccessOption(token oauth2.Token) grpc.DialOption {
	return grpc.WithPerRPCCredentials(oauth.NewOauthAccess(&token))
}
