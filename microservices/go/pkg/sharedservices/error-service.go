package sharedservices

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/apperrors"
)

type ErrorService interface {
	RuleErrorFromCode(code string, args ...any) apperrors.RuleError
}

type ErrorServiceImpl struct {
	errorCodeToMsgMap map[string]string
}

func (e ErrorServiceImpl) getMsgStrFromCode(code string, args ...any) string {
	if msg, ok := e.errorCodeToMsgMap[code]; ok {
		return fmt.Sprintf(msg, args...)
	}
	return code
}

func (e ErrorServiceImpl) RuleErrorFromCode(code string, args ...any) apperrors.RuleError {
	return apperrors.RuleError{
		Code:    code,
		Message: e.getMsgStrFromCode(code, args...),
	}
}

func NewErrorServiceImpl() *ErrorServiceImpl {
	errorCodeToMsgMap := map[string]string{
		apperrors.ErrCodeReqResourcesNotFound:   "Requested resources not found",
		apperrors.ErrCodeCannotBindJson:         "Unable to bind json",
		apperrors.ErrCodeResourceAlreadyCreated: "Resource already created",
		apperrors.ErrCodeUsernameTaken:          "Username is taken",
		apperrors.ErrCodeUserRequireFail:        "User is not found or incomplete",
		apperrors.ErrCodeIncorrectPasscode:      "Incorrect passcode",
		apperrors.ErrCodeInvalidSession:         "Invalid session",
	}
	return &ErrorServiceImpl{errorCodeToMsgMap: errorCodeToMsgMap}
}
