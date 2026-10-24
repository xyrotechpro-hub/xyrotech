//go:build windows

package main

import (
	"fmt"
	"syscall"
	"time"
)

const (
	smCXScreen     = 0
	smCYScreen     = 1
	mouseLeftDown  = 0x0002
	mouseLeftUp    = 0x0004
	mouseRightDown = 0x0008
	mouseRightUp   = 0x0010
	mouseWheel     = 0x0800
	keyeventKeyup  = 0x0002
	vkControl      = 0x11
	vkTab          = 0x09
	wheelDownDelta = ^uint32(119)
)

var (
	user32             = syscall.NewLazyDLL("user32.dll")
	getSystemMetrics   = user32.NewProc("GetSystemMetrics")
	getDoubleClickTime = user32.NewProc("GetDoubleClickTime")
	setCursorPos       = user32.NewProc("SetCursorPos")
	mouseEvent         = user32.NewProc("mouse_event")
	keybdEvent         = user32.NewProc("keybd_event")
)

func screenSize() (int, int) {
	width, _, _ := getSystemMetrics.Call(smCXScreen)
	height, _, _ := getSystemMetrics.Call(smCYScreen)
	return int(width), int(height)
}

func doubleClickInterval() int {
	interval, _, _ := getDoubleClickTime.Call()
	if interval == 0 {
		return 500
	}
	return int(interval)
}

func moveCursor(x, y int) error {
	result, _, err := setCursorPos.Call(uintptr(x), uintptr(y))
	if result == 0 {
		return fmt.Errorf("SetCursorPos failed: %w", err)
	}
	return nil
}

func performAction(action string) error {
	switch action {
	case "left":
		mouseEvent.Call(mouseLeftDown, 0, 0, 0, 0)
		mouseEvent.Call(mouseLeftUp, 0, 0, 0, 0)
	case "left-down":
		mouseEvent.Call(mouseLeftDown, 0, 0, 0, 0)
	case "left-up":
		mouseEvent.Call(mouseLeftUp, 0, 0, 0, 0)
	case "double-left":
		mouseEvent.Call(mouseLeftDown, 0, 0, 0, 0)
		mouseEvent.Call(mouseLeftUp, 0, 0, 0, 0)
		time.Sleep(90 * time.Millisecond)
		mouseEvent.Call(mouseLeftDown, 0, 0, 0, 0)
		mouseEvent.Call(mouseLeftUp, 0, 0, 0, 0)
	case "right":
		mouseEvent.Call(mouseRightDown, 0, 0, 0, 0)
		mouseEvent.Call(mouseRightUp, 0, 0, 0, 0)
	case "scroll-up":
		mouseEvent.Call(mouseWheel, 0, 0, 120, 0)
	case "scroll-down":
		mouseEvent.Call(mouseWheel, 0, 0, uintptr(wheelDownDelta), 0)
	case "tab":
		keybdEvent.Call(vkTab, 0, 0, 0)
		keybdEvent.Call(vkTab, 0, keyeventKeyup, 0)
	case "next-tab":
		keybdEvent.Call(vkControl, 0, 0, 0)
		keybdEvent.Call(vkTab, 0, 0, 0)
		keybdEvent.Call(vkTab, 0, keyeventKeyup, 0)
		keybdEvent.Call(vkControl, 0, keyeventKeyup, 0)
	default:
		return fmt.Errorf("unsupported action %q", action)
	}
	return nil
}
