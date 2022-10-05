package services

import (
	"context"
	"github.com/barweiss/go-tuple"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
)

type UserKeyService interface {
	CreateUserKey(
		ctx context.Context,
		userBo userbos.UserBo,
		passwordDto keydtos.PasscodeCreateDto,
	) single.Single[commondtos.SuccessDto]
}

type UserKeyServiceImpl struct {
	userKeyRepository repositories.UserKeyRepository
	appSecretService  AppSecretService
}

func (u UserKeyServiceImpl) CreateUserKey(
	ctx context.Context,
	userBo userbos.UserBo,
	passwordDto keydtos.PasscodeCreateDto,
) single.Single[commondtos.SuccessDto] {
	type tKey []byte
	type tKeyDerivationSalt []byte
	type tKeyHash []byte
	newKeySrc := single.FromSupplier(func() (tuple.T2[tKey, tKeyDerivationSalt], error) {
		key, keyDerivationSalt, err := cipherutils.DeriveKey([]byte(passwordDto.Passcode), nil)
		return tuple.New2(tKey(key), tKeyDerivationSalt(keyDerivationSalt)), err
	})
	newKeyAndHashSrc := single.MapWithError(
		newKeySrc,
		func(t tuple.T2[tKey, tKeyDerivationSalt]) (tuple.T2[tKeyDerivationSalt, tKeyHash], error) {
			key, keyDerivationSalt := t.V1, t.V2
			keyHash, err := cipherutils.HashKey(key)
			return tuple.New2(keyDerivationSalt, tKeyHash(keyHash)), err
		},
	)
	newUserKeySrc := single.Map(newKeyAndHashSrc, func(t tuple.T2[tKeyDerivationSalt, tKeyHash]) models.UserKey {
		keyDerivationSalt, keyHash := t.V1, t.V2
		return models.UserKey{UserId: userBo.Id, KeyDerivationSalt: keyDerivationSalt, KeyHash: keyHash}

	})

	userKeySaveSrc := single.FlatMap(newUserKeySrc, func(userKey models.UserKey) single.Single[models.UserKey] {
		return u.userKeyRepository.Create(ctx, userKey)
	})
	return single.Map(userKeySaveSrc, func(_ models.UserKey) commondtos.SuccessDto {
		return commondtos.NewSuccessTrue()
	})
}

func NewUserKeyServiceImpl(
	userKeyRepository repositories.UserKeyRepository,
	appSecretService AppSecretService,
) *UserKeyServiceImpl {
	return &UserKeyServiceImpl{userKeyRepository: userKeyRepository, appSecretService: appSecretService}
}
