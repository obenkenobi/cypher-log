package security

import "context"

// Auth0CustomClaims contains custom claims from an Auth0 token
type Auth0CustomClaims struct {
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c Auth0CustomClaims) Validate(ctx context.Context) error {
	return nil
}

// defaultAuth0CustomClaims are default claims for the custom claims of an auth0 token
func defaultAuth0CustomClaims() *Auth0CustomClaims {
	return &Auth0CustomClaims{Scope: "", Permissions: []string{}}
}
