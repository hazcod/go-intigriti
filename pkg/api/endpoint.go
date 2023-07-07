package api

import (
	"github.com/hazcod/go-intigriti/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

const (
	apiAllScopes = "offline_access company_external_api core_platform:read core_platform:write"
)

type Endpoint struct {
	logger *logrus.Logger

	clientID     string
	clientSecret string
	clientTag    string

	client     *http.Client
	oauthToken *oauth2.Token

	apiScopes []string
}

// New creates an Intigriti endpoint object to use
// this is the main object to interact with the SDK
func New(cfg config.Config) (Endpoint, error) {
	e := Endpoint{
		clientID:     cfg.Credentials.ClientID,
		clientSecret: cfg.Credentials.ClientSecret,
		clientTag:    clientTag,
		apiScopes:    cfg.APIScopes,
	}

	if len(e.apiScopes) == 0 {
		e.apiScopes = strings.Split(apiAllScopes, " ")
	}

	// initialize the logger to use
	if cfg.Logger == nil {
		e.logger = logrus.New()
	} else {
		e.logger = cfg.Logger
	}

	// prepare our oauth2-ed http client
	authenticator := &cfg.Authenticator
	if !cfg.OpenBrowser {
		authenticator = nil
	}

	httpClient, err := e.getClient(cfg.TokenCache, authenticator)
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
