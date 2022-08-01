package v2

import (
	"github.com/hazcod/go-intigriti/pkg/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

const (
	// TODO: inject on every HTTP request
	clientTag = "hazcod/go-intigriti/1.0"
)

type Endpoint struct {
	Logger *logrus.Logger
	URLAPI string

	clientID     string
	clientSecret string
	clientTag    string

	client     *http.Client
	oauthToken *oauth2.Token
}

func New(clientToken string, clientSecret string, tc *config.TokenCache) (Endpoint, error) {
	e := Endpoint{
		clientID:     clientToken,
		clientSecret: clientSecret,
		clientTag:    clientTag,
	}

	e.Logger = logrus.New()

	httpClient, err := e.getClient(tc)
	if err != nil {
		return e, errors.Wrap(err, "could not init client")
	}

	e.client = httpClient

	return e, nil
}
