package services

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf/authconf"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"gopkg.in/auth0.v5/management"
)

type AuthServerMgmtService interface {
	DeleteUser(authId string) single.Single[bool]
	DeleteAsync(authId string) single.Single[bool]
}

type AuthServerMgmtServiceImpl struct {
	auth0SecurityConf authconf.Auth0SecurityConf
}

func (a AuthServerMgmtServiceImpl) DeleteUser(authId string) single.Single[bool] {
	return single.FromSupplierCached(func() (bool, error) { return a.runDeleteUser(authId) })
}

func (a AuthServerMgmtServiceImpl) DeleteAsync(authId string) single.Single[bool] {
	return single.FromSupplierCached(func() (bool, error) { return a.runDeleteUser(authId) })
}

func (a AuthServerMgmtServiceImpl) runDeleteUser(authId string) (bool, error) {
	m, err := management.New(a.auth0SecurityConf.GetDomain(), management.WithClientCredentials(
		a.auth0SecurityConf.GetClientCredentialsId(),
		a.auth0SecurityConf.GetClientCredentialsSecret(),
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

func NewAuthServerMgmtServiceImpl(auth0ClientCredentialsConf authconf.Auth0SecurityConf) *AuthServerMgmtServiceImpl {
	return &AuthServerMgmtServiceImpl{auth0SecurityConf: auth0ClientCredentialsConf}
}
