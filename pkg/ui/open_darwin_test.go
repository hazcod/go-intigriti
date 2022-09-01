package ui

import "testing"

func TestOpen(t *testing.T) {
	if err := Open("/bin/ls"); err != nil {
		t.Error(err)
	}
}
