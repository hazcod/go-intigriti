package api

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

const (
	submissionUri = "/company/v2/programs/%s/submissions"
)

// GetSubmissions returns all submissions for the given program identifier
func (e *Endpoint) GetSubmissions(programId string) ([]Submission, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL+fmt.Sprintf(submissionUri, programId), nil)
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "could not read response")
	}

	var submissions []Submission

	if err := json.Unmarshal(b, &submissions); err != nil {
		return nil, errors.Wrap(err, "could not decode programs")
	}

	return submissions, nil
}

type Submission struct {
	Code              string      `json:"code"`
	InternalReference interface{} `json:"internalReference"`
	Title             string      `json:"title"`
	ProgramID         string      `json:"programId"`
	Severity          struct {
		ID     int         `json:"id"`
		Vector interface{} `json:"vector"`
		Value  string      `json:"value"`
	} `json:"severity"`
	State struct {
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
		CloseReason struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"closeReason"`
	} `json:"state"`
	TotalPayout struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"totalPayout"`
	CreatedAt        int  `json:"createdAt"`
	LastUpdatedAt    int  `json:"lastUpdatedAt"`
	AwaitingFeedback bool `json:"awaitingFeedback"`
	Destroyed        bool `json:"destroyed"`
	Assignee         struct {
		AvatarURL string `json:"avatarUrl"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		UserID    string `json:"userId"`
		Username  string `json:"userName"`
	} `json:"assignee"`
	Tags      []interface{} `json:"tags"`
	GroupID   interface{}   `json:"groupId"`
	Submitter struct {
		Ranking struct {
			Rank       int         `json:"rank"`
			Reputation int         `json:"reputation"`
			Streak     interface{} `json:"streak"`
		} `json:"ranking"`
		IdentityChecked bool   `json:"identityChecked"`
		UserID          string `json:"userId"`
		UserName        string `json:"userName"`
		AvatarURL       string `json:"avatarUrl"`
		Role            string `json:"role"`
	} `json:"submitter"`
	CollaboratorCount int `json:"collaboratorCount"`
	WebLinks          struct {
		Details string `json:"details"`
	} `json:"webLinks"`
}

func (s *Submission) IsClosed() bool {
	return s.State.CloseReason.Value != ""
}

func (s *Submission) IsActive() bool {
	switch strings.ToLower(s.State.Status.Value) {
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
