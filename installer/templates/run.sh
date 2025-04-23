#!/bin/bash
SCRIPT_DIR=$(dirname "$0")
BINARY_PATH="$SCRIPT_DIR/../Resources/$BinaryName"
CURRENT_DIR="$SCRIPT_DIR/../../../"
REQUIRES_SUDO=$Sudo

if [ ! -f "$BINARY_PATH" ]; then
    osascript -e 'display dialog "Error: Binary file not found." buttons {"OK"} default button "OK" with icon stop'
    exit 1
fi

ESCAPED_CURRENT_DIR=$(printf '%q' "$CURRENT_DIR")
ESCAPED_BINARY_PATH=$(printf '%q' "$BINARY_PATH")

if [ "$REQUIRES_SUDO" = "YES" ]; then
    TERMINAL_CMD="cd \"$ESCAPED_CURRENT_DIR\" && clear && \
    echo 'Enter your password to grant the application the necessary permissions for system-level changes.' && \
    sudo \"$ESCAPED_BINARY_PATH\""
else
    TERMINAL_CMD="cd \"$ESCAPED_CURRENT_DIR\" && clear && \
    \"$ESCAPED_BINARY_PATH\""
fi

ESCAPED_TERMINAL_CMD=$(echo "$TERMINAL_CMD" | sed 's/"/\\"/g')

osascript <<EOF
tell application "Terminal"
    set wasRunning to running
    set hadWindows to (exists window 1)
    activate
    if wasRunning and hadWindows then
        tell application "System Events" to keystroke "n" using command down
        delay 0.1
    end if
    
    set targetWindow to window 1
    do script "$ESCAPED_TERMINAL_CMD" in targetWindow
    
    try
        set custom title of targetWindow to "$BinaryName"
    end try
    set bounds of targetWindow to {100, 100, 800, 500}
    set frontmost of targetWindow to true
    
    repeat while busy of targetWindow
        delay 0.5
    end repeat
    
    try
        close targetWindow
    end try
    
    if not wasRunning and (count windows) is 0 then
        quit
    end if
end tell

Tell application "Finder" to activate
EOF