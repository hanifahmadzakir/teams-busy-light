# Teams Busy Light 🚦

A desktop application and hardware combo that acts as a physical status indicator for Microsoft Teams. Built with [Wails](https://wails.io/) (Go + Vue 3) and an ESP8266 (Wemos D1 Mini), this project allows you to manually control your status light or automatically sync it with your Microsoft Teams presence.

## ✨ Features

- **MS Teams Auto-Sync:** Automatically reads your Microsoft Teams local log files to detect status changes (Available, Busy, Away, Do Not Disturb, In a Meeting, Ringing) and updates the physical light.
- **Manual Override:** A sleek UI to manually set the light to Green, Red, Yellow, or trigger a "Calling" blink animation.
- **OLED Display Support:** The hardware features support for an SSD1306 OLED screen to display your current status textually (e.g., "Available", "Busy / DND").
- **System Tray Integration:** Minimizes to the system tray for unobtrusive background operation (specifically supported on Windows).
- **Serial Monitor:** Built-in serial monitor in the desktop app UI for easy hardware debugging.

---

## 🏗️ Architecture

The project is divided into three main components:

1. **Desktop App (Go Backend):** Built using Wails. Handles system tray logic, serial port communication, and file system monitoring (polling the MS Teams local cache logs).
2. **Frontend UI (Vue 3 + Tailwind CSS):** Provides the graphical interface for manual control, configuration (COM port selection, Auto Mode toggle), and hardware monitoring. 
3. **Hardware (PlatformIO / C++):** Runs on a Wemos D1 Mini (ESP8266). Listens for commands over USB Serial and updates the LEDs and OLED display accordingly.

---

## 🛠️ Hardware Requirements & Wiring

*   **Microcontroller:** Wemos D1 Mini (ESP8266)
*   **Display:** SSD1306 OLED Display (I2C)
*   **LEDs:** Green, Yellow, and Red LEDs (with appropriate resistors)

**Wiring Map:**
*   `D5` -> Green LED
*   `D6` -> Yellow LED
*   `D7` -> Red LED
*   `I2C Pins (D1/D2)` -> SSD1306 OLED (SCL/SDA)

---

## 🚀 Getting Started

### Prerequisites

To build and run this project, you will need:
*   [Go](https://golang.org/doc/install) (1.25.0 or later)
*   [Node.js and npm](https://nodejs.org/) (for the Vue frontend)
*   [Wails CLI](https://wails.io/docs/gettingstarted/installation) (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
*   [PlatformIO](https://platformio.org/) (VS Code Extension recommended for compiling the hardware code)

### 1. Flash the Hardware (Wemos D1 Mini)
1. Open the `BusyLight/` folder in VS Code (with PlatformIO installed).
2. Connect your Wemos D1 Mini via USB.
3. Build and Upload the code using the PlatformIO toolbar. 

### 2. Run the Desktop Application (Development Mode)
1. Navigate to the root folder of the repository.
2. Run the Wails dev server:
   ```bash
   wails dev
   ```
   *This will compile the Go backend and start the Vite development server for the Vue frontend with hot-reload enabled.*

### 3. Build for Production
To build a standalone executable for your operating system:
```bash
wails build
```
The compiled binary will be located in the `build/bin/` directory.

---

## ⚙️ How Auto-Sync Works (Windows)

When "Teams Sync" is enabled in the app, the Go backend continuously polls the Microsoft Teams local log directory (usually located at `%localappdata%\Packages\MSTeams_8wekyb3d8bbwe\LocalCache\Microsoft\MSTeams\Logs`). 

It scans the most recent `.log` file for `StatusIndicatorStateService` entries to parse your real-time status and sends the corresponding color command (`Green`, `Red`, `Yellow`, or `BlinkRed`) via USB Serial to the microcontroller.

*Note: You can override the default log path in the "Teams Sync" tab of the application settings if your Teams installation differs.*

---

## 📜 License
This project utilizes various open-source packages. Please refer to individual components for specific licensing (e.g., SIL Open Font License for the included Nunito font).