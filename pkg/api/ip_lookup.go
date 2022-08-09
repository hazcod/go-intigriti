package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	ipLookupURI       = "/company/v2/iplookup"
	ipLookupParamName = "ipAddress"
)

func (e *Endpoint) IsKnownIP(ip net.IP) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL+ipLookupURI, nil)
	if err != nil {
		return false, errors.Wrap(err, "could not create get programs")
	}

	req.URL.Query().Set(ipLookupParamName, ip.String())

	resp, err := e.client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "could not get programs")
	}

	if resp.StatusCode > 399 {
		return false, errors.Errorf("returned status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "could not read response")
	}

	var ipResponse lookupIPResponse
	if err := json.Unmarshal(b, &ipResponse); err != nil {
		return false, errors.Wrap(err, "could not decode response")
	}

	return ipResponse.Exists, nil
}

type lookupIPResponse struct {
	Exists bool `json:"exists"`
}
