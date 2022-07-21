package security

import "context"

// Auth0CustomClaims contains custom data we want from the token.
type Auth0CustomClaims struct {
	Scope       string   `json:"scope"`
	Permissions []string `json:"permissions"`
}

// Validate does nothing for this example, but we need
// it to satisfy validator.CustomClaims interface.
func (c Auth0CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func defaultAuth0CustomClaims() *Auth0CustomClaims {
	return &Auth0CustomClaims{Scope: "", Permissions: []string{}}
}
