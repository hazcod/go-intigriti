package ui

import "testing"

func TestOpen(t *testing.T) {
	browser := SystemBrowser{}
	if err := browser.OpenURL("/bin/ls"); err != nil {
		t.Error(err)
	}
}
