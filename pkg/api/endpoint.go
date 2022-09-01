package api

import (
	"net/http"

	"github.com/intigriti/sdk-go/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Endpoint struct {
	Logger *logrus.Logger

	clientID     string
	clientSecret string
	clientTag    string

	client     *http.Client
	oauthToken *oauth2.Token
}

// New creates an Intigriti endpoint object to use
// this is the main object to interact with the SDK
func New(clientToken string, clientSecret string, tc *config.TokenCache, logger *logrus.Logger) (Endpoint, error) {
	e := Endpoint{
		clientID:     clientToken,
		clientSecret: clientSecret,
		clientTag:    clientTag,
	}

	// initialize the logger to use
	if logger == nil {
		e.Logger = logrus.New()
	} else {
		e.Logger = logger
	}

	// prepare our oauth2-ed http client
	httpClient, err := e.getClient(tc)
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
