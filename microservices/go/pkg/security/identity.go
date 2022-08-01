package security

import (
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/obenkenobi/cypher-log/services/go/pkg/utils"
	"strings"
)

// Identity contains information regarding who is accessing a resource, usually
// via an HTTP request.
type Identity interface {

	// IsAnonymous checks if the identity is unknown. An example would be if a user
	// who is not signed in is accessing a resource via an HTTP request.
	IsAnonymous() bool

	// IsUser checks if the identity is the user.
	IsUser() bool

	// IsSystemClient checks if the identity is a system client. A system client
	// refers to a machine is trying to access a resource as opposed to a person.
	IsSystemClient() bool

	// GetAuthorities gets the authorities needed to access a given resource. An
	// authority is a general concept of an array of strings that indicate what
	// resources an identity can access. Examples of what can be considered an
	// authorities include user roles, security groups, or even permissions from the
	// user's roles.
	GetAuthorities() []string

	// GetAuthId returns the id/subject of the identity from the authentication
	// provider. It is not the same as a user's ID within this application. The two
	// are linked and an authID can be used to get a User's ID provided the identity
	// is a User. This is to prevent vendor lock with an identity provider.
	GetAuthId() string

	// ContainsAnyAuthorities Checks if at identity contains at least one authority
	// provided in the slice.
	ContainsAnyAuthorities(authoritiesToCheck []string) bool

	// ContainsAllAuthorities Checks if at identity contains all the authorities
	// provided in the slice.
	ContainsAllAuthorities(requiredAuthorities []string) bool
}

type identityAuth0Impl struct {
	isAnonymous     bool
	validatedClaims *validator.ValidatedClaims
	customClaims    *Auth0CustomClaims
}

func (i identityAuth0Impl) IsUser() bool {
	return strings.Contains(i.GetAuthId(), "|")
}

func (i identityAuth0Impl) IsSystemClient() bool {
	return !i.isAnonymous && !i.IsUser()
}

func (i identityAuth0Impl) IsAnonymous() bool {
	return i.isAnonymous
}

func (i identityAuth0Impl) GetAuthorities() []string {
	if i.customClaims == nil {
		return []string{}
	}
	return i.customClaims.Permissions
}

func (i identityAuth0Impl) GetAuthId() string {
	return i.validatedClaims.RegisteredClaims.Subject
}

func (i identityAuth0Impl) ContainsAnyAuthorities(authoritiesToCheck []string) bool {
	if len(authoritiesToCheck) == 0 {
		return true
	}
	authoritiesToCheckSet := map[string]bool{}
	for _, authToCheck := range authoritiesToCheck {
		authoritiesToCheckSet[authToCheck] = true
	}
	for _, authority := range i.GetAuthorities() {
		if _, ok := authoritiesToCheckSet[authority]; ok {
			return true
		}
	}
	return false
}

func (i identityAuth0Impl) ContainsAllAuthorities(requiredAuthorities []string) bool {
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

// GetIdentityFromGinContext retrieves the Identity of whoever is accessing an
// HTTP request implemented with Gin.
func GetIdentityFromGinContext(c *gin.Context) Identity {
	contextValue := c.Request.Context().Value(jwtmiddleware.ContextKey{})
	isAnonymous := utils.StringIsBlank(c.GetHeader("Authorization"))
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
	return &identityAuth0Impl{
		isAnonymous:     isAnonymous,
		validatedClaims: validatedClaims,
		customClaims:    customClaims,
	}
}
