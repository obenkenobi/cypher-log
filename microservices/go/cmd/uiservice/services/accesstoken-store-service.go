package services

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/uiservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	"strings"
	"time"
)

type AccessTokenStoreService interface {
	StoreToken(ctx context.Context, userId string, tokenId string, accessToken string) error
	// GetToken gets a token if it exists or returns an empty string if there is no token
	GetToken(ctx context.Context, userId string, tokenId string) (string, error)
}

type AccessTokenStoreServiceImpl struct {
	sessionConf                 conf.SessionConf
	accessTokenHolderRepository repositories.AccessTokenHolderRepository
}

func (a AccessTokenStoreServiceImpl) StoreToken(
	ctx context.Context,
	userId string,
	tokenId string,
	accessToken string,
) error {
	tokenHolder := models.AccessTokenHolder{}
	if err := tokenHolder.SetAccessToken(accessToken, a.sessionConf.GetAccessTokenKey()); err != nil {
		return err
	}
	repoKey := a.createRepoKey(userId, tokenId)
	_, err := a.accessTokenHolderRepository.Set(ctx, repoKey, tokenHolder, time.Hour*24)
	return err
}

func (a AccessTokenStoreServiceImpl) GetToken(
	ctx context.Context,
	userId string,
	tokenId string,
) (string, error) {
	repoKey := a.createRepoKey(userId, tokenId)
	maybe, err := a.accessTokenHolderRepository.Get(ctx, repoKey)
	if err != nil {
		return "", err
	}
	holder, ok := maybe.Get()
	if !ok {
		return "", nil
	}
	return holder.GetAccessToken(a.sessionConf.GetAccessTokenKey())
}

func (a AccessTokenStoreServiceImpl) createRepoKey(userId string, tokenId string) string {
	return strings.Join([]string{userId, tokenId}, "/")
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
