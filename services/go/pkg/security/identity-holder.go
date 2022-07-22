package security

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"strings"
)

type IdentityHolder interface {
	IsAnonymous() bool
	IsUser() bool
	IsSystemClient() bool
	GetAuthorities() []string
	GetIdFromProvider() string
	ContainsAnyAuthorities(requiredAuthorities []string) bool
	ContainsAllAuthorities(requiredAuthorities []string) bool
}

type identityHolderAuth0Impl struct {
	isAnonymous     bool
	validatedClaims *validator.ValidatedClaims
	customClaims    *Auth0CustomClaims
}

func (i identityHolderAuth0Impl) IsUser() bool {
	return strings.Contains(i.GetIdFromProvider(), "|")
}

func (i identityHolderAuth0Impl) IsSystemClient() bool {
	return !i.isAnonymous && !i.IsUser()
}

func (i identityHolderAuth0Impl) IsAnonymous() bool {
	return i.isAnonymous
}

func (i identityHolderAuth0Impl) GetAuthorities() []string {
	if i.customClaims == nil {
		return []string{}
	}
	return i.customClaims.Permissions
}

func (i identityHolderAuth0Impl) GetIdFromProvider() string {
	return i.validatedClaims.RegisteredClaims.Subject
}

func (i identityHolderAuth0Impl) ContainsAnyAuthorities(requiredAuthorities []string) bool {
	if len(requiredAuthorities) == 0 {
		return true
	}
	requiredAuthoritiesSet := map[string]bool{}
	for _, requiredAuthority := range requiredAuthorities {
		requiredAuthoritiesSet[requiredAuthority] = true
	}
	for _, authority := range i.GetAuthorities() {
		if _, ok := requiredAuthoritiesSet[authority]; ok {
			return true
		}
	}
	return false
}

func (i identityHolderAuth0Impl) ContainsAllAuthorities(requiredAuthorities []string) bool {
	if len(requiredAuthorities) == 0 {
		return true
	}
	authoritySet := map[string]bool{}
	for _, authority := range i.GetAuthorities() {
		authoritySet[authority] = true
	}
	for _, requiredAuthority := range requiredAuthorities {
		if _, ok := authoritySet[requiredAuthority]; !ok {
			return false
		}
	}
	return true
}

func NewIdentityHolderFromContext(c *gin.Context) IdentityHolder {
	contextValue := c.Request.Context().Value(jwtmiddleware.ContextKey{})
	isAnonymous := strings.TrimSpace(c.GetHeader("Authorization")) == ""
	validatedClaims, ok := contextValue.(*validator.ValidatedClaims)
	if !ok {
		validatedClaims = &validator.ValidatedClaims{
			RegisteredClaims: validator.RegisteredClaims{},
			CustomClaims:     defaultAuth0CustomClaims(),
		}
	}
	customClaims, ok := validatedClaims.CustomClaims.(*Auth0CustomClaims)
	if !ok {
		customClaims = defaultAuth0CustomClaims()
	}
	return &identityHolderAuth0Impl{
		isAnonymous:     isAnonymous,
		validatedClaims: validatedClaims,
		customClaims:    customClaims,
	}
}
