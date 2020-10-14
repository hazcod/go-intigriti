package intigriti

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type submissionsResponse []struct {
	SubmissionDetailUrl string `json:"submissionDetailUrl"`
	Code              string `json:"code"`
	InternalReference struct {
		Reference string `json:"reference"`
		URL       string `json:"url"`
	} `json:"internalReference"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Program       struct {
		Name                 string      `json:"name"`
		Handle               string      `json:"handle"`
		LogoURL              interface{} `json:"logoUrl"`
		ConfidentialityLevel struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"confidentialityLevel"`
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
	} `json:"program"`
	Type struct {
		Name     string      `json:"name"`
		Category interface{} `json:"category"`
		Cwe      interface{} `json:"cwe"`
	} `json:"type"`
	Severity struct {
		ID     int         `json:"id"`
		Vector interface{} `json:"vector"`
		Value  string      `json:"value"`
	} `json:"severity"`
	Domain struct {
		Type struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"type"`
		Value          string `json:"value"`
		Motivation     string `json:"motivation"`
		BusinessImpact struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"businessImpact"`
	} `json:"domain"`
	EndpointVulnerableComponent string `json:"endpointVulnerableComponent"`
	State                       struct {
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
		CloseReason struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"closeReason"`
	} `json:"state"`
	TotalPayout      float64     `json:"totalPayout"`
	CreatedAt        uint64         `json:"createdAt"`
	LastUpdatedAt    uint64         `json:"lastUpdatedAt"`
	ValidatedAt      uint64         `json:"validatedAt"`
	AcceptedAt       uint64         `json:"acceptedAt"`
	ClosedAt         uint64         `json:"closedAt"`
	ArchivedAt       uint64         `json:"archivedAt"`
	AwaitingFeedback bool        `json:"awaitingFeedback"`
	Assignee         struct {
		UserName  string      `json:"userName"`
		AvatarURL string      `json:"avatarUrl"`
		Email     string      `json:"email"`
	} `json:"assignee"`
	Researcher struct {
		UserName  string      `json:"userName"`
		AvatarURL string      `json:"avatarUrl"`
		Ranking   struct {
			Rank       int `json:"rank"`
			Reputation int `json:"reputation"`
			Streak     struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"streak"`
		} `json:"ranking"`
		IdentityChecked bool `json:"identityChecked"`
	} `json:"researcher"`
	LastUpdater struct {
		UserName  string      `json:"userName"`
		AvatarURL interface{} `json:"avatarUrl"`
		Email     string      `json:"email"`
	} `json:"lastUpdater"`
}

type Researcher struct {
	Username 	string
	AvatarURL 	string
}

type Program struct {
	Handle string
	Name   string
}

type Submission struct {
	Program     Program
	Researcher  Researcher

	DateCreated uint64
	DateClosed  uint64

	Type 		string
	ID			string
	URL			string
	Title		string
	Severity	string
	Timestamp	time.Time
	Endpoint	string
	State 		string
	Payout      float64
}

func (f *Submission) IsReady() bool {
	switch strings.ToLower(f.State) {
	case "triage":
		return false
	case "closed":
		return false
	case "accepted":
		return false
	case "archived":
		return false
	default:
		return true
	}
}

func getBearerTokenHeader(authToken string) string {
	return "Bearer " + authToken
}

func (e *Endpoint) GetSubmissions() ([]Submission, error) {
	var findings []Submission

	if err := authenticate(e); err != nil {
		return  findings, errors.Wrap(err, "could not authenticate to intigriti API")
	}

	req, err := http.NewRequest(http.MethodGet, apiSubmissions, nil)
	if err != nil {
		return findings, errors.Wrap(err, "could not create http request to intigriti")
	}

	req.Header.Set("Content-Type", mimeFormUrlEncoded)
	req.Header.Set("X-Client", clientTag)
	req.Header.Set("Authorization", getBearerTokenHeader(e.authToken))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return findings, errors.Wrap(err, "fetching to intigriti failed")
	}

	defer func(){ _ = resp.Body.Close() }()
	if resp.StatusCode > 399 {
		return findings, errors.Errorf("fetch from intigriti returned status code: %d", resp.StatusCode)
	}

	var fetchResp submissionsResponse
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return findings, errors.Wrap(err, "could not read response")
	}

	if err := json.Unmarshal(respBytes, &fetchResp); err != nil {
		return findings, errors.Wrap(err, "could not decode intigriti response")
	}

	for _, entry := range fetchResp {
		findings = append(findings, Submission{
			Program: Program{
				Handle: entry.Program.Handle,
				Name:   entry.Program.Name,
			},

			Researcher: Researcher{
				Username:  entry.Researcher.UserName,
				AvatarURL: entry.Researcher.AvatarURL,
			},

			State: 		entry.State.Status.Value,
			Type:		entry.Type.Name,
			ID:			entry.Code,
			// TODO wait for them to implement FE view url
			URL:		"https://intigriti.com/",
			Title:      entry.Title,
			Severity:   entry.Severity.Value,
			DateCreated:entry.CreatedAt,
			DateClosed: entry.ClosedAt,
			Endpoint:   entry.EndpointVulnerableComponent,
		})
	}

	return findings, nil
}