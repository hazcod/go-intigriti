package v2

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const (
	programURI = "/company/v2/programs"
)

func (e *Endpoint) GetPrograms() ([]Program, error) {
	req, err := http.NewRequest(http.MethodGet, e.URLAPI+programURI, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not create get programs")
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get programs")
	}

	if resp.StatusCode > 399 {
		return nil, errors.Errorf("returned status %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response")
	}

	var programs []Program

	if err := json.Unmarshal(b, &programs); err != nil {
		return nil, errors.Wrap(err, "could not decode programs")
	}

	return programs, nil
}

type Program struct {
	ID            string `json:"id"`
	Handle        string `json:"handle"`
	CompanyID     string `json:"companyId"`
	CompanyHandle string `json:"companyHandle"`
	LogoURL       string `json:"logoUrl"`
	Name          string `json:"name"`
	Status        struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"status"`
	ConfidentialityLevel struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"confidentialityLevel"`
	WebLinks struct {
		Details string `json:"details"`
	} `json:"webLinks"`
	Type struct {
		ID    int    `json:"id"`
		Value string `json:"value"`
	} `json:"type"`
}
