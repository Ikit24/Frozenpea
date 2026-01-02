# Frozenpea

A productivity app which locks the user out from the machine after a set of minutes to prevent RSI related problems and support healthier work.

## Features

- Automatic screen lock after 50 minutes (customizable)
- Background window display during rest periods
- Helps prevent repetitive strain injuries (RSI)
- Promotes healthier work habits

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

Once the application is running, you will be greeted with a background picture window when it is time to take a rest (after 50 minutes by default).

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
nano ~/startup.sh
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
