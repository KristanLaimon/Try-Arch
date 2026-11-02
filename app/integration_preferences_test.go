package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationPreferencesMigrateWithoutChangingRollbackSettings(t *testing.T) {
	dir := t.TempDir()
	if err := saveSettings(settingsPath(dir), settings{MemoryMiB: 4096}); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(settingsPath(dir))
	p, err := loadIntegrationPreferences(dir)
	if err != nil || !p.CaptureFriendly || p.DiscordRPC {
		t.Fatalf("fresh defaults: %+v %v", p, err)
	}
	if err := os.WriteFile(filepath.Join(dir, discordEnabledFilename), nil, 0600); err != nil {
		t.Fatal(err)
	}
	p, err = loadIntegrationPreferences(dir)
	if err != nil || !p.DiscordRPC {
		t.Fatal("legacy Discord opt-in lost", err)
	}
	p.CaptureFriendly = false
	p.DiscordRPC = false
	if err := saveIntegrationPreferences(dir, p); err != nil {
		t.Fatal(err)
	}
	got, err := loadIntegrationPreferences(dir)
	if err != nil || got != p {
		t.Fatalf("roundtrip: %+v %v", got, err)
	}
	after, _ := os.ReadFile(settingsPath(dir))
	if string(before) != string(after) {
		t.Fatal("new preferences changed rollback settings")
	}
}

func TestIntegrationPreferencesRejectDamage(t *testing.T) {
	dir := t.TempDir()
	for _, data := range []string{"null", `{}`, `{"schemaVersion":2}`, `{"schemaVersion":1,"other":true}`, `{"schemaVersion":1} {}`, strings.Repeat(" ", 16385)} {
		if err := os.WriteFile(filepath.Join(dir, integrationPreferencesFilename), []byte(data), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadIntegrationPreferences(dir); err == nil {
			t.Fatalf("accepted damaged preferences: %.80s", data)
		}
	}
}
