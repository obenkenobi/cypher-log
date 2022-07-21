package security

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type IdentityHolder interface {
	IsAnonymous() bool
	GetAuthorities() []string
	GetSubject() string
	ValidateAuthoritiesAny(requiredAuthorities []string) bool
	ValidateAuthoritiesAll(requiredAuthorities []string) bool
}

type identityHolderAuth0Impl struct {
	isAnonymous     bool
	validatedClaims *validator.ValidatedClaims
	customClaims    *Auth0CustomClaims
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

func (i identityHolderAuth0Impl) GetSubject() string {
	return i.validatedClaims.RegisteredClaims.Subject
}

func (i identityHolderAuth0Impl) ValidateAuthoritiesAny(requiredAuthorities []string) bool {
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

func (i identityHolderAuth0Impl) ValidateAuthoritiesAll(requiredAuthorities []string) bool {
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
	validatedClaims, ok := contextValue.(*validator.ValidatedClaims)
	isAnonymous := !ok
	if isAnonymous {
		isAnonymous = true
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
