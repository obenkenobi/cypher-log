package middlewares

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionMiddleware interface {
	SessionHandler() gin.HandlerFunc
}

type SessionMiddlewareImpl struct {
}

func (s SessionMiddlewareImpl) SessionHandler() gin.HandlerFunc {
	//Todo: get a secret via env variable
	gob.Register(map[string]interface{}{})
	store := cookie.NewStore([]byte("secret"))
	return sessions.Sessions("auth-session", store)
}

func NewSessionMiddlewareImpl() *SessionMiddlewareImpl {
	return &SessionMiddlewareImpl{}
}
