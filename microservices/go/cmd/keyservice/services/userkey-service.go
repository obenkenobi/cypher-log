package services

import (
	"context"
	"errors"
	"github.com/barweiss/go-tuple"
	"github.com/google/uuid"
	bos "github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/businessobjects"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
)

type UserKeyService interface {
	UserKeyExists(ctx context.Context, userBo userbos.UserBo) single.Single[commondtos.ExistsDto]

	CreateUserKey(
		ctx context.Context,
		userBo userbos.UserBo,
		passwordDto keydtos.PasscodeCreateDto,
	) single.Single[commondtos.SuccessDto]

	NewKeySession(
		ctx context.Context,
		userBo userbos.UserBo,
		dto keydtos.PasscodeDto,
	) single.Single[keydtos.UserKeySessionDto]

	GetKeyFromSession(
		ctx context.Context,
		sessionDto keydtos.UserKeySessionDto,
	) single.Single[keydtos.UserKeyDto]
}

type UserKeyServiceImpl struct {
	userKeyGeneratorRepository repositories.UserKeyGeneratorRepository
	userKeySessionRepository   repositories.UserKeySessionRepository
	userKeyBr                  businessrules.UserKeyBr
	appSecretService           AppSecretService
	errorService               sharedservices.ErrorService
	keyConf                    conf.KeyConf
}

func (u UserKeyServiceImpl) UserKeyExists(
	ctx context.Context,
	userBo userbos.UserBo,
) single.Single[commondtos.ExistsDto] {
	userFindSrc := u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id)
	return single.Map(userFindSrc, func(maybe option.Maybe[models.UserKeyGenerator]) commondtos.ExistsDto {
		return commondtos.NewExistsDto(maybe.IsPresent())
	})
}

func (u UserKeyServiceImpl) CreateUserKey(
	ctx context.Context,
	userBo userbos.UserBo,
	passcodeDto keydtos.PasscodeCreateDto,
) single.Single[commondtos.SuccessDto] {
	type derivedKey struct {
		key               []byte
		keyDerivationSalt []byte
	}
	type keyGeneration struct {
		keyHash           []byte
		keyDerivationSalt []byte
	}
	newKeySrc := single.FromSupplier(func() (derivedKey, error) {
		key, keyDerivationSalt, err := cipherutils.DeriveAESKeyFromPassword([]byte(passcodeDto.Passcode), nil)
		return derivedKey{key: key, keyDerivationSalt: keyDerivationSalt}, err
	})
	newKeyAndHashSrc := single.MapWithError(newKeySrc, func(dk derivedKey) (keyGeneration, error) {
		keyHash, err := cipherutils.HashKeyBcrypt(dk.key)
		return keyGeneration{keyHash: keyHash, keyDerivationSalt: dk.keyDerivationSalt}, err
	})
	newUserKeyGenSrc := single.Map(newKeyAndHashSrc, func(kg keyGeneration) models.UserKeyGenerator {
		return models.UserKeyGenerator{UserId: userBo.Id, KeyDerivationSalt: kg.keyDerivationSalt, KeyHash: kg.keyHash}
	})
	userKeyGenSaveSrc := single.FlatMap(newUserKeyGenSrc,
		func(userKeyGen models.UserKeyGenerator) single.Single[models.UserKeyGenerator] {
			return u.userKeyGeneratorRepository.Create(ctx, userKeyGen)
		},
	)
	return single.Map(userKeyGenSaveSrc, func(_ models.UserKeyGenerator) commondtos.SuccessDto {
		return commondtos.NewSuccessTrue()
	})
}

func (u UserKeyServiceImpl) NewKeySession(
	ctx context.Context,
	userBo userbos.UserBo,
	dto keydtos.PasscodeDto,
) single.Single[keydtos.UserKeySessionDto] {
	userKeyGenSrc := u.getUserKeyGenerator(ctx, userBo)
	KeySrc := single.FlatMap(userKeyGenSrc,
		func(userKeyGen models.UserKeyGenerator) single.Single[[]byte] {
			newKeySrc := single.FromSupplier(func() ([]byte, error) {
				key, _, err := cipherutils.DeriveAESKeyFromPassword([]byte(dto.Passcode), userKeyGen.KeyDerivationSalt)
				return key, err
			})
			keyValidated := single.FlatMap(newKeySrc, func(key []byte) single.Single[[]apperrors.RuleError] {
				return u.userKeyBr.ValidateKeyFromPassword(userKeyGen, key)
			})
			return single.FlatMap(keyValidated, func(_ []apperrors.RuleError) single.Single[[]byte] {
				return newKeySrc
			})
		})
	proxyKeySrc := single.FromSupplier(cipherutils.GenerateRandomKeyAES)
	appSecretSrc := u.appSecretService.GetPrimaryAppSecret(ctx)
	return single.FlatMap(single.Zip3(proxyKeySrc, appSecretSrc, KeySrc),
		func(t tuple.T3[[]byte, bos.AppSecretBo, []byte]) single.Single[keydtos.UserKeySessionDto] {
			proxyKey, appSecret, key := t.V1, t.V2, t.V3
			tokenBytesSrc := single.FromSupplier(func() ([]byte, error) {
				return cipherutils.EncryptAES(appSecret.Key, proxyKey)
			})
			tokenHashSrc := single.MapWithError(tokenBytesSrc, cipherutils.HashWithSaltSHA256)
			tokenSrc := single.Map(tokenBytesSrc, encodingutils.EncodeBase64String)
			proxyKidSrc := single.MapWithError(
				single.FromSupplier(uuid.NewRandom),
				func(uuid uuid.UUID) (string, error) {
					proxyKid := uuid.String()
					if utils.StringIsBlank(proxyKid) {
						return proxyKid, errors.New("generated proxy KID is blank")
					}
					return proxyKid, nil
				},
			)
			encryptedKeySrc := single.FromSupplier(func() ([]byte, error) {
				return cipherutils.EncryptAES(proxyKey, key)
			})
			return single.FlatMap(single.Zip4(tokenHashSrc, tokenSrc, proxyKidSrc, encryptedKeySrc),
				func(t tuple.T4[[]byte, string, string, []byte]) single.Single[keydtos.UserKeySessionDto] {
					tokenHash, token, proxyKid, encryptedKey := t.V1, t.V2, t.V3, t.V4
					keySessionModel := models.UserKeySession{
						EncryptedKey: encryptedKey,
						TokenHash:    tokenHash,
						AppSecretKid: appSecret.Kid,
					}
					savedKeySession := u.userKeySessionRepository.Set(ctx, proxyKid,
						keySessionModel, u.keyConf.GetTokenSessionDuration())
					return single.Map(savedKeySession, func(_ models.UserKeySession) keydtos.UserKeySessionDto {
						return keydtos.UserKeySessionDto{Token: token, ProxyKid: proxyKid}
					})
				},
			)
		},
	)
}

