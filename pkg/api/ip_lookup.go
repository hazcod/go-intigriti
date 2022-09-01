package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
)

const (
	ipLookupURI       = "/company/v2/iplookup"
	ipLookupParamName = "ipAddress"
)

// IsKnownIP verifies whether the IP address is known to the Intigriti platform
// this can be as a researcher or company account
func (e *Endpoint) IsKnownIP(ip net.IP) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL+ipLookupURI, nil)
	if err != nil {
		return false, errors.Wrap(err, "could not create get programs")
	}

	queryValues := req.URL.Query()
	queryValues.Set(ipLookupParamName, ip.String())
	req.URL.RawQuery = queryValues.Encode()

	resp, err := e.client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "could not get programs")
	}

	if resp.StatusCode > 399 {
		return false, errors.Errorf("returned status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
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
