package v1

import "testing"

func TestNew(t *testing.T) {
	clientToken := "foo"
	clientSecret := "bar"

	endpoint := New(clientToken, clientSecret)

	if endpoint.clientToken != clientToken {
		t.Fail()
	}

	if endpoint.clientSecret != clientSecret {
		t.Fail()
	}
}
