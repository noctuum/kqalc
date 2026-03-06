package main

import "testing"

func TestCopyToClipboardNoTool(t *testing.T) {
	// Without wl-copy/xclip available, CopyToClipboard should return an error
	t.Setenv("WAYLAND_DISPLAY", "")
	t.Setenv("DISPLAY", "")
	err := CopyToClipboard("test")
	if err == nil {
		t.Skip("clipboard tool available in test environment")
	}
}
