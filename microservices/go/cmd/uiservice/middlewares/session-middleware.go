package middlewares

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
)

type SessionMiddleware interface {
	SessionHandler() gin.HandlerFunc
}

type SessionMiddlewareImpl struct {
	sessionConf conf.SessionConf
}

func (s SessionMiddlewareImpl) SessionHandler() gin.HandlerFunc {
	// Todo: configure redis as the sessions store
	gob.Register(map[string]interface{}{})
	secret := s.sessionConf.GetSessionStoreSecret()
	store := cookie.NewStore([]byte(secret))
	return sessions.Sessions("auth-session", store)
}

func NewSessionMiddlewareImpl(sessionConf conf.SessionConf) *SessionMiddlewareImpl {
	return &SessionMiddlewareImpl{sessionConf: sessionConf}
}
