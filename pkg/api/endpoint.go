package api

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Endpoint struct {
	logger *logrus.Logger

	clientID     string
	clientSecret string
	clientTag    string

	client     *http.Client
	oauthToken *oauth2.Token
}

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
	OpenBrowser bool

	// optional token cache if caching previous credentials
	TokenCache *CachedToken

	// optional logger instance
	Logger *logrus.Logger
}

// New creates an Intigriti endpoint object to use
// this is the main object to interact with the SDK
func New(cfg Config) (Endpoint, error) {
	e := Endpoint{
		clientID:     cfg.Credentials.ClientID,
		clientSecret: cfg.Credentials.ClientSecret,
		clientTag:    clientTag,
	}

	// initialize the logger to use
	if cfg.Logger == nil {
		e.logger = logrus.New()
	} else {
		e.logger = cfg.Logger
	}

	// prepare our oauth2-ed http client
	httpClient, err := e.getClient(cfg.TokenCache, cfg.OpenBrowser)
	if err != nil {
		return e, errors.Wrap(err, "could not init client")
	}

	e.client = httpClient

	// ensure our current token is fetched or renewed if expired
	if _, err = e.getToken(); err != nil {
		return e, errors.Wrap(err, "could not prepare token")
	}

	return e, nil
}

// IsAuthenticated returns whether the current SDK instance has successfully authenticated
func (e *Endpoint) IsAuthenticated() bool {
	if e.oauthToken == nil {
		return false
	}

	return e.oauthToken.Valid()
}
