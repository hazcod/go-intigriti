package ui

import "os/exec"

func (SystemBrowser) OpenURL(url string) error {
	return exec.Command("open", url).Run()
}
