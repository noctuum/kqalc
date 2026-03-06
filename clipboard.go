package main

import (
	"os"
	"os/exec"
	"strings"
)

// CopyToClipboard copies text to the system clipboard.
// Uses wl-copy on Wayland, xclip on X11.
func CopyToClipboard(text string) error {
	var cmd *exec.Cmd
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		cmd = exec.Command("wl-copy", text)
	} else {
		cmd = exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
	}
	return cmd.Run()
}
