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

type Config struct {
	// required authentication credentials
	Credentials struct {
		ClientID     string
		ClientSecret string
	}

	// optional open a browser to complete authentication if user interaction is required
	// TODO: move ui/ to the commandline client subpackage
	OpenBrowser bool

	// optional token cache if caching previous credentials
	TokenCache *CachedToken

	// optional logger instance
	Logger *logrus.Logger
}
