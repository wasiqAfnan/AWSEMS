package cognito

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
	"time"

	"golang.org/x/oauth2"
)

type OAuthSession struct {
	State        string
	CodeVerifier string
	ExpiresAt    time.Time
}

var (
	oauthSessions = make(map[string]OAuthSession)
	oauthMutex    sync.Mutex // protects concurrent access to oauthSessions
)

func GenerateRandomValue() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func CreateOAuthSession() (string, error) {
	sessionID, err := GenerateRandomValue()
	if err != nil {
		return "", err
	}

	state, err := GenerateRandomValue()
	if err != nil {
		return "", err
	}

	verifier := oauth2.GenerateVerifier()

	StoreOAuthSession(sessionID, OAuthSession{
		State:        state,
		CodeVerifier: verifier,
		ExpiresAt:    time.Now().Add(2 * time.Minute),
	})

	return sessionID, nil
}

func StoreOAuthSession(sessionID string, session OAuthSession) {
	oauthMutex.Lock()
	defer oauthMutex.Unlock()

	oauthSessions[sessionID] = session
}

func GetOAuthSession(sessionID string) (OAuthSession, bool) {
	oauthMutex.Lock()
	defer oauthMutex.Unlock()

	session, exists := oauthSessions[sessionID]

	if !exists || time.Now().After(session.ExpiresAt) {
		delete(oauthSessions, sessionID)
		return OAuthSession{}, false
	}

	return session, true
}

func DeleteOAuthSession(sessionID string) {
	oauthMutex.Lock()
	defer oauthMutex.Unlock()

	delete(oauthSessions, sessionID)
}
