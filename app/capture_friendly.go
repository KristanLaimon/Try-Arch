package main

import (
	"os"
	"path/filepath"
)

const captureFriendlyFilename = "capture-friendly-immersive"

func captureFriendlyEnabled(dir string) bool {
	info, err := os.Lstat(filepath.Join(dir, captureFriendlyFilename))
	return err == nil && info.Mode().IsRegular()
}
