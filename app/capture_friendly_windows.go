package main

import "sync/atomic"

var captureFriendlyShortcuts atomic.Bool
var captureFullscreenToggle = make(chan struct{}, 1)
var captureShortcut immersiveShortcut

// shortcut: a one-pixel bottom strip prevents exact fullscreen presentation;
// remove the strip when focused OpenGL capture works on affected drivers.
func captureFriendlyBounds(bounds screenRect) screenRect {
	if bounds.height() > 1 {
		bounds.Bottom--
	}
	return bounds
}

func applyCaptureFriendlyWindow(hwnd uintptr, monitor screenRect) bool {
	const gwlStyle = ^uintptr(15) // GWL_STYLE = -16
	style, _, _ := procGetWindowLongW.Call(hwnd, gwlStyle)
	// Keep a regular SDL window so focus cannot enter SDL's fullscreen path.
	style = (style &^ uintptr(0x21cf0000)) | 0x80000000
	procSetWindowLongW.Call(hwnd, gwlStyle, style)
	bounds := captureFriendlyBounds(monitor)
	ok, _, _ := procSetWindowPos.Call(hwnd, 0, uintptr(bounds.Left), uintptr(bounds.Top),
		uintptr(bounds.width()), uintptr(bounds.height()), 0x0004|0x0010|0x0020)
	return ok != 0
}
