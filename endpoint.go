package intigriti

import (
	"github.com/sirupsen/logrus"
	"time"
)

const (
	apiSubmissions = "https://api.intigriti.com/external/submission"
	apiAuth = "https://login.intigriti.com/connect/token"

	clientTag = "Go intigriti library"
)

type Endpoint struct {
	Logger 			*logrus.Logger

	clientToken		string
	clientSecret 	string
	clientTag 		string

	authToken		string
	authTokenExp    time.Time

	apiAuth 		string
	apiSubmissions 	string
}

func New(clientToken string, clientSecret string) Endpoint {
	e := Endpoint{
		clientToken: clientToken,
		clientSecret: clientSecret,
		clientTag: clientTag,
	}

	e.Logger = logrus.New()
	e.apiAuth = apiAuth
	e.apiSubmissions = apiSubmissions

	return e
}