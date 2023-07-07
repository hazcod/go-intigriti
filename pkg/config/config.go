package config

import (
	"github.com/sirupsen/logrus"
	"time"
)

type CachedToken struct {
	RefreshToken string
	AccessToken  string
	ExpiryDate   time.Time
	Type         string
}

type InteractiveAuthenticator interface {
	OpenURL(url string) error
}

type Config struct {
	// Required: API authentication credentials
	Credentials struct {
		ClientID     string
		ClientSecret string
	}

	// Optional: open a browser to complete authentication if user interaction is required
	OpenBrowser   bool
	Authenticator InteractiveAuthenticator

	// Optional: token cache if caching previous credentials
	TokenCache *CachedToken

	// Optional: logger instance
	Logger *logrus.Logger

	// Optional: the API scope permissions that the token should have
	// limit this as much as possible to limit token leakage impact
	// https://intigriti.readme.io/reference/api-token-scopes
	APIScopes []string
}
