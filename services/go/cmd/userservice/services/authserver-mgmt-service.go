package services

import (
	"github.com/obenkenobi/cypher-log/services/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/services/go/pkg/extensions/streamx/single"
	"gopkg.in/auth0.v5/management"
)

type AuthServerMgmtService interface {
	DeleteUser(authId string) single.Single[bool]
	DeleteAsync(authId string) single.Single[bool]
}

type AuthServerMgmtServiceImpl struct {
	auth0ClientCredentialsConf authconf.Auth0ClientCredentialsConf
}

func (a AuthServerMgmtServiceImpl) DeleteUser(authId string) single.Single[bool] {
	return single.FromSupplier(func() (bool, error) { return a.runDeleteUser(authId) })
}

func (a AuthServerMgmtServiceImpl) DeleteAsync(authId string) single.Single[bool] {
	return single.FromSupplier(func() (bool, error) { return a.runDeleteUser(authId) })
}

func (a AuthServerMgmtServiceImpl) runDeleteUser(authId string) (bool, error) {
	m, err := management.New(a.auth0ClientCredentialsConf.GetDomain(), management.WithClientCredentials(
		a.auth0ClientCredentialsConf.GetClientId(),
		a.auth0ClientCredentialsConf.GetClientSecret(),
	))
	if err != nil {
		return false, err
	}
	// login.auth0.com/api/v2
	if err := m.User.Delete(authId); err != nil {
		return false, err
	}
	return true, nil
}

func NewAuthServerMgmtService(auth0ClientCredentialsConf authconf.Auth0ClientCredentialsConf) AuthServerMgmtService {
	return &AuthServerMgmtServiceImpl{auth0ClientCredentialsConf: auth0ClientCredentialsConf}
}
