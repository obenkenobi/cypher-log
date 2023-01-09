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
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/wrappers/option"
	"time"
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
	) single.Single[commondtos.UKeySessionDto]

	GetKeyFromSession(
		ctx context.Context,
		sessionDto commondtos.UKeySessionDto,
	) single.Single[keydtos.UserKeyDto]

	DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64]
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
	userFindSrc := single.FromSupplierCached(func() (option.Maybe[models.UserKeyGenerator], error) {
		return u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id)
	})
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
	newKeySrc := single.FromSupplierCached(func() (derivedKey, error) {
		key, keyDerivationSalt, err := cipherutils.DeriveAESKeyFromPassword([]byte(passcodeDto.Passcode), nil)
		return derivedKey{key: key, keyDerivationSalt: keyDerivationSalt}, err
	})
	newKeyAndHashSrc := single.MapWithError(newKeySrc, func(dk derivedKey) (keyGeneration, error) {
		keyHash, err := cipherutils.HashKeyBcrypt(dk.key)
		return keyGeneration{keyHash: keyHash, keyDerivationSalt: dk.keyDerivationSalt}, err
	})
	newUserKeyGenSrc := single.Map(newKeyAndHashSrc, func(kg keyGeneration) models.UserKeyGenerator {
		return models.UserKeyGenerator{
			UserId:            userBo.Id,
			KeyDerivationSalt: kg.keyDerivationSalt,
			KeyHash:           kg.keyHash,
			KeyVersion:        0,
		}
	})
	userKeyGenSaveSrc := single.MapWithError(newUserKeyGenSrc,
		func(userKeyGen models.UserKeyGenerator) (models.UserKeyGenerator, error) {
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
) single.Single[commondtos.UKeySessionDto] {
	userKeyGenSrc := u.getUserKeyGenerator(ctx, userBo).ScheduleLazyAndCache(ctx)
	KeySrc := single.FlatMap(userKeyGenSrc,
		func(userKeyGen models.UserKeyGenerator) single.Single[[]byte] {
			newKeySrc := single.FromSupplierCached(func() ([]byte, error) {
				logger.Log.Debugf("Generating key from password")
				key, _, err := cipherutils.DeriveAESKeyFromPassword([]byte(dto.Passcode), userKeyGen.KeyDerivationSalt)
				return key, err
			})
			keyValidated := single.MapWithError(newKeySrc, func(key []byte) (any, error) {
				err := u.userKeyBr.ValidateKeyFromPassword(userKeyGen, key)
				return any(true), err
			})
			return single.FlatMap(keyValidated, func(_ any) single.Single[[]byte] {
				return newKeySrc
			})
		})
	proxyKeySrc := single.FromSupplierCached(cipherutils.GenerateRandomKeyAES)
	appSecretSrc := u.appSecretService.GetPrimaryAppSecret(ctx)
	return single.FlatMap(single.Zip4(proxyKeySrc, appSecretSrc, KeySrc, userKeyGenSrc),
		func(t tuple.T4[[]byte, bos.AppSecretBo, []byte, models.UserKeyGenerator]) single.Single[commondtos.UKeySessionDto] {
			proxyKey, appSecret, key, userKeyGen := t.V1, t.V2, t.V3, t.V4
			tokenBytesSrc := single.FromSupplierCached(func() ([]byte, error) {
				return cipherutils.EncryptAES(appSecret.Key, proxyKey)
			})
			tokenHashSrc := single.MapWithError(tokenBytesSrc, cipherutils.HashWithSaltSHA256)
			tokenSrc := single.Map(tokenBytesSrc, encodingutils.EncodeBase64String)
			proxyKidSrc := single.MapWithError(
				single.FromSupplierCached(uuid.NewRandom),
				func(uuid uuid.UUID) (string, error) {
					proxyKid := uuid.String()
					if utils.StringIsBlank(proxyKid) {
						return proxyKid, errors.New("generated proxy KID is blank")
					}
					return proxyKid, nil
				},
			)
			keyCipherSrc := single.FromSupplierCached(func() ([]byte, error) {
				return cipherutils.EncryptAES(proxyKey, key)
			})
			userIdCipherSrc := single.FromSupplierCached(func() ([]byte, error) {
				return cipherutils.EncryptAES(proxyKey, []byte(userBo.Id))
			})
			keyVersionCipherSrc := single.FromSupplierCached(func() ([]byte, error) {
				return cipherutils.EncryptAES(proxyKey, []byte(utils.Int64ToStr(userKeyGen.KeyVersion)))
			})
			return single.FlatMap(
				single.Zip6(tokenHashSrc, tokenSrc, proxyKidSrc, keyCipherSrc, userIdCipherSrc, keyVersionCipherSrc),
				func(t tuple.T6[[]byte, string, string, []byte, []byte, []byte]) single.Single[commondtos.UKeySessionDto] {
					tokenHash, token, proxyKid := t.V1, t.V2, t.V3
					keyCipher, userIdCipher, keyVersionCipher := t.V4, t.V5, t.V6
					keySessionModel := models.UserKeySession{
						KeyCipher:        keyCipher,
						TokenHash:        tokenHash,
						AppSecretKid:     appSecret.Kid,
						UserIdCipher:     userIdCipher,
						KeyVersionCipher: keyVersionCipher,
					}
					startTime := time.Now().UnixMilli()
					sessionDuration := u.keyConf.GetTokenSessionDuration()
					savedKeySession := u.userKeySessionRepository.Set(ctx, proxyKid, keySessionModel, sessionDuration)
					return single.Map(savedKeySession, func(_ models.UserKeySession) commondtos.UKeySessionDto {
						return commondtos.UKeySessionDto{
							Token:         token,
							ProxyKid:      proxyKid,
							UserId:        userBo.Id,
							KeyVersion:    userKeyGen.KeyVersion,
							StartTime:     startTime,
							DurationMilli: sessionDuration.Milliseconds(),
						}
					})
				},
			)
		},
	)
}

func (u UserKeyServiceImpl) GetKeyFromSession(
	ctx context.Context,
	sessionDto commondtos.UKeySessionDto,
) single.Single[keydtos.UserKeyDto] {
	storedSessionSrc := single.FlatMap(u.userKeySessionRepository.Get(ctx, sessionDto.ProxyKid),
		func(maybe option.Maybe[models.UserKeySession]) single.Single[models.UserKeySession] {
			return option.Map(maybe, single.Just[models.UserKeySession]).
				OrElseGet(func() single.Single[models.UserKeySession] {
					ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession)
					return single.Error[models.UserKeySession](apperrors.NewBadReqErrorFromRuleError(ruleErr))
				})
		},
	).ScheduleLazyAndCache(ctx)

	tokenBytesSrc := single.MapWithError(single.Just(sessionDto.Token), encodingutils.DecodeBase64String)
	tokenHashVerifiedSrc := single.MapWithError(single.Zip2(storedSessionSrc, tokenBytesSrc),
		func(t tuple.T2[models.UserKeySession, []byte]) (any, error) {
			session, tokenBytes := t.V1, t.V2
			return any(true), u.userKeyBr.ValidateSessionTokenHash(session, tokenBytes)
		},
	)
	appSecretSrc := single.FlatMap(single.Zip2(tokenHashVerifiedSrc, storedSessionSrc),
		func(t tuple.T2[any, models.UserKeySession]) single.Single[bos.AppSecretBo] {
			session := t.V2
			return u.appSecretService.GetAppSecret(ctx, session.AppSecretKid)
		})
	proxyKeySrc := single.MapWithError(single.Zip2(appSecretSrc, tokenBytesSrc),
		func(t tuple.T2[bos.AppSecretBo, []byte]) ([]byte, error) {
			appSecret, tokenBytes := t.V1, t.V2
			return cipherutils.DecryptAES(appSecret.Key, tokenBytes)
		},
	)
	keyBytesSrc := single.FlatMap(single.Zip2(storedSessionSrc, proxyKeySrc),
		func(t tuple.T2[models.UserKeySession, []byte]) single.Single[[]byte] {
			session, proxyKey := t.V1, t.V2
			err := u.userKeyBr.ValidateProxyKeyCiphersFromSession(
				ctx,
				proxyKey,
				sessionDto.UserId,
				sessionDto.KeyVersion,
				session,
			)
			if err != nil {
				return single.Error[[]byte](err)
			}
			return single.FromSupplierCached(func() ([]byte, error) {
				return cipherutils.DecryptAES(proxyKey, session.KeyCipher)
			})
		},
	)
	return single.Map(keyBytesSrc, func(keyBytes []byte) keydtos.UserKeyDto {
		return keydtos.NewUserKeyDto(keyBytes, sessionDto.KeyVersion)
	})
}

