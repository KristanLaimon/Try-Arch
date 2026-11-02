package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCaptureFriendlyMarkerIsOptIn(t *testing.T) {
	dir := t.TempDir()
	if captureFriendlyEnabled(dir) {
		t.Fatal("capture mode enabled without opt-in")
	}
	if err := os.WriteFile(filepath.Join(dir, captureFriendlyFilename), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if !captureFriendlyEnabled(dir) {
		t.Fatal("capture mode ignored opt-in")
	}
	if !backupNameAllowed(captureFriendlyFilename) {
		t.Fatal("backup omits capture preference")
	}
}

func TestCaptureFriendlyImmersiveKeepsGPUWithoutSDLFullscreen(t *testing.T) {
	cfg := &config{fullscreen: true, captureFriendly: true, useGpu: true, displays: 1, memMiB: 4096, cpus: 2}
	args := buildQemuArgs(cfg, "root=/dev/vda")
	if slices.Contains(args, "-full-screen") {
		t.Fatal("capture mode requested SDL fullscreen")
	}
	if !slices.Contains(args, sdlDisplay(true, false)) {
		t.Fatal("capture mode disabled GPU display")
	}
	cfg.captureFriendly = false
	if !slices.Contains(buildQemuArgs(cfg, "root=/dev/vda"), "-full-screen") {
		t.Fatal("ordinary Immersive lost fullscreen")
	}
	monitor := screenRect{-1920, -100, 0, 980}
	bounds := captureFriendlyBounds(monitor)
	if bounds.Left != monitor.Left || bounds.Top != monitor.Top || bounds.Right != monitor.Right || bounds.height() != 1079 {
		t.Fatalf("unexpected capture bounds: %+v", bounds)
	}
}
