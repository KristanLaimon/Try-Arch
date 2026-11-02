package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const integrationPreferencesFilename = "integration-preferences.json"
const captureFriendlyFilename = "capture-friendly-immersive"

type integrationPreferences struct {
	SchemaVersion   int  `json:"schemaVersion"`
	CaptureFriendly bool `json:"captureFriendly"`
	DiscordRPC      bool `json:"discordRPC"`
}

func loadIntegrationPreferences(dir string) (integrationPreferences, error) {
	p := integrationPreferences{SchemaVersion: 1, CaptureFriendly: true}
	f, err := os.Open(filepath.Join(dir, integrationPreferencesFilename))
	if os.IsNotExist(err) {
		p.DiscordRPC = discordBridgeEnabled(dir)
		return p, nil
	}
	if err != nil {
		return p, err
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil {
		return p, err
	}
	if len(data) > 16384 {
		return p, fmt.Errorf("integration preferences are too large")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	p.SchemaVersion = 0
	if err := decoder.Decode(&p); err != nil {
		return p, err
	}
	if p.SchemaVersion != 1 || decoder.Decode(&struct{}{}) != io.EOF {
		return p, fmt.Errorf("invalid integration preferences")
	}
	return p, nil
}

func saveIntegrationPreferences(dir string, p integrationPreferences) error {
	p.SchemaVersion = 1
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return writeUpdateFile(filepath.Join(dir, integrationPreferencesFilename), append(data, '\n'))
}
