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

func (t TaggedRoundTripper) RoundTrip(req *http.Request) (res *http.Response, e error) {
	req.Header.Set("user-agent", clientTag)
	resp, err := t.Proxied.RoundTrip(req)

	if t.Logger != nil && t.Logger.IsLevelEnabled(logrus.TraceLevel) {
		dumped, err := httputil.DumpRequest(resp.Request, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http request")
		} else {
			t.Logger.Trace(string(dumped))
		}

		dumped, err = httputil.DumpResponse(resp, true)
		if err != nil {
			t.Logger.WithError(err).Error("could not dump http response")
		} else {
			t.Logger.Trace(string(dumped))
		}
	}

	return resp, err
}
