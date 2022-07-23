package apperrors

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

func GetErrorCodeToMsgMap() map[string]string {
	return errorCodeToMsgMap
}
