package main

import (
	"strings"
	"testing"
)

func TestUninstallKeyNameIsStablePerInstall(t *testing.T) {
	def := `C:\Users\me\AppData\Local\TryArch`
	if got := uninstallKeyName(def, def); got != "TryArch" {
		t.Fatalf("default install key: %q", got)
	}
	a := uninstallKeyName(`E:\Arch\TryArch`, def)
	b := uninstallKeyName(`e:\arch\tryarch\`, def)
	if a != b || !strings.HasPrefix(a, "TryArch-") || len(a) != len("TryArch-")+8 {
		t.Fatalf("alternate install keys differ or are malformed: %q %q", a, b)
	}
	if uninstallKeyName(`D:\Other`, def) == a {
		t.Fatal("different folders must not share a key")
	}
	if got := uninstallDisplayName(`E:\Arch\TryArch`, def); got != `Try Arch (E:\Arch\TryArch)` {
		t.Fatalf("display name: %q", got)
	}
}

func TestUninstallCommandQuotesPaths(t *testing.T) {
	got := uninstallCommand(`E:\My Arch\TryArch.exe`, `E:\My Arch`)
	if got != `"E:\My Arch\TryArch.exe" -dir "E:\My Arch" -uninstall` {
		t.Fatalf("command: %s", got)
	}
	if displayVersion("v0.0.12-preview") != "0.0.12-preview" {
		t.Fatal("display version keeps the v prefix")
	}
}
