package ui

import "os/exec"

func Open(url string) error {
	return exec.Command("open", url).Run()
}
