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
	httpTimeoutSec = 15
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

	if tc.AccessToken != "" {
		e.logger.Debug("using cached access token")
		e.oauthToken = &oauth2.Token{
			AccessToken:  tc.AccessToken,
			RefreshToken: tc.RefreshToken,
			Expiry:       tc.ExpiryDate,
			TokenType:    tc.Type,
		}
	}

	if e.oauthToken.Valid() {
		e.logger.Debug("cached access token is valid, skipping authentication")
	} else {
		e.logger.Debug("access token is invalid or expired, authenticating for new token")

		authzCode, err := e.authenticate(ctx, &conf, auth, e.oauthToken.AccessToken)
		if err != nil {
			return nil, errors.Wrap(err, "failed to authenticate")
		}

		if authzCode != "" {
			e.logger.WithField("code", authzCode).Debug("exchanging code")
			e.oauthToken, err = conf.Exchange(ctx, authzCode)
			if err != nil {
				return nil, errors.Wrap(err, "could not exchange code")
			}
		}
	}

	// Ensure our HTTP client uses the OAuth2 credentials
	authHttpClient := conf.Client(ctx, e.oauthToken)

	// Inject a logging middleware into the HTTP client
	authHttpClient.Transport = TaggedRoundTripper{Proxied: authHttpClient.Transport, Logger: e.logger}
	e.logger.Debug("successfully created client")

	return authHttpClient, nil
}

// authenticate authenticates with the Intigriti API using either an access token or interactive OAuth.
func (e *Endpoint) authenticate(ctx context.Context, oauth2Config *oauth2.Config, auth *config.InteractiveAuthenticator, accessToken string) (string, error) {
	// If an access token is provided, validate it
	if accessToken != "" {
		e.logger.Info("validating provided access token")

		// Check token validity by calling an endpoint (e.g., user info or token introspection)
		client := oauth2Config.Client(ctx, &oauth2.Token{AccessToken: accessToken})
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.intigriti.com/v1/userinfo", nil) // Example endpoint
		if err != nil {
			e.logger.WithError(err).Error("failed to create validation request")
			return "", err
		}

		resp, err := client.Do(req)
		if err != nil {
			e.logger.WithError(err).Warn("access token validation failed, proceeding to interactive authentication")
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				e.logger.Debug("access token is valid")
				return accessToken, nil
			} else {
				e.logger.WithField("status", resp.StatusCode).Warn("access token invalid, proceeding to interactive authentication")
			}
		}
	}

	// No valid access token provided, start interactive authentication flow
	e.logger.Info("starting interactive authentication flow")
	state := randomString(stateLengthLetters)
	resultChan := make(chan callbackResult, 1)

	// Set a timeout for the interactive authentication flow
	ctx, cancel := context.WithTimeout(ctx, time.Second*callbackTimeoutSec)
	defer func() { go func() { cancel(); resultChan <- callbackResult{} }() }()

	// Start a listener to handle the callback
	go e.listenForCallback(localCallbackURI, localCallbackHost, localCallbackPort, state, resultChan)

	// Generate the authentication URL
	url := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	// Log the URL only if no valid access token was provided
	e.logger.Warnf("Please authenticate: %s", url)

	// Attempt to open the system browser for authentication
	if auth != nil {
		e.logger.Info("opening system browser to authenticate")
		authenticator := *auth
		if err := authenticator.OpenURL(url); err != nil {
			e.logger.WithField("url", url).WithError(err).Warnf("could not open browser")
		}
	}

	e.logger.Debug("waiting for callback click")

	// Wait for the callback or timeout
	var chanResult callbackResult
	select {
	case <-ctx.Done():
		resultChan <- callbackResult{Error: errors.New("timeout")}
	case chanResult = <-resultChan:
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
