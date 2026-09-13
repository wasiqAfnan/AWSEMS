package handler

import (
	"net/http"

	"awsems/internal/cognito"
	"awsems/internal/config"
	"awsems/internal/utils"

	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	Cognito *cognito.Client
	Config  *config.Config
}

func NewAuthHandler(cognitoClient *cognito.Client, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		Cognito: cognitoClient,
		Config:  cfg,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	sessionID, err := cognito.CreateOAuthSession()
	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"Failed to create OAuth session",
		)
		return
	}

	session, ok := cognito.GetOAuthSession(sessionID)
	if !ok {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"Failed to create OAuth session",
		)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   120,
	})

	authURL := h.Cognito.OAuthConfig.AuthCodeURL(
		session.State,
		oauth2.S256ChallengeOption(session.CodeVerifier),
	)

	http.Redirect(w, r, authURL, http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" || state == "" {
		utils.Error(w, http.StatusBadRequest, "Missing code or state")
		return
	}

	cookie, err := r.Cookie("oauth_session")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "OAuth session not found")
		return
	}

	session, ok := cognito.GetOAuthSession(cookie.Value)
	if !ok {
		utils.Error(w, http.StatusBadRequest, "OAuth session expired or invalid")
		return
	}

	if state != session.State {
		utils.Error(w, http.StatusBadRequest, "Invalid state")
		return
	}

	token, err := h.Cognito.OAuthConfig.Exchange(
		r.Context(),
		code,
		oauth2.VerifierOption(session.CodeVerifier),
	)
	if err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"Failed to exchange authorization code",
		)
		return
	}

	// log.Printf("OAuth session cookie: %s", cookie.Value)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		utils.Error(w, http.StatusBadRequest, "ID token not found")
		return
	}

	idToken, err := h.Cognito.Provider.Verifier(&oidc.Config{
		ClientID: h.Cognito.OAuthConfig.ClientID,
	}).Verify(r.Context(), rawIDToken)

	if err != nil {
		utils.Error(w, http.StatusUnauthorized, "Invalid ID token")
		return
	}

	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
	}

	if err := idToken.Claims(&claims); err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"Failed to read user information",
		)
		return
	}

	cognito.DeleteOAuthSession(cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  token.Expiry,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    token.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, h.Config.FrontendURL, http.StatusFound)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	// Delete access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Delete refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Delete temporary OAuth session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	// Build Cognito logout URL
	logoutURL := h.Cognito.OAuthConfig.Endpoint.AuthURL

	u, err := url.Parse(logoutURL)
	if err != nil {
		utils.Error(
			w,
			http.StatusInternalServerError,
			"Failed to create logout URL",
		)
		return
	}

	u.Path = "/logout"

	query := u.Query()
	query.Set("client_id", h.Cognito.OAuthConfig.ClientID)
	query.Set("logout_uri", h.Config.CognitoLogoutURI)
	u.RawQuery = query.Encode()

	// Redirect to Cognito
	http.Redirect(w, r, u.String(), http.StatusFound)
}
