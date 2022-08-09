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

func New(clientToken string, clientSecret string, tc *config.TokenCache, logger *logrus.Logger) (Endpoint, error) {
	e := Endpoint{
		clientID:     clientToken,
		clientSecret: clientSecret,
		clientTag:    clientTag,
	}

	if logger == nil {
		e.Logger = logrus.New()
	} else {
		e.Logger = logger
	}

	httpClient, err := e.getClient(tc)
	if err != nil {
		return e, errors.Wrap(err, "could not init client")
	}

	e.client = httpClient

	return e, nil
}
