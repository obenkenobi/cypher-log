package errormgmt

import "github.com/obenkenobi/cypher-log/services/go/pkg/dtos/errordtos"

const ErrCodeReqItemsNotFound = "ReqItemsNotFound"
const ErrCodeCannotBindJson = "CannotBindJson"
const ErrCodeUserAlreadyCreated = "UserAlreadyCreated"
const ErrCodeUsernameTaken = "UsernameTaken"

var errorCodeToMsgMap = map[string]string{
	ErrCodeReqItemsNotFound:   "Requested item(s) not found",
	ErrCodeCannotBindJson:     "Unable to bind json",
	ErrCodeUserAlreadyCreated: "User already created",
	ErrCodeUsernameTaken:      "Username is taken",
}

func getMsgFromCode(code string) string {
	if msg, ok := errorCodeToMsgMap[code]; ok {
		return msg
	}
	return code
}

func GetAppErrorDtosFromCodes(codes []string) []errordtos.AppErrorDto {
	codeToMsgMap := map[string]string{}
	for _, code := range codes {
		codeToMsgMap[code] = getMsgFromCode(code)
	}
	var appErrors []errordtos.AppErrorDto
	for code, msg := range codeToMsgMap {
		appErrors = append(appErrors, errordtos.AppErrorDto{Code: code, Message: msg})
	}
	return appErrors
}
