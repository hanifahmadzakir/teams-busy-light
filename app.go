package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	// "github.com/getlantern/systray"
	// "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.bug.st/serial"
)

type App struct {
	ctx          context.Context
	serialPort   serial.Port
	autoMode     bool
	logPath      string
	stopLogCheck chan bool
}

func NewApp() *App {
	return &App{
		stopLogCheck: make(chan bool),
	}
}

// OnStartup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize System Tray in a separate goroutine
	// go systray.Run(a.onReady, a.onExit)
	setupSystemTray(a)
}

// --- SYSTEM TRAY LOGIC ---
// func (a *App) onReady() {
// 	systray.SetTitle("Busy Light")
// 	systray.SetTooltip("Teams Busy Light Controller")
	
// 	mOpen := systray.AddMenuItem("Open", "Open the application UI")
// 	mQuit := systray.AddMenuItem("Quit", "Quit the application")

// 	for {
// 		select {
// 		case <-mOpen.ClickedCh:
// 			runtime.WindowShow(a.ctx)
// 		case <-mQuit.ClickedCh:
// 			systray.Quit()
// 			runtime.Quit(a.ctx)
// 			return
// 		}
// 	}
// }

// func (a *App) onExit() {
// 	// Clean up tray icon
// }

// --- IPC METHODS CALLED BY VUE ---

// GetPorts returns a list of available COM ports
func (a *App) GetPorts() []string {
    portMap := make(map[string]bool)

    // 1. Try getting ports from the serial library
    ports, err := serial.GetPortsList()
    if err == nil {
        for _, port := range ports {
            portMap[port] = true
        }
    }

    // 2. Scan OS-specific locations (like Windows Registry) as a fallback
    scanOSSpecificPorts(portMap)

    // 3. Convert the map back to a slice for the frontend
    var finalPorts []string
    for port := range portMap {
        finalPorts = append(finalPorts, port)
    }

    return finalPorts
}

// ConnectSerial connects to the selected Wemos D1 port
func (a *App) ConnectSerial(portName string) string {
	// Close any existing connection before opening a new one
	if a.serialPort != nil {
		a.serialPort.Close()
	}

	mode := &serial.Mode{BaudRate: 9600}
	port, err := serial.Open(portName, mode)
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	a.serialPort = port

	// Start a background goroutine to read incoming serial data
	go a.listenSerial()

	return "Connected to " + portName
}

// listenSerial continuously reads from the serial port and sends data to the frontend
func (a *App) listenSerial() {
	if a.serialPort == nil {
		return
	}
	
	// Create a scanner to read the serial stream line by line
	scanner := bufio.NewScanner(a.serialPort)
	for scanner.Scan() {
		line := scanner.Text()
		// Emit an event named "serial-data" to the Vue frontend
		runtime.EventsEmit(a.ctx, "serial-data", line)
	}
	
	// If the scanner stops (e.g., device unplugged)
	if err := scanner.Err(); err != nil {
		// Ignore the error if we intentionally closed the port to reconnect
		if !strings.Contains(err.Error(), "Port has been closed") {
			runtime.EventsEmit(a.ctx, "serial-data", fmt.Sprintf("Serial connection lost: %v", err))
		}
	}
}

// DisconnectSerial closes the active serial connection
func (a *App) DisconnectSerial() string {
	if a.serialPort != nil {
		err := a.serialPort.Close()
		a.serialPort = nil // Clear the reference
		if err != nil {
			return fmt.Sprintf("Error disconnecting: %v", err)
		}
		return "Disconnected"
	}
	return "Already disconnected"
}

// SendCommand sends a string (e.g., "Red", "Green") to the Wemos
func (a *App) SendCommand(cmd string) string {
	if a.serialPort == nil {
		return "Port not connected"
	}
	_, err := a.serialPort.Write([]byte(cmd + "\n"))
	if err != nil {
		return fmt.Sprintf("Error writing to serial: %v", err)
	}
	return "Success"
}

// SetAutoMode toggles reading the Teams Log
func (a *App) SetAutoMode(enabled bool, inputPath string) string {
	a.autoMode = enabled
	
	if !enabled {
		a.stopLogCheck <- true
		return "Manual mode enabled"
	}

	// Handle Windows %localappdata% resolution
	if inputPath == "" {
		inputPath = "%localappdata%\\Packages\\MSTeams_8wekyb3d8bbwe\\LocalCache\\Microsoft\\MSTeams\\Logs"
	}
	
	if strings.Contains(inputPath, "%localappdata%") {
		localAppData := os.Getenv("LOCALAPPDATA")
		inputPath = strings.Replace(inputPath, "%localappdata%", localAppData, 1)
	}

	a.logPath = inputPath
	
	go a.watchTeamsLog() // Start background watcher
	return "Auto mode enabled, watching logs..."
}

// --- LOG WATCHER (Background Goroutine) ---
func (a *App) watchTeamsLog() {
    ticker := time.NewTicker(3 * time.Second)
    defer ticker.Stop()

    lastStatus := ""

    for {
        select {
        case <-a.stopLogCheck:
            return
        case <-ticker.C:
            // Find the latest .log file in the directory
            files, err := filepath.Glob(filepath.Join(a.logPath, "*.log"))
            if err != nil || len(files) == 0 {
                continue
            }

            // Find the most recently modified file
            var newestFile string
            var newestTime int64
            for _, file := range files {
                info, err := os.Stat(file)
                if err == nil && info.ModTime().Unix() > newestTime {
                    newestTime = info.ModTime().Unix()
                    newestFile = file
                }
            }

            if newestFile != "" {
                currentStatus := a.readLatestStatus(newestFile)
                if currentStatus != "" && currentStatus != lastStatus {
                    a.SendCommand(currentStatus)
                    lastStatus = currentStatus
                }
            }
        }
    }
}

// Helper function to read the file and extract the status
func (a *App) readLatestStatus(filePath string) string {
    file, err := os.Open(filePath)
    if err != nil {
        return ""
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    latestColor := ""

    // Scan through the file. In a production app, you might want to 
    // seek to the end and tail it, but for small logs, scanning is fine.
    for scanner.Scan() {
        line := scanner.Text()
        
        // This keyword is based on standard MS Teams log structures
        if strings.Contains(line, "StatusIndicatorStateService: Added") {
            if strings.Contains(line, "Available") {
                latestColor = "Green"
            } else if strings.Contains(line, "Busy") || strings.Contains(line, "InAMeeting") || strings.Contains(line, "OnThePhone") || strings.Contains(line, "DoNotDisturb") || strings.Contains(line, "Presenting") {
                latestColor = "Red"
            } else if strings.Contains(line, "BeRightBack") || strings.Contains(line, "Away") || strings.Contains(line, "Offline") {
                latestColor = "Yellow"
            }
        } else if strings.Contains(line, "reportIncomingCall") {
            latestColor = "BlinkRed"
        }
    }
    return latestColor
}

// HideWindow allows the frontend to hide the window to the tray
// func (a *App) HideWindow() {
// 	runtime.WindowHide(a.ctx)
// }