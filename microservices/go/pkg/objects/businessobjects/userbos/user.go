package userbos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/embedded/embeddeduser"
)

type UserBo struct {
	embeddeduser.BaseUserPublicDto
}
