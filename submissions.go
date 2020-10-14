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
	Code              string `json:"code"`
	InternalReference struct {
		Reference string `json:"reference"`
		URL       string `json:"url"`
	} `json:"internalReference"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originalTitle"`
	Program       struct {
		Name                 string `json:"name"`
		Handle               string `json:"handle"`
		LogoURL              string `json:"logoUrl"`
		ConfidentialityLevel struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"confidentialityLevel"`
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
		StatusTrigger struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"statusTrigger"`
	} `json:"program"`
	Type struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Cwe      string `json:"cwe"`
	} `json:"type"`
	Severity struct {
		ID     int    `json:"id"`
		Vector string `json:"vector"`
		Value  string `json:"value"`
	} `json:"severity"`
	Domain struct {
		Value      string `json:"value"`
		Motivation string `json:"motivation"`
		Type       struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"type"`
		BountyTable struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"bountyTable"`
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
	TotalPayout      float64 `json:"totalPayout"`
	CreatedAt        int64   `json:"createdAt"`
	LastUpdatedAt    int64   `json:"lastUpdatedAt"`
	ClosedAt         int64   `json:"closedAt"`
	ValidatedAt      int64   `json:"validatedAt"`
	AcceptedAt       int64   `json:"acceptedAt"`
	ArchivedAt       int64   `json:"archivedAt"`
	AwaitingFeedback bool    `json:"awaitingFeedback"`
	Assignee         struct {
		UserName  string `json:"userName"`
		AvatarURL string `json:"avatarUrl"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	} `json:"assignee"`
	Researcher struct {
		UserName  string `json:"userName"`
		AvatarURL string `json:"avatarUrl"`
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
		UserName  string `json:"userName"`
		AvatarURL string `json:"avatarUrl"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	} `json:"lastUpdater"`
	Links struct {
		Details string `json:"details"`
	} `json:"links"`
	WebLinks struct {
		Details string `json:"details"`
	} `json:"webLinks"`
	SubmissionDetailURL string `json:"submissionDetailUrl"`
}

type Researcher struct {
	Username  string
	AvatarURL string
}

type Program struct {
	Handle string
	Name   string
}

type Submission struct {
	Program    Program
	Researcher Researcher

	DateLastUpdated time.Time
	DateCreated     time.Time
	DateClosed      time.Time

	CWE      string
	Type     string
	Category string

	InternalReference string

	ID       string
	URL      string
	Title    string
	Severity string
	Endpoint string
	State    string
	Payout   float64

	CloseReason string
}

func (s *Submission) IsReady() bool {
	switch strings.ToLower(s.State) {
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

func (s *Submission) IsClosed() bool {
	return s.CloseReason != ""
}

func getBearerTokenHeader(authToken string) string {
	return "Bearer " + authToken
}

func (e *Endpoint) GetSubmissions() ([]Submission, error) {
	var findings []Submission

	if err := authenticate(e); err != nil {
		return findings, errors.Wrap(err, "could not authenticate to intigriti API")
	}

	req, err := http.NewRequest(http.MethodGet, e.apiSubmissions, nil)
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

	defer func() { _ = resp.Body.Close() }()
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

			State:    entry.State.Status.Value,
			Type:     entry.Type.Name,
			CWE:      entry.Type.Cwe,
			Category: entry.Type.Category,
			ID:       entry.Code,

			URL:      entry.SubmissionDetailURL,
			Title:    entry.Title,
			Severity: entry.Severity.Value,

			DateCreated:     time.Unix(entry.CreatedAt, 0),
			DateClosed:      time.Unix(entry.ClosedAt, 0),
			DateLastUpdated: time.Unix(entry.LastUpdatedAt, 0),

			Endpoint: entry.EndpointVulnerableComponent,

			InternalReference: entry.InternalReference.Reference,
			CloseReason: entry.State.CloseReason.Value,
		})
	}

	return findings, nil
}
