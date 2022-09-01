package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

const (
	clientTag = "intigriti/sdk-go"
)

type TaggedRoundTripper struct {
	Proxied http.RoundTripper
	Logger  *logrus.Logger
}

// RoundTrip injects a http request header on every request and logs request/response
func (t TaggedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("user-agent", clientTag)
	resp, err := t.Proxied.RoundTrip(req)

	if t.Logger != nil && t.Logger.IsLevelEnabled(logrus.TraceLevel) {
		dumped, err := httputil.DumpRequest(req, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http request")
		} else {
			t.Logger.Trace(string(dumped))
		}

		dumped, err = httputil.DumpResponse(req.Response, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http response")
		} else {
			t.Logger.Trace(string(dumped))
		}
	}

	return resp, err
}
