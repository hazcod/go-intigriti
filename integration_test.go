package intigriti

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {
	token := os.Getenv("TOKEN")
	secret := os.Getenv("SECRET")
	apiAuth := os.Getenv("AUTH_API")
	apiSub := os.Getenv("SUB_API")

	if token == "" || secret == "" {
		t.Fatal("no token or secret supplied to test via env")
	}

	inti := New(token, secret)
	inti.Logger = logrus.New()
	inti.Logger.SetLevel(logrus.DebugLevel)
	if apiAuth != "" { inti.apiAuth= apiAuth }
	if apiSub != "" { inti.apiSubmissions = apiSub }

	subs, err := inti.GetSubmissions()
	if err != nil {
		t.Errorf("could not fetch submissions: %v", err)
	}

	if len(subs) == 0 {
		t.Error("no submissions returned")
	}

	for _, sub := range subs {
		if sub.ID == "" || sub.Endpoint == "" || sub.URL == "" {
			t.Error("empty submission details")
		}
	}
}
