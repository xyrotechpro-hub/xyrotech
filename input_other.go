//go:build !windows

package main

import "errors"

func screenSize() (int, int)   { return 0, 0 }
func doubleClickInterval() int { return 500 }
func moveCursor(int, int) error {
	return errors.New("system pointer control is currently available on Windows")
}
func performAction(string) error {
	return errors.New("system pointer actions are currently available on Windows")
}