func (u UserKeyServiceImpl) getUserKeyGenerator(
	ctx context.Context,
	userBo userbos.UserBo,
) single.Single[models.UserKeyGenerator] {
	userKeyFind := single.FromSupplierCached(func() (option.Maybe[models.UserKeyGenerator], error) {
		return u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id)
	})
	return single.FlatMap(userKeyFind,
		func(maybe option.Maybe[models.UserKeyGenerator]) single.Single[models.UserKeyGenerator] {
			return option.Map(maybe, single.Just[models.UserKeyGenerator]).
				OrElseGet(func() single.Single[models.UserKeyGenerator] {
					ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
					return single.Error[models.UserKeyGenerator](apperrors.NewBadReqErrorFromRuleError(ruleErr))
				})
		},
	)
}

func (u UserKeyServiceImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) single.Single[int64] {
	return single.FromSupplierCached(func() (int64, error) {
		return u.userKeyGeneratorRepository.DeleteByUserIdAndGetCount(ctx, userId)
	})
}

func NewUserKeyServiceImpl(
	userKeyGeneratorRepository repositories.UserKeyGeneratorRepository,
	appSecretService AppSecretService,
	errorService sharedservices.ErrorService,
	userKeySessionRepository repositories.UserKeySessionRepository,
	userKeyBr businessrules.UserKeyBr,
	keyConf conf.KeyConf,
) *UserKeyServiceImpl {
	return &UserKeyServiceImpl{
		userKeyGeneratorRepository: userKeyGeneratorRepository,
		appSecretService:           appSecretService,
		errorService:               errorService,
		userKeySessionRepository:   userKeySessionRepository,
		userKeyBr:                  userKeyBr,
		keyConf:                    keyConf,
	}
}
