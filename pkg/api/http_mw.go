package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

const (
	clientTag = "hazcod/go-intigriti"
)

type TaggedRoundTripper struct {
	Proxied http.RoundTripper
	Logger  *logrus.Logger
}

// RoundTrip injects a http request header on every request and logs request/response
func (t TaggedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("user-agent", clientTag)

	if t.Logger != nil && t.Logger.IsLevelEnabled(logrus.TraceLevel) && req != nil {
		dumped, err := httputil.DumpRequest(req, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http request")
		} else {
			t.Logger.Trace(string(dumped))
		}
	}

	resp, err := t.Proxied.RoundTrip(req)

	if t.Logger != nil && t.Logger.IsLevelEnabled(logrus.TraceLevel) && resp != nil {
		dumped, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http response")
		} else {
			t.Logger.Trace(string(dumped))
		}
	}

	return resp, err
}
