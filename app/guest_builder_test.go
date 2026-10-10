package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestPinnedVanillaReleaseMetadataContract(t *testing.T) {
	sums, err := parseVerifiedSums(defaultSums, defaultSumsSHA256)
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := os.ReadFile("testdata/guest-manifest.arch.json")
	if err != nil {
		t.Fatal(err)
	}
	if testSHA256(metadata) != sums["guest-manifest.json"] {
		t.Fatal("release guest metadata does not match its authenticated manifest")
	}
	if _, err := parseGuestArtifactSizes(metadata, sums); err != nil {
		t.Fatalf("installer rejected actual release metadata: %v", err)
	}
}

// Use the builder's real output rather than an independently handwritten schema.
func TestVanillaBuilderManifestInstallsThroughLauncher(t *testing.T) {
	python, err := exec.LookPath("python")
	if err != nil {
		python, err = exec.LookPath("python3")
	}
	if err != nil {
		t.Fatal("Python is required to validate the guest builder contract")
	}
	directory := t.TempDir()
	files := factoryPayloadFixture(t, []byte("tiny vanilla factory for installer validation"))
	files["packages.txt"] = []byte("base 3-3\nlinux 7.2.9.arch1-1\n")
	for _, name := range []string{"rootfs.ext4", "rootfs.ext4.zst", "vmlinuz-linux", "initramfs-linux.img", "packages.txt"} {
		if err := os.WriteFile(filepath.Join(directory, name), files[name], 0600); err != nil {
			t.Fatal(err)
		}
	}
	if output, err := exec.Command(python, "../guest-build/spec.py", directory).CombinedOutput(); err != nil {
		t.Fatalf("guest builder: %v\n%s", err, output)
	}
	for _, name := range []string{"build-spec.json", "guest-manifest.json"} {
		files[name], err = os.ReadFile(filepath.Join(directory, name))
		if err != nil {
			t.Fatal(err)
		}
	}
	setFixtureSums(files)
	cfg, _, server, _ := onDemandFixtureFiles(t, false, files)
	if err := ensureFactory(cfg, server.URL, testSHA256(files["SHA256SUMS"])); err != nil {
		t.Fatalf("installing builder-produced manifest: %v", err)
	}
	installed, err := os.ReadFile(filepath.Join(cfg.guestDir, "rootfs.ext4"))
	if err != nil || !bytes.Equal(installed, files["rootfs.ext4"]) {
		t.Fatalf("factory installation produced incorrect bytes: %v", err)
	}
}
