package intigriti

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	expectedTokenType  = "bearer"
	expectedTokenScope = "external_api"
	authRefreshNeeded  = time.Second * 10
	httpTimeout 	   = time.Second * 5
	mimeFormUrlEncoded = "application/x-www-form-urlencoded"
)

type authResponse struct {
	AccessToken			string	`json:"access_token"`
	ExpiresAtSeconds	int		`json:"expires_in"`
	TokenType			string 	`json:"token_type"`
	Scope 				string	`json:"scope"`
}

func needsAuthRefresh(token, testTime time.Time) bool {
	return token.Add( authRefreshNeeded ).Before( testTime )
}

func getNewAuthToken(apiUrl, clientId, clientSecret string) (authResponse authResponse, err error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("client_id", clientId)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequest(http.MethodPost, apiUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return authResponse, errors.Wrap(err, "could not create http request")
	}

	req.Header.Set("Content-Type", mimeFormUrlEncoded)
	req.Header.Set("X-Client", clientTag)

	client := http.Client{Timeout: httpTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return authResponse, errors.Wrap(err, "http request failed")
	}

	defer resp.Body.Close()
	if resp.StatusCode > 399 {
		return authResponse, errors.Errorf("received error code: %d", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return authResponse, errors.Wrap(err, "could not read response")
	}

	if err := json.Unmarshal(respBytes, &authResponse); err != nil {
		logrus.Debugf("%+v", string(respBytes))
		return authResponse, errors.Wrap(err, "could not decode auth response")
	}

	return authResponse, nil
}

func setNewAuthExpiration(e *Endpoint, tokenSeconds int) error {
	now := time.Now()

	newExpTime := now.Add(time.Second * time.Duration( tokenSeconds ))
	if newExpTime.Before(now) {
		return errors.Errorf("new expiration time %s is before %s", newExpTime, now)
	}

	e.authTokenExp = newExpTime

	logrus.WithField("token_exp", newExpTime).Debug("new token expiration set")

	return nil
}

func authenticate(e *Endpoint) error {
	now := time.Now()

	if ! needsAuthRefresh(e.authTokenExp, now) {
		logrus.WithField("auth_token_exp", e.authTokenExp).
			Debug("no need to refresh intigriti auth token")
		return nil
	}

	authResponse, err := getNewAuthToken(apiAuth, e.clientToken, e.clientSecret)
	if err != nil {
		return errors.Wrap(err, "could not retrieve new intigriti auth token")
	}

	if err := setNewAuthExpiration(e, authResponse.ExpiresAtSeconds); err != nil {
		return errors.Wrap(err, "invalid token received")
	}

	e.authToken = authResponse.AccessToken

	if strings.ToLower(authResponse.TokenType) != expectedTokenType {
		logrus.WithField("token_type", authResponse.TokenType).Warn("unexpected token type")
	}

	if strings.ToLower(authResponse.Scope) != expectedTokenScope {
		logrus.WithField("token_scope", authResponse.Scope).Warn("unexpected token scope")
	}

	logrus.Debug("authenticated to intigriti")

	return nil
}
