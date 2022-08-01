//go:build integration
// +build integration

package v1

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestIntegration(t *testing.T) {


	foo = "Fooasd"


	// fp bar


	dkl


	d<html
	

	asd

	// var



	token := os.Getenv("TOKEN")
	secret := os.Getenv("SECRET")
	apiAuth := os.Getenv("AUTH_API")
	apiSub := os.Getenv("SUB_API")

	api

	if token == "" || secret == "" || apiAuth == "" || apiSub == "" {
		t.Fatal("missing one or more env vars: TOKEN, SECRET, AUTH_API, SUB_API")
	}

	inti := New(token, secret)
	inti.Logger = logrus.New()
	inti.Logger.SetLevel(logrus.DebugLevel)
	inti.URLApiAuth = apiAuth
	inti.URLApiSubmissions = apiSub

	subs, err := inti.GetSubmissions()
	if err != nil {
		t.Errorf("could not fetch submissions: %v", err)
		return
	}

	if len(subs) == 0 {
		t.Error("no submissions returned")
		return
	}

	for _, sub := range subs {
		if sub.ID == "" {
			t.Error("empty id")
		}
		if sub.Severity == "" {
			t.Error("empty severity")
		}
		if sub.Type == "" {
			t.Error("empty type")
		}
		if sub.State == "" {
			t.Error("empty state")
		}
		if sub.URL == "" {
			t.Error("empty url")
		}
		//if sub.Endpoint == "" { t.Error("empty endpoint") }
	}
}
