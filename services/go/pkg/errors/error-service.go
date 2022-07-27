package errors

import (
	"fmt"
)

type ErrorService interface {
	RuleErrorFromCode(code string, args ...any) RuleError
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

func (e errorMessageServiceImpl) RuleErrorFromCode(code string, args ...any) RuleError {
	return RuleError{
		Code:    code,
		Message: e.getMsgStrFromCode(code, args...),
	}
}

func NewErrorServiceImpl() ErrorService {
	errorCodeToMsgMap := map[string]string{
		ErrCodeReqItemsNotFound:   "Requested item(s) not found",
		ErrCodeCannotBindJson:     "Unable to bind json",
		ErrCodeUserAlreadyCreated: "User already created",
		ErrCodeUsernameTaken:      "Username is taken",
	}
	return &errorMessageServiceImpl{errorCodeToMsgMap: errorCodeToMsgMap}
}
