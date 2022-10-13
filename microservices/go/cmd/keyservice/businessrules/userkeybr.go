package businessrules

import (
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/reactive/single"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
)

type UserKeyBr interface {
	ValidateSessionTokenHash(session models.UserKeySession, tokenBytes []byte) single.Single[[]apperrors.RuleError]
	ValidateKeyFromSession(userKeyGen models.UserKeyGenerator, key []byte) single.Single[[]apperrors.RuleError]
	ValidateKeyFromPassword(userKeyGen models.UserKeyGenerator, key []byte) single.Single[[]apperrors.RuleError]
}

type UserKeyBrImpl struct {
	errorService sharedservices.ErrorService
}

func (u UserKeyBrImpl) ValidateSessionTokenHash(
	session models.UserKeySession,
	tokenBytes []byte,
) single.Single[[]apperrors.RuleError] {
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

func (u UserKeyBrImpl) ValidateKeyFromSession(
	userKeyGen models.UserKeyGenerator,
	key []byte,
) single.Single[[]apperrors.RuleError] {
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

func (u UserKeyBrImpl) ValidateKeyFromPassword(
	userKeyGen models.UserKeyGenerator,
	key []byte,
) single.Single[[]apperrors.RuleError] {
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

func NewUserKeyBrImpl(errorService sharedservices.ErrorService) *UserKeyBrImpl {
	return &UserKeyBrImpl{errorService: errorService}
}
