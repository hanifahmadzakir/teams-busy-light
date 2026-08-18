//go:build windows

package main

import (
	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func setupSystemTray(a *App) {
	go systray.Run(a.onReady, a.onExit)
}

// --- SYSTEM TRAY LOGIC ---
func (a *App) onReady() {
	systray.SetTitle("Busy Light")
	systray.SetTooltip("Teams Busy Light Controller")
	
	mOpen := systray.AddMenuItem("Open", "Open the application UI")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	for {
		select {
		case <-mOpen.ClickedCh:
			runtime.WindowShow(a.ctx)
		case <-mQuit.ClickedCh:
			systray.Quit()
			runtime.Quit(a.ctx)
			return
		}
	}
}

func (a *App) onExit() {
	// Clean up tray icon
}

// HideWindow allows the frontend to hide the window to the tray
func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
}