package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

type discordPipe struct{ *os.File }

func (pipe discordPipe) Write(data []byte) (int, error) {
	if err := pipe.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return 0, err
	}
	return pipe.File.Write(data)
}

func dialDiscordPipe() (io.ReadWriteCloser, error) {
	for index := 0; index < 10; index++ {
		name := fmt.Sprintf(`\\.\pipe\discord-ipc-%d`, index)
		if pipe, err := openDiscordPipe(name); err == nil {
			return pipe, nil
		}
	}
	return nil, fmt.Errorf("Windows Discord IPC is unavailable; open the Discord desktop app")
}

func openDiscordPipe(name string) (io.ReadWriteCloser, error) {
	path, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(path, syscall.GENERIC_READ|syscall.GENERIC_WRITE, 0, nil,
		syscall.OPEN_EXISTING, syscall.FILE_FLAG_OVERLAPPED, 0)
	if err != nil {
		return nil, err
	}
	return discordPipe{os.NewFile(uintptr(handle), name)}, nil
}

func startDiscordBridge() func() {
	listener, err := net.Listen("tcp4", fmt.Sprintf("127.0.0.1:%d", discordBridgePort))
	if err != nil {
		logf("Discord bridge could not start: %v", err)
		return func() {}
	}
	logf("Discord Rich Presence bridge listening on %d", discordBridgePort)
	return serveDiscordBridge(listener, dialDiscordPipe)
}
