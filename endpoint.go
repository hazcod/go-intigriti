package intigriti

import (
	"time"
)

const (
	apiSubmissions = "https://api.intigriti.com/external/submission"
	apiAuth = "https://login.intigriti.com/connect/token"

	clientTag = "Go intigriti library"
)

type Endpoint struct {
	clientToken		string
	clientSecret 	string
	clientTag 		string

	authToken		string
	authTokenExp    time.Time
}

func New(clientToken string, clientSecret string) Endpoint {
	return Endpoint{
		clientToken: clientToken,
		clientSecret: clientSecret,
		clientTag: clientTag,
	}
}