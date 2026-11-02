package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOrdinaryLauncherUpgradePreservesInstalledBootAssets(t *testing.T) {
	dir := t.TempDir()
	guest := filepath.Join(dir, "guest")
	if err := os.MkdirAll(guest, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(guest, "vmlinuz-linux")
	if err := os.WriteFile(path, []byte("installed kernel matching user modules"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg := &config{dir: dir, guestDir: guest}
	pins := pinnedPayloadUpdate{Release: "https://example.com/new", Digest: testSHA256([]byte("new guest"))}
	if ownPayloadUpdate(cfg, true, false, nil, pins) != nil {
		t.Fatal("ordinary upgrade staged replacement guest assets")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "installed kernel matching user modules" {
		t.Fatal("installed kernel changed", err)
	}
}
