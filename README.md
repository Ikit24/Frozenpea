# Frozenpea

A productivity app which locks the user out from the machine after a set of minutes to prevent RSI related problems and support healthier work-life.

## Features

- Notifaction sound upon startup and before break starts
- Automatic screen lock after the set amount of minutes
- Background window display during rest periods
- Helps prevent repetitive strain injuries (RSI)
- Promotes healthier work habits

## System Requirements
- Linux with X11 display server
- Go 1.21 or higher
- Pop!_OS, Ubuntu, Debian, Fedora (X11 session), or other X11-based distributions

**Note:** Wayland is not currently supported. Ensure you're running an X11 session.

## Prerequisites

- Go 1.21+ installed: https://go.dev/dl/
- Required dependencies:
```bash
  go get fyne.io/fyne/v2
  go get github.com/BurntSushi/xgb
```

## Installation

1. **Clone the repository**:
```bash
   git clone https://github.com/Ikit24/Frozenpea.git
   cd Frozenpea
```

2. **Install dependencies**:
```bash
   go mod download
```

3. **Run the application**:
```bash
   go run .
```

## Usage

On startup you will have the option to set break and workduration.
Once the settings are confirmed the application starts, you will be notified when the break is about to happen.
After break starts you cannot interact with your machine.

## Automatic Startup Configuration (Linux)

Follow these steps to make FrozenPea run automatically on system startup.

### 1. Compile the Program
```bash
cd ~/projects/code/FrozenPea
go build -o main main.go
chmod +x main
```

### 2. Create Startup Script
```bash
nvim ~/startup.sh
```

Add the following content (replace `ati` with your username):
```bash
#!/bin/bash
/home/ati/projects/code/FrozenPea/main >> /home/ati/projects/code/FrozenPea/startup.log 2>&1
```

Save and exit (Ctrl+X, Y, Enter)

### 3. Make Script Executable
```bash
chmod +x ~/startup.sh
```

### 4. Add to Crontab
```bash
crontab -e
```

Add this line at the end:
```
@reboot /home/ati/startup.sh
```

Save and exit.

### 5. Verify Setup

Check that crontab was saved:
```bash
crontab -l
```

You should see: `@reboot /home/ati/startup.sh`

### Testing

Test the script without rebooting:
```bash
~/startup.sh
```

Check the log file for any errors:
```bash
cat ~/projects/code/FrozenPea/startup.log
```

### Managing the Autostart

**Stop the program:**
```bash
pkill main
```

**Remove from startup:**
```bash
crontab -e
# Delete the @reboot line, then save
```
