package embeddeduser

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded"

type BaseUserCommon struct {
	UserName    string `json:"userName" binding:"required"`
	DisplayName string `json:"displayName" binding:"required"`
}

type BaseUserAuthId struct {
	AuthId string `json:"authId"`
}

type BaseUserPublicDto struct {
	BaseUserCommon
	embedded.BaseCRUDObject
}
