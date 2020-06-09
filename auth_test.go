package intigriti

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO implement real auth integration tests

func TestNeedsAuthRefresh(t *testing.T) {
	now := time.Now()

	okTime := now.Add( time.Minute )
	if needsAuthRefresh(okTime, now) {
		t.Error("token needs refresh while it doesnt")
	}

	badTime := now.Add( time.Minute * -1 )
	if ! needsAuthRefresh(badTime, now) {
		t.Error("token never needs refresh")
	}
}

func TestSetNewAuthExpiration(t *testing.T) {
	e := Endpoint{}
	okToken  := time.Hour * 4

	if err := setNewAuthExpiration(&e, int(okToken.Seconds())); err != nil {
		t.Error(err)
	}

	badToken := time.Second * -1
	if err := setNewAuthExpiration(&e, int(badToken.Seconds())); err == nil {
		t.Error("can set invalid auth expiration time")
	}
}

func getTestAuthTokenServer(token string, expiresAtSec int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		authResponse := authResponse{
			AccessToken:	  token,
			ExpiresAtSeconds: expiresAtSec,
			TokenType:        expectedTokenType,
			Scope:            expectedTokenScope,
		}
		bytes, err := json.Marshal(&authResponse)
		if err != nil { w.WriteHeader(http.StatusBadRequest) }
		w.Write(bytes)
	}))
}

func TestGetNewAuthToken(t *testing.T) {
	now := time.Now()
	ts := getTestAuthTokenServer("foo", int(now.Add(time.Hour).Unix()))
	defer ts.Close()

	if _, err := getNewAuthToken(ts.URL, "foo", "bar"); err != nil {
		t.Error(err)
	}
}