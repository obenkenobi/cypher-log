package userbos

import (
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/sharedobjects/embedded/embeddeduser"
)

type UserBo struct {
	embeddeduser.BaseUserCommon
	embedded.BaseId
	embedded.BaseTimestamp
}
