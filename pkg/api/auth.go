package api

import (
	"context"
	"fmt"
	"github.com/hazcod/go-intigriti/pkg/config"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	// timeout of every http request
	httpTimeoutSec = 5
	// the length of our Oauth2 state parameter
	stateLengthLetters = 10
	// timeout of the local callback listener
	callbackTimeoutSec = 120

	// local callback url listener
	localCallbackPort = 1337
	localCallbackHost = "localhost"
	localCallbackURI  = "/"

	// default production API endpoints
	defaultApiTokenURL = "https://login.intigriti.com/connect/token"
	defaultApiAuthzURL = "https://login.intigriti.com/connect/authorize"
	defaultApiEndpoint = "https://api.intigriti.com/external"
)

var (
	// used to override the API endpoints at runtime for testing
	tokenURL = os.Getenv("INTI_TOKEN_URL")
	authzURL = os.Getenv("INTI_AUTH_URL")
	apiURL   = os.Getenv("INTI_API_URL")
)

// used if we do local testing to non-production endpoints
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

// retrieve the oauth2 configuration to use
func (e *Endpoint) getOauth2Config(apiScopes []string) oauth2.Config {
	e.logger.WithField("api_url", apiURL).Debug("set api url")

	oauthConfig := oauth2.Config{
		ClientID:     e.clientID,
		ClientSecret: e.clientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: tokenURL,
			AuthURL:  authzURL,
		},
		RedirectURL: fmt.Sprintf("http://%s:%d%s", localCallbackHost, localCallbackPort, localCallbackURI),
		Scopes:      apiScopes,
	}

	e.logger.Tracef("%+v", oauthConfig)

	return oauthConfig
}

// fetch the latest (valid) oauth2 access and refresh token
func (e *Endpoint) getToken() (*oauth2.Token, error) {
	// don't do anything when the token is ok
	if e.oauthToken != nil && e.oauthToken.Valid() {
		return e.oauthToken, nil
	}

	// get out oauth2 config to use
	conf := e.getOauth2Config(e.apiScopes)

	// get valid refresh and access tokens
	tokenSrc := conf.TokenSource(context.Background(), e.oauthToken)
	token, err := tokenSrc.Token()
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve refresh token")
	}

	return token, nil
}

// return the http client which automatically injects the right authentication credentials
func (e *Endpoint) getClient(tc *config.CachedToken, auth *config.InteractiveAuthenticator) (*http.Client, error) {
	ctx := context.Background()

	conf := e.getOauth2Config(e.apiScopes)

	httpClient := &http.Client{Timeout: httpTimeoutSec * time.Second}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	e.oauthToken = &oauth2.Token{}

	if tc == nil {
		tc = &config.CachedToken{}
	}

	// if our configuration contains a cached token, re-use it
	if tc.RefreshToken != "" {
		e.logger.Debug("trying to use cached token")
		e.oauthToken.AccessToken = tc.AccessToken
		e.oauthToken.RefreshToken = tc.RefreshToken
		e.oauthToken.Expiry = tc.ExpiryDate
		e.oauthToken.TokenType = tc.Type
	}

	// if the current token is invalid, fetch a new one
	if !e.oauthToken.Valid() {
		e.logger.Debug("authenticating for new token")

		authzCode, err := e.authenticate(ctx, &conf, auth)
		if err != nil {
			return nil, errors.Wrap(err, "failed to authenticate")
		}

		e.logger.WithField("code", authzCode).Debug("exchanging code")

		e.oauthToken, err = conf.Exchange(ctx, authzCode)
		if err != nil {
			return nil, errors.Wrap(err, "could not exchange code")
		}
	}

	// ensure our http client uses our oauth2 credentials
	authHttpClient := conf.Client(ctx, e.oauthToken)

	// inject a logging middleware into our http client
	authHttpClient.Transport = TaggedRoundTripper{Proxied: authHttpClient.Transport, Logger: e.logger}
	e.logger.Debug("successfully created client")

	return authHttpClient, nil
}

// authenticate versus the Intigriti API, this requires user interaction
func (e *Endpoint) authenticate(ctx context.Context, oauth2Config *oauth2.Config, auth *config.InteractiveAuthenticator) (string, error) {
	state := randomString(stateLengthLetters)

	resultChan := make(chan callbackResult, 1)

	ctx, cancel := context.WithTimeout(ctx, time.Second*callbackTimeoutSec)

	go e.listenForCallback(localCallbackURI, localCallbackHost, localCallbackPort, state, resultChan)
	defer func() { go func() { cancel(); resultChan <- callbackResult{} }() }()

	// Redirect user to consent page to ask for permission 	for the scopes specified above.
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	e.logger.Warnf("Please authenticate: %s", url)

	if auth != nil {
		e.logger.Info("opening system browser to authenticate")

		authenticator := *auth
		if err := authenticator.OpenURL(url); err != nil {
			e.logger.WithField("url", url).WithError(err).Warnf("could not open browser")
		}
	}

	e.logger.Debug("waiting for callback click")

	var chanResult callbackResult
	select {
	case <-ctx.Done():
		resultChan <- callbackResult{Error: errors.New("timeout")}
		break
	case chanResult = <-resultChan:
		break
	}

	e.logger.WithField("result", chanResult).Debug("received callback result")

	if chanResult.Error != nil {
		return "", chanResult.Error
	}

	if chanResult.Code == "" {
		return "", errors.New("got empty code")
	}

	e.logger.WithField("code", chanResult.Code).Debug("successfully retrieved new code")
	return chanResult.Code, nil
}
