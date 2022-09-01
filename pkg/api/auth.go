package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/intigriti/sdk-go/pkg/config"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	httpTimeoutSec     = 5
	stateLengthLetters = 10

	localCallbackURL = "http://localhost:8999/"

	defaultApiTokenURL = "https://login.intigriti.com/connect/token"
	defaultApiAuthzURL = "https://login.intigriti.com/connect/authorize"
	defaultApiEndpoint = "https://api.intigriti.com/external"
)

var (
	tokenURL = os.Getenv("INTI_TOKEN_URL")
	authzURL = os.Getenv("INTI_AUTH_URL")
	apiURL   = os.Getenv("INTI_API_URL")
)

func init() {
	if tokenURL == "" {
		tokenURL = defaultApiTokenURL
	}

	if authzURL == "" {
		authzURL = defaultApiAuthzURL
	}

	if apiURL == "" {
		apiURL = defaultApiEndpoint
	}
}

func (e *Endpoint) getOauth2Config() oauth2.Config {
	e.Logger.WithField("api_url", apiURL).Debug("set api url")

	return oauth2.Config{
		ClientID:     e.clientID,
		ClientSecret: e.clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
			AuthURL:  authzURL,
		},
		RedirectURL: localCallbackURL,
		Scopes:      []string{"external_api", "offline_access"},
	}
}

func (e *Endpoint) GetToken() (*oauth2.Token, error) {
	if e.oauthToken != nil && e.oauthToken.Valid() {
		return e.oauthToken, nil
	}

	conf := e.getOauth2Config()

	tokenSrc := conf.TokenSource(context.Background(), e.oauthToken)
	token, err := tokenSrc.Token()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve refresh token")
	}

	return token, nil
}

func (e *Endpoint) getClient(tc *config.TokenCache) (*http.Client, error) {
	ctx := context.Background()

	conf := e.getOauth2Config()

	httpClient := &http.Client{Timeout: httpTimeoutSec * time.Second}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	e.oauthToken = &oauth2.Token{}
	if tc.RefreshToken != "" {
		e.Logger.Debug("trying to use cached token")
		e.oauthToken.AccessToken = tc.AccessToken
		e.oauthToken.RefreshToken = tc.RefreshToken
		e.oauthToken.Expiry = tc.ExpiryDate
		e.oauthToken.TokenType = tc.Type
	}

	if !e.oauthToken.Valid() {
		e.Logger.Debug("authenticating for new token")

		authzCode, err := e.authenticate(ctx, &conf)
		if err != nil {
			return nil, errors.Wrap(err, "failed to authenticate")
		}

		e.Logger.Debug("exchanging code")

		e.oauthToken, err = conf.Exchange(ctx, authzCode)
		if err != nil {
			return nil, errors.Wrap(err, "could not exchange code")
		}
	}

	authHttpClient := conf.Client(ctx, e.oauthToken)
	authHttpClient.Transport = TaggedRoundTripper{Proxied: authHttpClient.Transport, Logger: e.Logger}
	e.Logger.Debug("successfully created client")

	return authHttpClient, nil
}

func (e *Endpoint) authenticate(ctx context.Context, oauth2Config *oauth2.Config) (string, error) {
	state := randomString(stateLengthLetters)

	resultChan := make(chan callbackResult, 1)

	go e.listenForCallback(8999, state, resultChan)
	defer func() { go func() { resultChan <- callbackResult{} }() }()

	// Redirect user to consent page to ask for permission 	for the scopes specified above.
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	e.Logger.Warnf("Please authenticate: %s", url)

	e.Logger.Debug("waiting for callback click")

	var chanResult callbackResult
	select {
	case <-ctx.Done():
		resultChan <- callbackResult{Error: errors.New("timeout")}
		break
	case chanResult = <-resultChan:
		break
	}

	e.Logger.WithField("result", chanResult).Debug("received callback result")

	if chanResult.Error != nil {
		return "", chanResult.Error
	}

	if chanResult.Code == "" {
		return "", errors.New("got empty code")
	}

	e.Logger.WithField("code", chanResult.Code).Debug("successfully retrieved new code")
	return chanResult.Code, nil
}
