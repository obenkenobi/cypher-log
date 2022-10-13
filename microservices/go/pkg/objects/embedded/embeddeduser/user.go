package embeddeduser

import "github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded"

type BaseUserCommon struct {
	UserName    string `json:"userName" binding:"required,alphanumunicode,min=4,max=255"`
	DisplayName string `json:"displayName" binding:"required,min=4,max=255"`
}

type BaseUserAuthId struct {
	AuthId string `json:"authId"`
}

type BaseUserPublicDto struct {
	BaseUserCommon
	embedded.BaseCRUDObject
}
