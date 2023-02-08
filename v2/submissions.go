package v2

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	submissionUri = "/company/v2/submissions"
)

func (e *Endpoint) GetSubmissions() ([]Submission, error) {
	req, err := http.NewRequest(http.MethodGet, e.URLAPI+submissionUri, nil)
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

	var submissions []Submission

	if err := json.Unmarshal(b, &submissions); err != nil {
		return nil, errors.Wrap(err, "could not decode programs")
	}

	return submissions, nil
}

// json-to-go from https://api.intigriti.com/external/swagger/index.html?urls.primaryName=V1.2#/Submissions/Submissions_Get
type Submission struct {
	Code        string `json:"code"`
	Originators struct {
		ProgramID   string `json:"programId"`
		PentestCode string `json:"pentestCode"`
	} `json:"originators"`
	InternalReference struct {
		Reference string `json:"reference"`
		URL       string `json:"url"`
	} `json:"internalReference"`
	Title  string `json:"title"`
	Report struct {
		OriginalTitle string `json:"originalTitle"`
		Type          struct {
			Name     string `json:"name"`
			Category string `json:"category"`
			Cwe      string `json:"cwe"`
		} `json:"type"`
		Questions []struct {
			Question string `json:"question"`
			Type     struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"type"`
			Answer string `json:"answer"`
		} `json:"questions"`
		Domain struct {
			Name       string `json:"name"`
			Motivation string `json:"motivation"`
			Type       struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"type"`
			Tier struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"tier"`
			Description string `json:"description"`
		} `json:"domain"`
		EndpointVulnerableComponent string `json:"endpointVulnerableComponent"`
		PocDescription              string `json:"pocDescription"`
		Impact                      string `json:"impact"`
		PersonalData                bool   `json:"personalData"`
		RecommendedSolution         string `json:"recommendedSolution"`
		Attachments                 []struct {
			URL  string `json:"url"`
			Code int    `json:"code"`
		} `json:"attachments"`
		IP string `json:"ip"`
	} `json:"report"`
	State struct {
		Status struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"status"`
		CloseReason struct {
			ID    int    `json:"id"`
			Value string `json:"value"`
		} `json:"closeReason"`
		DuplicateInfo struct {
			ParentSubmissionCode string   `json:"parentSubmissionCode"`
			ChildSubmissionCodes []string `json:"childSubmissionCodes"`
		} `json:"duplicateInfo"`
		ValidatedAt int `json:"validatedAt"`
		AcceptedAt  int `json:"acceptedAt"`
		ClosedAt    int `json:"closedAt"`
		ArchivedAt  int `json:"archivedAt"`
	} `json:"state"`
	Severity struct {
		ID     int    `json:"id"`
		Vector string `json:"vector"`
		Value  string `json:"value"`
	} `json:"severity"`
	AwaitingFeedback bool `json:"awaitingFeedback"`
	Reward           struct {
		TotalPayout struct {
			Value    int    `json:"value"`
			Currency string `json:"currency"`
		} `json:"totalPayout"`
		TotalBountyPayout struct {
			Value    int    `json:"value"`
			Currency string `json:"currency"`
		} `json:"totalBountyPayout"`
		TotalBonusPayout struct {
			Value    int    `json:"value"`
			Currency string `json:"currency"`
		} `json:"totalBonusPayout"`
		PossibleBounty struct {
			Value    int    `json:"value"`
			Currency string `json:"currency"`
		} `json:"possibleBounty"`
	} `json:"reward"`
	CreatedAt TimeStamp  `json:"createdAt"`
	Destroyed *Destroyed `json:"-,omitempty"`
	Assignee  struct {
		UserID    string `json:"userId"`
		UserName  string `json:"userName"`
		AvatarURL string `json:"avatarUrl"`
		Role      string `json:"role"`
		Email     string `json:"email"`
	} `json:"assignee"`
	Tags      []string `json:"tags"`
	GroupID   string   `json:"groupId"`
	Submitter struct {
		UserID    string `json:"userId"`
		UserName  string `json:"userName"`
		AvatarURL string `json:"avatarUrl"`
		Role      string `json:"role"`
		Ranking   struct {
			Rank       int `json:"rank"`
			Reputation int `json:"reputation"`
			Streak     struct {
				ID    int    `json:"id"`
				Value string `json:"value"`
			} `json:"streak"`
		} `json:"ranking"`
		IdentityChecked bool `json:"identityChecked"`
	} `json:"submitter"`
	LastUpdated struct {
		LastUpdater struct {
			UserID    string `json:"userId"`
			UserName  string `json:"userName"`
			AvatarURL string `json:"avatarUrl"`
			Role      string `json:"role"`
		} `json:"lastUpdater"`
		LastUpdatedAt int `json:"lastUpdatedAt"`
	} `json:"lastUpdated"`
	AttachmentCount int `json:"attachmentCount"`
	WebLinks        struct {
		Details string `json:"details"`
	} `json:"webLinks"`
	IntegrationCount int `json:"integrationCount"`
}

type Destroyed struct {
	DestroyedBy *DestroyedBy `json:"destroyedBy,omitempty"`
	DestroyedAt TimeStamp    `json:"destroyedAt,omitempty"`
}

type DestroyedBy struct {
	UserID    string `json:"userId"`
	UserName  string `json:"userName"`
	AvatarURL string `json:"avatarUrl"`
	Role      string `json:"role"`
	Email     string `json:"email"`
}

type TimeStamp struct {
	time.Time
}

func (s *TimeStamp) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}
	s.Time = time.Unix(raw, 0)
	return nil
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
