package ui

import (
	"os/exec"
)

func (SystemBrowser) OpenURL(url string) error {
	return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Run()
}
