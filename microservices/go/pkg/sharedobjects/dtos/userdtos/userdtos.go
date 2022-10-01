package userdtos

import (
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

type UserChangeAction bool

const (
	UserSave   UserChangeAction = true
	UserDelete UserChangeAction = false
)

type UserChangeEventDto struct {
	embeddeduser.BaseUserPublicDto
	embeddeduser.BaseUserAuthId
	Action UserChangeAction `json:"action"`
}

type UserSaveDto struct {
	embeddeduser.BaseUserCommon
}
