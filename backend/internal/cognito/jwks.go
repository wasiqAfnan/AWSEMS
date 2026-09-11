package cognito

import (
	"fmt"
	"strings"

	"awsems/internal/config"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	JWKS     *keyfunc.JWKS
	Issuer   string
	ClientID string
}

func NewJWTVerifier(cfg *config.Config) (*JWTVerifier, error) {
	jwksURL := strings.TrimRight(
		cfg.CognitoOpenIDConfigURL,
		"/",
	) + "/.well-known/jwks.json"

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to load Cognito JWKS: %w", err)
	}

	return &JWTVerifier{
		JWKS:     jwks,
		Issuer:   cfg.CognitoOpenIDConfigURL,
		ClientID: cfg.CognitoClientID,
	}, nil
}

func (v *JWTVerifier) VerifyAccessToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(
		tokenString,
		v.JWKS.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.Issuer),
	)
}
