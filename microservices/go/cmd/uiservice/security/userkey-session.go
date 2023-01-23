package security

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/objects/dtos/commondtos"
)

func WriteUKeySessionDtoToSession(c *gin.Context, sessionDto commondtos.UKeySessionDto) {
	session := sessions.Default(c)
	session.Set(UKeySessionKey, sessionDto)
}

func ReadUKeySessionDtoFromSession(c *gin.Context) (commondtos.UKeySessionDto, bool) {
	session := sessions.Default(c)
	sessDto, ok := session.Get(UKeySessionKey).(commondtos.UKeySessionDto)
	return sessDto, ok
}
