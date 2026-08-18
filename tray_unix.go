//go:build !windows

package main

import "fmt"

// setupSystemTray is a no-op on Linux to prevent GTK main loop conflicts with Wails
func setupSystemTray(a *App) {
	fmt.Println("System tray disabled on Linux during development.")
}

// HideWindow is mocked on Linux so the app doesn't disappear permanently
func (a *App) HideWindow() {
	fmt.Println("Hide to tray ignored on Linux (no tray available).")
}