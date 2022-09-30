package userdtos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded/embeddeduser"
)

type UserIdentityDto struct {
	UserReadDto
	embeddeduser.BaseUserAuthId
	Authorities []string `json:"authorities"`
}

type UserReadDto struct {
	embeddeduser.BaseUserPublicDto
	Exists bool `json:"exists"`
}

type DistUserSaveDto struct {
	embeddeduser.BaseUserPublicDto
	embeddeduser.BaseUserAuthId
}

type DistUserDeleteDto struct {
	embedded.BaseId
}

type UserSaveDto struct {
	embeddeduser.BaseUserCommon
}