func (u UserKeyServiceImpl) GetKeyFromSession(
	ctx context.Context,
	sessionDto keydtos.UserKeySessionDto,
) single.Single[keydtos.UserKeyDto] {
	storedSessionSrc := single.FlatMap(u.userKeySessionRepository.Get(ctx, sessionDto.ProxyKid),
		func(maybe option.Maybe[models.UserKeySession]) single.Single[models.UserKeySession] {
			return option.Map(maybe, single.Just[models.UserKeySession]).
				OrElseGet(func() single.Single[models.UserKeySession] {
					ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession)
					return single.Error[models.UserKeySession](apperrors.NewBadReqErrorFromRuleError(ruleErr))
				})
		},
	)
	tokenBytesSrc := single.FromSupplier(func() ([]byte, error) {
		return encodingutils.DecodeBase64String(sessionDto.Token)
	})
	tokenHashVerifiedSrc := single.FlatMap(single.Zip2(storedSessionSrc, tokenBytesSrc),
		func(t tuple.T2[models.UserKeySession, []byte]) single.Single[[]apperrors.RuleError] {
			session, tokenBytes := t.V1, t.V2
			return u.userKeyBr.ValidateSessionTokenHash(session, tokenBytes)
		},
	)
	appSecretSrc := single.FlatMap(single.Zip2(tokenHashVerifiedSrc, storedSessionSrc),
		func(t tuple.T2[[]apperrors.RuleError, models.UserKeySession]) single.Single[bos.AppSecretBo] {
			session := t.V2
			return u.appSecretService.GetAppSecret(ctx, session.AppSecretKid)
		})
	proxyKeySrc := single.FlatMap(appSecretSrc, func(appSecret bos.AppSecretBo) single.Single[[]byte] {
		proxyKeyCipherSrc := single.MapWithError(single.Just(sessionDto.Token), encodingutils.DecodeBase64String)
		return single.MapWithError(proxyKeyCipherSrc, func(proxyKeyCipher []byte) ([]byte, error) {
			return cipherutils.DecryptAES(appSecret.Key, proxyKeyCipher)
		})
	})
	keyBytesSrc := single.MapWithError(single.Zip2(storedSessionSrc, proxyKeySrc),
		func(t tuple.T2[models.UserKeySession, []byte]) ([]byte, error) {
			session, proxyKey := t.V1, t.V2
			return cipherutils.DecryptAES(proxyKey, session.EncryptedKey)
		},
	)
	return single.Map(keyBytesSrc, keydtos.NewUserKeyDto)
}

func (u UserKeyServiceImpl) getUserKeyGenerator(
	ctx context.Context,
	userBo userbos.UserBo,
) single.Single[models.UserKeyGenerator] {
	return single.FlatMap(u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id),
		func(maybe option.Maybe[models.UserKeyGenerator]) single.Single[models.UserKeyGenerator] {
			return option.Map(maybe, single.Just[models.UserKeyGenerator]).
				OrElseGet(func() single.Single[models.UserKeyGenerator] {
					ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
					return single.Error[models.UserKeyGenerator](apperrors.NewBadReqErrorFromRuleError(ruleErr))
				})
		},
	)
}

func NewUserKeyServiceImpl(
	userKeyGeneratorRepository repositories.UserKeyGeneratorRepository,
	appSecretService AppSecretService,
	errorService sharedservices.ErrorService,
	userKeySessionRepository repositories.UserKeySessionRepository,
	userKeyBr businessrules.UserKeyBr,
) *UserKeyServiceImpl {
	return &UserKeyServiceImpl{
		userKeyGeneratorRepository: userKeyGeneratorRepository,
		appSecretService:           appSecretService,
		errorService:               errorService,
		userKeySessionRepository:   userKeySessionRepository,
		userKeyBr:                  userKeyBr,
	}
}
