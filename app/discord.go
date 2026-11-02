package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const discordFrameLimit = 64 << 10
const discordEnabledFilename = "discord-rpc-enabled"

// A separate opt-in marker keeps settings readable by older launchers.
func discordBridgeEnabled(dir string) bool {
	info, err := os.Lstat(filepath.Join(dir, discordEnabledFilename))
	return err == nil && info.Mode().IsRegular()
}

func readDiscordFrame(reader io.Reader) (uint32, []byte, error) {
	var header [8]byte
	if _, err := io.ReadFull(reader, header[:]); err != nil {
		return 0, nil, err
	}
	opcode, size := binary.LittleEndian.Uint32(header[:4]), binary.LittleEndian.Uint32(header[4:])
	if opcode > 4 || size > discordFrameLimit {
		return 0, nil, fmt.Errorf("invalid Discord frame")
	}
	data := make([]byte, int(size))
	_, err := io.ReadFull(reader, data)
	return opcode, data, err
}

func writeDiscordFrame(writer io.Writer, opcode uint32, data []byte) error {
	frame := make([]byte, 8+len(data))
	binary.LittleEndian.PutUint32(frame[:4], opcode)
	binary.LittleEndian.PutUint32(frame[4:8], uint32(len(data)))
	copy(frame[8:], data)
	_, err := io.Copy(writer, bytes.NewReader(frame))
	return err
}

// Only Rich Presence crosses this bridge; account, voice and authorization RPCs
// are deliberately unavailable to the guest.
func validDiscordRequest(opcode uint32, data []byte, handshake bool) bool {
	if handshake {
		var request struct {
			Version int    `json:"v"`
			ID      string `json:"client_id"`
		}
		if opcode != 0 || json.Unmarshal(data, &request) != nil || request.Version != 1 || len(request.ID) < 1 || len(request.ID) > 20 {
			return false
		}
		for _, digit := range request.ID {
			if digit < '0' || digit > '9' {
				return false
			}
		}
		return true
	}
	if opcode == 2 || opcode == 3 || opcode == 4 {
		return true
	}
	var request struct {
		Command string `json:"cmd"`
	}
	return opcode == 1 && json.Unmarshal(data, &request) == nil && request.Command == "SET_ACTIVITY"
}

func discordActivityForHost(data []byte) ([]byte, error) {
	var request map[string]json.RawMessage
	if err := json.Unmarshal(data, &request); err != nil {
		return nil, err
	}
	var args map[string]json.RawMessage
	if err := json.Unmarshal(request["args"], &args); err != nil {
		return nil, err
	}
	if args == nil {
		return nil, fmt.Errorf("missing Discord activity arguments")
	}
	// A Linux PID does not identify a process on Windows.
	args["pid"], _ = json.Marshal(os.Getpid())
	var activity map[string]json.RawMessage
	if raw := args["activity"]; len(raw) != 0 && !bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if err := json.Unmarshal(raw, &activity); err != nil {
			return nil, err
		}
		// Guest apps can put the Linux account name in their location/status.
		activity["state"], _ = json.Marshal("In arch-linux")
		args["activity"], _ = json.Marshal(activity)
	}
	request["args"], _ = json.Marshal(args)
	return json.Marshal(request)
}

func relayDiscord(guest net.Conn, dial func() (io.ReadWriteCloser, error)) error {
	defer guest.Close()
	_ = guest.SetReadDeadline(time.Now().Add(5 * time.Second))
	opcode, data, err := readDiscordFrame(guest)
	if err != nil || !validDiscordRequest(opcode, data, true) {
		return fmt.Errorf("invalid Discord handshake")
	}
	pipe, err := dial()
	if err != nil {
		return err
	}
	defer pipe.Close()
	if err := writeDiscordFrame(pipe, opcode, data); err != nil {
		return err
	}
	_ = guest.SetReadDeadline(time.Time{})
	done := make(chan error, 1)
	go func() {
		defer guest.Close()
		for {
			opcode, data, err := readDiscordFrame(pipe)
			if err == nil {
				_ = guest.SetWriteDeadline(time.Now().Add(5 * time.Second))
				err = writeDiscordFrame(guest, opcode, data)
			}
			if err != nil {
				done <- err
				return
			}
		}
	}()
	for {
		opcode, data, err := readDiscordFrame(guest)
		if err == nil && !validDiscordRequest(opcode, data, false) {
			err = fmt.Errorf("unsupported Discord RPC command")
		}
		if err == nil && opcode == 1 {
			data, err = discordActivityForHost(data)
		}
		if err == nil {
			err = writeDiscordFrame(pipe, opcode, data)
		}
		if err != nil {
			_ = guest.Close()
			_ = pipe.Close()
			<-done
			return err
		}
	}
}

func serveDiscordBridge(listener net.Listener, dial func() (io.ReadWriteCloser, error)) func() {
	var mu sync.Mutex
	clients := make(map[net.Conn]bool)
	stopped := false
	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			if stopped || len(clients) >= 16 {
				mu.Unlock()
				connection.Close()
				continue
			}
			clients[connection] = true
			mu.Unlock()
			go func() {
				if err := relayDiscord(connection, dial); err != nil && err != io.EOF {
					logf("Discord bridge: %v", err)
				}
				mu.Lock()
				delete(clients, connection)
				mu.Unlock()
			}()
		}
	}()
	return func() {
		mu.Lock()
		defer mu.Unlock()
		stopped = true
		listener.Close()
		for connection := range clients {
			connection.Close()
		}
	}
}
