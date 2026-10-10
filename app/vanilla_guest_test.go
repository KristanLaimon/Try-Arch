package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnpinnedGuestCannotInstallOrModifyDisk(t *testing.T) {
	dir := t.TempDir()
	disk := filepath.Join(dir, "disk.raw")
	if err := os.WriteFile(disk, []byte("existing VM"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := ensureGuest(&config{dir: dir, disk: disk, fresh: true}, "", defaultSumsSHA256)
	if err == nil || !strings.Contains(err.Error(), "vanilla Arch guest release") {
		t.Fatalf("missing guest pin returned %v", err)
	}
	data, err := os.ReadFile(disk)
	if err != nil || string(data) != "existing VM" {
		t.Fatalf("unpinned guest changed existing disk: %q, %v", data, err)
	}
}
