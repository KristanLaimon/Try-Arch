package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestDiscordNativeNamedPipe(t *testing.T) {
	name := fmt.Sprintf(`\\.\pipe\tryarch-discord-test-%d`, os.Getpid())
	path, _ := syscall.UTF16PtrFromString(name)
	kernel := syscall.NewLazyDLL("kernel32.dll")
	handle, _, err := kernel.NewProc("CreateNamedPipeW").Call(uintptr(unsafe.Pointer(path)),
		3, 0, 1, 65536, 65536, 0, 0)
	if handle == uintptr(syscall.InvalidHandle) {
		t.Fatal(err)
	}
	defer syscall.CloseHandle(syscall.Handle(handle))
	pipe, err := openDiscordPipe(name)
	if err != nil {
		t.Fatal(err)
	}
	defer pipe.Close()
	done := make(chan error, 1)
	go func() {
		result, _, connectErr := kernel.NewProc("ConnectNamedPipe").Call(handle, 0)
		if result == 0 && connectErr != syscall.Errno(535) {
			done <- connectErr
			return
		}
		var data [4]byte
		var count uint32
		readErr := syscall.ReadFile(syscall.Handle(handle), data[:], &count, nil)
		if readErr == nil && (count != 4 || string(data[:]) != "ping") {
			readErr = io.ErrUnexpectedEOF
		}
		if readErr == nil {
			readErr = syscall.WriteFile(syscall.Handle(handle), []byte("pong"), &count, nil)
		}
		done <- readErr
	}()
	if _, err := pipe.Write([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	if err := pipe.(discordPipe).SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	var response [4]byte
	if _, err := io.ReadFull(pipe, response[:]); err != nil || string(response[:]) != "pong" {
		t.Fatalf("%s %v", response, err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}
