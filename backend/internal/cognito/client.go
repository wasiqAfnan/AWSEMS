package cognito

import (
	"context"
	"log"

	"awsems/internal/config"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Client is a wrapper around the OIDC provider and OAuth2 config which
// can be used to authenticate users.
type Client struct {
	OAuthConfig *oauth2.Config
	Provider    *oidc.Provider
}

func NewClient(cfg *config.Config) *Client {
	ctx := context.Background()

	// Create a new OIDC provider
	provider, err := oidc.NewProvider(
		ctx,
		cfg.CognitoOpenIDConfigURL,
	)
	if err != nil {
		log.Fatalf("failed to create Cognito OIDC provider: %v", err)
	}

	// Create a new OAuth2 config
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.CognitoClientID,
		ClientSecret: cfg.CognitoClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.CognitoRedirectURI,
		Scopes:       []string{"openid", "email", "profile"},
	}

	return &Client{
		OAuthConfig: oauthConfig,
		Provider:    provider,
	}
}
