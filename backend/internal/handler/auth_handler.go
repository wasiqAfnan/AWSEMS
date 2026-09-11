package handler

import (
	"log"
	"net/http"

	"awsems/internal/cognito"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type AuthHandler struct {
	Cognito *cognito.Client
}

func NewAuthHandler(cognitoClient *cognito.Client) *AuthHandler {
	return &AuthHandler{
		Cognito: cognitoClient,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	sessionID, err := cognito.CreateOAuthSession()
	if err != nil {
		http.Error(w, "Failed to create OAuth session", http.StatusInternalServerError)
		return
	}

	session, ok := cognito.GetOAuthSession(sessionID)
	if !ok {
		http.Error(w, "Failed to create OAuth session", http.StatusInternalServerError)
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
		http.Error(w, "Missing code or state", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("oauth_session")
	if err != nil {
		http.Error(w, "OAuth session not found", http.StatusBadRequest)
		return
	}

	session, ok := cognito.GetOAuthSession(cookie.Value)
	if !ok {
		http.Error(w, "OAuth session expired or invalid", http.StatusBadRequest)
		return
	}

	if state != session.State {
		http.Error(w, "Invalid state", http.StatusBadRequest)
		return
	}

	token, err := h.Cognito.OAuthConfig.Exchange(
		r.Context(),
		code,
		oauth2.VerifierOption(session.CodeVerifier),
	)
	if err != nil {
		http.Error(w, "Failed to exchange authorization code", http.StatusBadRequest)
		return
	}

	log.Printf("OAuth session cookie: %s", cookie.Value)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "ID token not found", http.StatusBadRequest)
		return
	}

	idToken, err := h.Cognito.Provider.Verifier(&oidc.Config{
		ClientID: h.Cognito.OAuthConfig.ClientID,
	}).Verify(r.Context(), rawIDToken)

	if err != nil {
		http.Error(w, "Invalid ID token", http.StatusUnauthorized)
		return
	}

	var claims struct {
		Subject string `json:"sub"`
		Email   string `json:"email"`
	}

	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to read user information", http.StatusInternalServerError)
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

	http.Redirect(w, r, "/", http.StatusFound)
}
