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
	token, err := jwt.Parse(
		tokenString,
		v.JWKS.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.Issuer),
	)

	if err != nil {
		return nil, err
	}

	// Verify the token is valid for the given issuer
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token has the correct type
	tokenUse, ok := claims["token_use"].(string)
	if !ok || tokenUse != "access" {
		return nil, fmt.Errorf("invalid token type")
	}

	// Verify the token has the correct client
	clientID, ok := claims["client_id"].(string)
	if !ok || clientID != v.ClientID {
		return nil, fmt.Errorf("invalid client")
	}

	return token, nil
}
