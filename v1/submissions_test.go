package intigriti

import "testing"

// TODO implement real submission integration tests

func TestGetBearerToken(t *testing.T) {
	has := "myCoolToken"
	wants := "Bearer myCoolToken"

	if wants != getBearerTokenHeader(has) {
		t.Error("invalid bearer token")
	}
}

// TODO test get submissions