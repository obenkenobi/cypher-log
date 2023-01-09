package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/businessrules"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/conf"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/repositories"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/businessobjects/userbos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/keydtos"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/encodingutils"
	"time"
)

type UserKeyService interface {
	UserKeyExists(ctx context.Context, userBo userbos.UserBo) (commondtos.ExistsDto, error)

	CreateUserKey(
		ctx context.Context,
		userBo userbos.UserBo,
		passwordDto keydtos.PasscodeCreateDto,
	) (commondtos.SuccessDto, error)

	NewKeySession(
		ctx context.Context,
		userBo userbos.UserBo,
		dto keydtos.PasscodeDto,
	) (commondtos.UKeySessionDto, error)

	GetKeyFromSession(ctx context.Context, sessionDto commondtos.UKeySessionDto) (keydtos.UserKeyDto, error)

	DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error)
}

type UserKeyServiceImpl struct {
	userKeyGeneratorRepository repositories.UserKeyGeneratorRepository
	userKeySessionRepository   repositories.UserKeySessionRepository
	userKeyBr                  businessrules.UserKeyBr
	appSecretService           AppSecretService
	errorService               sharedservices.ErrorService
	keyConf                    conf.KeyConf
}

func (u UserKeyServiceImpl) UserKeyExists(ctx context.Context, userBo userbos.UserBo) (commondtos.ExistsDto, error) {
	userFind, err := u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id)
	if err != nil {
		return commondtos.ExistsDto{}, err
	}
	return commondtos.NewExistsDto(userFind.IsPresent()), nil
}

func (u UserKeyServiceImpl) CreateUserKey(
	ctx context.Context,
	userBo userbos.UserBo,
	passcodeDto keydtos.PasscodeCreateDto,
) (commondtos.SuccessDto, error) {
	key, keyDerivationSalt, err := cipherutils.DeriveAESKeyFromPassword([]byte(passcodeDto.Passcode), nil)
	if err != nil {
		return commondtos.SuccessDto{}, err
	}
	keyHash, err := cipherutils.HashKeyBcrypt(key)
	if err != nil {
		return commondtos.SuccessDto{}, err
	}

	newUserKeyGen := models.UserKeyGenerator{
		UserId:            userBo.Id,
		KeyDerivationSalt: keyDerivationSalt,
		KeyHash:           keyHash,
		KeyVersion:        0,
	}

	if _, err := u.userKeyGeneratorRepository.Create(ctx, newUserKeyGen); err != nil {
		return commondtos.SuccessDto{}, err
	}
	return commondtos.NewSuccessTrue(), nil
}

func (u UserKeyServiceImpl) NewKeySession(
	ctx context.Context,
	userBo userbos.UserBo,
	dto keydtos.PasscodeDto,
) (commondtos.UKeySessionDto, error) {
	userKeyGen, err := u.getUserKeyGenerator(ctx, userBo)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}

	logger.Log.Debugf("Generating key from password")
	key, _, err := cipherutils.DeriveAESKeyFromPassword([]byte(dto.Passcode), userKeyGen.KeyDerivationSalt)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}

	if err := u.userKeyBr.ValidateKeyFromPassword(userKeyGen, key); err != nil {
		return commondtos.UKeySessionDto{}, err
	}

	proxyKey, err := cipherutils.GenerateRandomKeyAES()
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	appSecret, err := u.appSecretService.GetPrimaryAppSecret(ctx)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	tokenBytes, err := cipherutils.EncryptAES(appSecret.Key, proxyKey)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	tokenHash, err := cipherutils.HashWithSaltSHA256(tokenBytes)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	token := encodingutils.EncodeBase64String(tokenBytes)

	proxyKidUUID, err := uuid.NewRandom()
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	proxyKid := proxyKidUUID.String()
	if utils.StringIsBlank(proxyKid) {
		return commondtos.UKeySessionDto{}, errors.New("generated proxy KID is blank")
	}
	keyCipher, err := cipherutils.EncryptAES(proxyKey, key)
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	userIdCipher, err := cipherutils.EncryptAES(proxyKey, []byte(userBo.Id))
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}
	keyVersionCipher, err := cipherutils.EncryptAES(proxyKey, []byte(utils.Int64ToStr(userKeyGen.KeyVersion)))
	if err != nil {
		return commondtos.UKeySessionDto{}, err
	}

	keySessionModel := models.UserKeySession{
		KeyCipher:        keyCipher,
		TokenHash:        tokenHash,
		AppSecretKid:     appSecret.Kid,
		UserIdCipher:     userIdCipher,
		KeyVersionCipher: keyVersionCipher,
	}

	startTime := time.Now().UnixMilli()
	sessionDuration := u.keyConf.GetTokenSessionDuration()

	if _, err := u.userKeySessionRepository.Set(ctx, proxyKid, keySessionModel, sessionDuration); err != nil {
		return commondtos.UKeySessionDto{}, err
	}

	sessionDto := commondtos.UKeySessionDto{
		Token:         token,
		ProxyKid:      proxyKid,
		UserId:        userBo.Id,
		KeyVersion:    userKeyGen.KeyVersion,
		StartTime:     startTime,
		DurationMilli: sessionDuration.Milliseconds(),
	}
	return sessionDto, nil
}

