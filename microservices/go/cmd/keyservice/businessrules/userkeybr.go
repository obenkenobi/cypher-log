package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
)

type UserKeyBr interface {
	ValidateSessionTokenHash(session models.UserKeySession, tokenBytes []byte) single.Single[any]
	ValidateKeyFromSession(userKeyGen models.UserKeyGenerator, key []byte) single.Single[any]
	ValidateKeyFromPassword(userKeyGen models.UserKeyGenerator, key []byte) single.Single[any]
	ValidateProxyKeyCiphersFromSession(
		ctx context.Context,
		proxyKey []byte,
		userId string,
		keyVersion int64,
		session models.UserKeySession,
	) single.Single[any]
}

type UserKeyBrImpl struct {
	errorService sharedservices.ErrorService
	userService  sharedservices.UserService
}

func (u UserKeyBrImpl) ValidateSessionTokenHash(session models.UserKeySession, tokenBytes []byte) single.Single[any] {
	hashCheckSrc := single.FromSupplierCached(func() (bool, error) {
		return cipherutils.VerifyHashWithSaltSHA256(session.TokenHash, tokenBytes)
	})
	tokenHashValidate := single.Map(hashCheckSrc, func(verified bool) []apperrors.RuleError {
		var ruleErrs []apperrors.RuleError
		if !verified {
			ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
		}
		return ruleErrs
	})
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(tokenHashValidate)

}

func (u UserKeyBrImpl) ValidateKeyFromSession(userKeyGen models.UserKeyGenerator, key []byte) single.Single[any] {
	verifiedKeyHashSrc := single.FromSupplierCached(func() (bool, error) {
		return cipherutils.VerifyKeyHashBcrypt(userKeyGen.KeyHash, key)
	})
	keyHashValidationSrc := single.Map(verifiedKeyHashSrc, func(verified bool) []apperrors.RuleError {
		var ruleErrs []apperrors.RuleError
		if !verified {
			ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
		}
		return ruleErrs
	})
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(keyHashValidationSrc)
}

func (u UserKeyBrImpl) ValidateProxyKeyCiphersFromSession(
	ctx context.Context,
	proxyKey []byte,
	userId string,
	keyVersion int64,
	session models.UserKeySession,
) single.Single[any] {
	validateProxyKeyCiphersSrc := single.FromSupplierCached(func() ([]apperrors.RuleError, error) {
		var ruleErrs []apperrors.RuleError

		savedUserIdBytes, err := cipherutils.DecryptAES(proxyKey, session.UserIdCipher)
		if err != nil {
			logger.Log.WithError(err).Debug()
			return ruleErrs, err
		}
		userIdInvalid := string(savedUserIdBytes) != userId

		savedKeyVersionBytes, err := cipherutils.DecryptAES(proxyKey, session.KeyVersionCipher)
		if err != nil {
			logger.Log.WithError(err).Debug()
			return ruleErrs, err
		}
		keyInvalid := string(savedKeyVersionBytes) != utils.Int64ToStr(keyVersion)

		if userIdInvalid || keyInvalid {
			ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
		}
		return ruleErrs, nil
	})
	validateUserExists := single.MapWithError(u.userService.UserExistsWithId(ctx, userId),
		func(exists bool) ([]apperrors.RuleError, error) {
			var ruleErrs []apperrors.RuleError
			if !exists {
				ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
			}
			return ruleErrs, nil
		})
	ruleErrs := validationutils.ConcatSinglesOfRuleErrs(validateProxyKeyCiphersSrc, validateUserExists)
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(ruleErrs)
}

func (u UserKeyBrImpl) ValidateKeyFromPassword(
	userKeyGen models.UserKeyGenerator,
	key []byte,
) single.Single[any] {
	verifiedKeyHashSrc := single.FromSupplierCached(func() (bool, error) {
		return cipherutils.VerifyKeyHashBcrypt(userKeyGen.KeyHash, key)
	})
	keyHashValidationSrc := single.Map(verifiedKeyHashSrc, func(verified bool) []apperrors.RuleError {
		var ruleErrs []apperrors.RuleError
		if !verified {
			ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeIncorrectPasscode))
		}
		return ruleErrs
	})
	return validationutils.PassRuleErrorsIfEmptyElsePassBadReqError(keyHashValidationSrc)
}

func NewUserKeyBrImpl(errorService sharedservices.ErrorService, userService sharedservices.UserService) *UserKeyBrImpl {
	return &UserKeyBrImpl{errorService: errorService, userService: userService}
}
