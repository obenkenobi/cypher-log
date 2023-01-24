package middlewares

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/microservices/go/pkg/conf"
	csrf "github.com/utrack/gin-csrf"
	"net/http"
)

type SessionMiddleware interface {
	// SessionHandler should be added first before any other middleware that uses sessions
	SessionHandler() gin.HandlerFunc
	CsrfHandler() gin.HandlerFunc
}

type SessionMiddlewareImpl struct {
	sessionConf conf.SessionConf
}

func init() {
	gob.Register(map[string]any{})
}

func (s SessionMiddlewareImpl) SessionHandler() gin.HandlerFunc {
	secret := s.sessionConf.GetSessionStoreSecret()
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
	return sessions.Sessions("session", store)
}

func (s SessionMiddlewareImpl) CsrfHandler() gin.HandlerFunc {
	return csrf.Middleware(csrf.Options{
		Secret: s.sessionConf.GetCSRFSecret(),
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	})
}

func NewSessionMiddlewareImpl(sessionConf conf.SessionConf) *SessionMiddlewareImpl {
	return &SessionMiddlewareImpl{sessionConf: sessionConf}
}
