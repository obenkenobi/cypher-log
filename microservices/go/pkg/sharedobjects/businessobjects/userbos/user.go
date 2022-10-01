package userbos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded/embeddeduser"
)

type UserBo struct {
	embeddeduser.BaseUserPublicDto
}
