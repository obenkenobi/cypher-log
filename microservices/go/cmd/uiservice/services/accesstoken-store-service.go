package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"time"
)

type AccessTokenStoreService interface {
	StoreToken(ctx context.Context, tokenId string, accessToken string) error
	// GetToken gets a token if it exists or returns an empty string if there is no token
	GetToken(ctx context.Context, tokenId string) (string, error)
	DeleteToken(ctx context.Context, tokenId string) error
}

type AccessTokenStoreServiceImpl struct {
	sessionConf                 conf.SessionConf
	accessTokenHolderRepository repositories.AccessTokenHolderRepository
}

func (a AccessTokenStoreServiceImpl) StoreToken(ctx context.Context, tokenId string, accessToken string) error {
	tokenHolder := models.AccessTokenHolder{}
	if err := tokenHolder.SetAccessToken(accessToken, a.sessionConf.GetAccessTokenKey()); err != nil {
		return err
	}
	_, err := a.accessTokenHolderRepository.Set(ctx, tokenId, tokenHolder, time.Hour*24)
	return err
}

func (a AccessTokenStoreServiceImpl) GetToken(ctx context.Context, tokenId string) (string, error) {
	maybe, err := a.accessTokenHolderRepository.Get(ctx, tokenId)
	if err != nil {
		return "", err
	}
	holder, ok := maybe.Get()
	if !ok {
		return "", nil
	}
	return holder.GetAccessToken(a.sessionConf.GetAccessTokenKey())
}

func (a AccessTokenStoreServiceImpl) DeleteToken(ctx context.Context, tokenId string) error {
	return a.accessTokenHolderRepository.Del(ctx, tokenId)
}

func NewAccessTokenStoreServiceImpl(
	sessionConf conf.SessionConf,
	accessTokenHolderRepository repositories.AccessTokenHolderRepository,
) *AccessTokenStoreServiceImpl {
	return &AccessTokenStoreServiceImpl{
		sessionConf:                 sessionConf,
		accessTokenHolderRepository: accessTokenHolderRepository,
	}
}
