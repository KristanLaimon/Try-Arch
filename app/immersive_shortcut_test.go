package main

import "testing"

func TestCaptureShortcutDebouncesAndOwnsReleaseAfterFocusChange(t *testing.T) {
	var key immersiveShortcut
	if toggle, swallow := key.route(true, true, true, true, false, true); !toggle || !swallow {
		t.Fatal("chord was not captured")
	}
	if toggle, swallow := key.route(true, true, true, true, false, true); toggle || !swallow {
		t.Fatal("repeat toggled fullscreen")
	}
	if toggle, swallow := key.route(true, false, false, false, false, false); toggle || !swallow {
		t.Fatal("owned key release leaked after focus changed")
	}
	if _, swallow := key.route(true, false, true, true, false, true); swallow {
		t.Fatal("host key was captured")
	}
	if _, swallow := key.route(true, true, true, true, true, true); swallow {
		t.Fatal("AltGr key was captured")
	}
	if _, swallow := key.route(false, true, true, true, false, true); swallow {
		t.Fatal("legacy fullscreen was intercepted")
	}
}
