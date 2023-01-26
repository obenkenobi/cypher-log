package businessrules

import (
	"context"
	"github.com/obenkenobi/cypher-log/microservices/go/cmd/keyservice/models"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors/validationutils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/logger"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedservices"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/utils/cipherutils"
)

type UserKeyBr interface {
	ValidateSessionTokenHash(session models.UserKeySession, tokenBytes []byte) error
	ValidateKeyFromSession(userKeyGen models.UserKeyGenerator, key []byte) error
	ValidateKeyFromPassword(userKeyGen models.UserKeyGenerator, key []byte) error
	ValidateProxyKeyCiphersFromSession(
		ctx context.Context,
		proxyKey []byte,
		userId string,
		keyVersion int64,
		session models.UserKeySession,
	) error
}

type UserKeyBrImpl struct {
	errorService sharedservices.ErrorService
	userService  sharedservices.UserService
}

func (u UserKeyBrImpl) ValidateSessionTokenHash(session models.UserKeySession, tokenBytes []byte) error {
	var ruleErrs []apperrors.RuleError
	verified, err := cipherutils.VerifyHashWithSaltSHA256(session.TokenHash, tokenBytes)
	if err != nil {
		return err
	}
	if !verified {
		ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
	}
	return validationutils.MergeRuleErrors(ruleErrs)

}

func (u UserKeyBrImpl) ValidateKeyFromSession(userKeyGen models.UserKeyGenerator, key []byte) error {
	var ruleErrs []apperrors.RuleError

	verified, err := cipherutils.VerifyKeyHashBcrypt(userKeyGen.KeyHash, key)
	if err != nil {
		return err
	}
	if !verified {
		ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
	}

	return validationutils.MergeRuleErrors(ruleErrs)
}

func (u UserKeyBrImpl) ValidateProxyKeyCiphersFromSession(
	ctx context.Context,
	proxyKey []byte,
	userId string,
	keyVersion int64,
	session models.UserKeySession,
) error {
	var ruleErrs []apperrors.RuleError

	// Validate Proxy Key Ciphers
	savedUserIdBytes, err := cipherutils.DecryptAES(proxyKey, session.UserIdCipher)
	if err != nil {
		logger.Log.WithContext(ctx).WithError(err).Debug()
		return err
	}
	userIdInvalid := string(savedUserIdBytes) != userId
	savedKeyVersionBytes, err := cipherutils.DecryptAES(proxyKey, session.KeyVersionCipher)
	if err != nil {
		logger.Log.WithContext(ctx).WithError(err).Debug()
		return err
	}
	keyInvalid := string(savedKeyVersionBytes) != utils.Int64ToStr(keyVersion)

	if userIdInvalid || keyInvalid {
		ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
	}

	// Validate User Exists
	userExists, err := u.userService.UserExistsWithId(ctx, userId)
	if !userExists {
		ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeInvalidSession))
	}

	return validationutils.MergeRuleErrors(ruleErrs)
}

func (u UserKeyBrImpl) ValidateKeyFromPassword(userKeyGen models.UserKeyGenerator, key []byte) error {
	var ruleErrs []apperrors.RuleError
	verified, err := cipherutils.VerifyKeyHashBcrypt(userKeyGen.KeyHash, key)
	if err != nil {
		return err
	} else if !verified {
		ruleErrs = append(ruleErrs, u.errorService.RuleErrorFromCode(apperrors.ErrCodeIncorrectPasscode))
	}
	return validationutils.MergeRuleErrors(ruleErrs)
}

func NewUserKeyBrImpl(errorService sharedservices.ErrorService, userService sharedservices.UserService) *UserKeyBrImpl {
	return &UserKeyBrImpl{errorService: errorService, userService: userService}
}