func (u UserKeyServiceImpl) GetKeyFromSession(
	ctx context.Context,
	sessionDto commondtos.UKeySessionDto,
) (keydtos.UserKeyDto, error) {
	findStoredSession, err := u.userKeySessionRepository.Get(ctx, sessionDto.ProxyKid)
	if err != nil {
		return keydtos.UserKeyDto{}, err
	}
	session, sessionPresent := findStoredSession.Get()
	if !sessionPresent {
		ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession)
		return keydtos.UserKeyDto{}, apperrors.NewBadReqErrorFromRuleError(ruleErr)
	}
	tokenBytes, err := encodingutils.DecodeBase64String(sessionDto.Token)
	if err != nil {
		return keydtos.UserKeyDto{}, err
	}
	if err := u.userKeyBr.ValidateSessionTokenHash(session, tokenBytes); err != nil {
		return keydtos.UserKeyDto{}, err
	}
	appSecret, err := u.appSecretService.GetAppSecret(ctx, session.AppSecretKid)
	if err != nil {
		return keydtos.UserKeyDto{}, err
	}
	proxyKey, err := cipherutils.DecryptAES(appSecret.Key, tokenBytes)
	if err != nil {
		return keydtos.UserKeyDto{}, err
	}
	if err := u.userKeyBr.ValidateProxyKeyCiphersFromSession(
		ctx,
		proxyKey,
		sessionDto.UserId,
		sessionDto.KeyVersion,
		session,
	); err != nil {
		return keydtos.UserKeyDto{}, err
	}
	keyBytes, err := cipherutils.DecryptAES(proxyKey, session.KeyCipher)
	if err != nil {
		return keydtos.UserKeyDto{}, err
	}
	return keydtos.NewUserKeyDto(keyBytes, sessionDto.KeyVersion), nil
}

func (u UserKeyServiceImpl) getUserKeyGenerator(
	ctx context.Context,
	userBo userbos.UserBo,
) (models.UserKeyGenerator, error) {
	userKeyFind, err := u.userKeyGeneratorRepository.FindOneByUserId(ctx, userBo.Id)
	if err != nil {
		return models.UserKeyGenerator{}, err
	}
	if userKey, ok := userKeyFind.Get(); ok {
		return userKey, nil
	} else {
		ruleErr := u.errorService.RuleErrorFromCode(apperrors.ErrCodeReqResourcesNotFound)
		return models.UserKeyGenerator{}, apperrors.NewBadReqErrorFromRuleError(ruleErr)
	}
}

func (u UserKeyServiceImpl) DeleteByUserIdAndGetCount(ctx context.Context, userId string) (int64, error) {
	return u.userKeyGeneratorRepository.DeleteByUserIdAndGetCount(ctx, userId)
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
