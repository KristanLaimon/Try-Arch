package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDiscordOptInKeepsRollbackSettingsReadable(t *testing.T) {
	dir := t.TempDir()
	if discordBridgeEnabled(dir) {
		t.Fatal("bridge enabled without opt-in")
	}
	if err := saveSettings(settingsPath(dir), settings{MemoryMiB: 4096}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, discordEnabledFilename), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if !discordBridgeEnabled(dir) {
		t.Fatal("opt-in marker ignored")
	}
	if _, err := loadSettings(settingsPath(dir)); err != nil {
		t.Fatal("bridge changed rollback settings", err)
	}
	if !backupNameAllowed(discordEnabledFilename) {
		t.Fatal("backup excludes opt-in")
	}
}

func TestDiscordPresenceRelay(t *testing.T) {
	guest, bridge := net.Pipe()
	pipe, discord := net.Pipe()
	defer guest.Close()
	defer discord.Close()
	guest.SetDeadline(time.Now().Add(5 * time.Second))
	discord.SetDeadline(time.Now().Add(5 * time.Second))
	done := make(chan error, 1)
	go func() { done <- relayDiscord(bridge, func() (io.ReadWriteCloser, error) { return pipe, nil }) }()
	server := make(chan error, 1)
	go func() {
		opcode, data, err := readDiscordFrame(discord)
		if err == nil && (opcode != 0 || !validDiscordRequest(opcode, data, true)) {
			err = io.ErrUnexpectedEOF
		}
		if err == nil {
			err = writeDiscordFrame(discord, 1, []byte(`{"evt":"READY"}`))
		}
		if err == nil {
			opcode, data, err = readDiscordFrame(discord)
			if err == nil && (opcode != 1 || !bytes.Contains(data, []byte(`"activity":null`))) {
				err = io.ErrUnexpectedEOF
			}
			var request struct {
				Args struct {
					PID int `json:"pid"`
				} `json:"args"`
			}
			if err == nil {
				err = json.Unmarshal(data, &request)
				if request.Args.PID != os.Getpid() {
					err = io.ErrUnexpectedEOF
				}
			}
		}
		if err == nil {
			err = writeDiscordFrame(discord, 1, []byte(`{"cmd":"SET_ACTIVITY","nonce":"1"}`))
		}
		server <- err
	}()
	if err := writeDiscordFrame(guest, 0, []byte(`{"v":1,"client_id":"123456789012345678"}`)); err != nil {
		t.Fatal(err)
	}
	if _, data, err := readDiscordFrame(guest); err != nil || !bytes.Contains(data, []byte("READY")) {
		t.Fatalf("ready: %s %v", data, err)
	}
	if err := writeDiscordFrame(guest, 1, []byte(`{"cmd":"SET_ACTIVITY","args":{"pid":123,"activity":null},"nonce":"1"}`)); err != nil {
		t.Fatal(err)
	}
	if _, data, err := readDiscordFrame(guest); err != nil || !bytes.Contains(data, []byte("SET_ACTIVITY")) {
		t.Fatalf("reply: %s %v", data, err)
	}
	if err := <-server; err != nil {
		t.Fatal(err)
	}
	guest.Close()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("relay did not stop after app exit")
	}
}

func TestDiscordRejectsAccountCommandsAndBadFrames(t *testing.T) {
	for _, command := range []string{"AUTHORIZE", "AUTHENTICATE", "GET_GUILDS", "SUBSCRIBE", "SELECT_VOICE_CHANNEL"} {
		if validDiscordRequest(1, []byte(`{"cmd":"`+command+`"}`), false) {
			t.Fatalf("allowed %s", command)
		}
	}
	if validDiscordRequest(0, []byte(`{"v":1,"client_id":"../../file"}`), true) {
		t.Fatal("accepted invalid app ID")
	}
	if validDiscordRequest(0, []byte(`{"v":2,"client_id":"123"}`), true) {
		t.Fatal("accepted wrong RPC version")
	}
	var header [8]byte
	binary.LittleEndian.PutUint32(header[4:], discordFrameLimit+1)
	if _, _, err := readDiscordFrame(bytes.NewReader(header[:])); err == nil {
		t.Fatal("accepted oversized payload")
	}
	if _, _, err := readDiscordFrame(bytes.NewReader(header[:3])); err == nil {
		t.Fatal("accepted truncated header")
	}
}

func TestDiscordBadHandshakeDoesNotOpenWindowsDiscord(t *testing.T) {
	guest, bridge := net.Pipe()
	defer guest.Close()
	called := make(chan bool, 1)
	done := make(chan error, 1)
	go func() {
		done <- relayDiscord(bridge, func() (io.ReadWriteCloser, error) { called <- true; return nil, io.EOF })
	}()
	writeDiscordFrame(guest, 1, []byte(`{"cmd":"AUTHORIZE"}`))
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("accepted bad handshake")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("relay blocked")
	}
	select {
	case <-called:
		t.Fatal("opened Discord before validation")
	default:
	}
}
