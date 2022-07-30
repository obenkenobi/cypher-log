package errorservices

import (
	"fmt"
	"github.com/obenkenobi/cypher-log/services/go/pkg/apperrors"
)

type ErrorService interface {
	RuleErrorFromCode(code string, args ...any) apperrors.RuleError
}

type errorMessageServiceImpl struct {
	errorCodeToMsgMap map[string]string
}

func (e errorMessageServiceImpl) getMsgStrFromCode(code string, args ...any) string {
	if msg, ok := e.errorCodeToMsgMap[code]; ok {
		return fmt.Sprintf(msg, args...)
	}
	return code
}

func (e errorMessageServiceImpl) RuleErrorFromCode(code string, args ...any) apperrors.RuleError {
	return apperrors.RuleError{
		Code:    code,
		Message: e.getMsgStrFromCode(code, args...),
	}
}

func NewErrorService() ErrorService {
	errorCodeToMsgMap := map[string]string{
		apperrors.ErrCodeReqItemsNotFound:   "Requested item(s) not found",
		apperrors.ErrCodeCannotBindJson:     "Unable to bind json",
		apperrors.ErrCodeUserAlreadyCreated: "User already created",
		apperrors.ErrCodeUsernameTaken:      "Username is taken",
	}
	return &errorMessageServiceImpl{errorCodeToMsgMap: errorCodeToMsgMap}
}
