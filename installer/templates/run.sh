#!/bin/bash
SCRIPT_DIR=$(dirname "$0")
BINARY_PATH="$SCRIPT_DIR/../Resources/$BinaryName"
CURRENT_DIR="$SCRIPT_DIR/../../../"

# Проверяем, существует ли бинарник
if [ ! -f "$BINARY_PATH" ]; then
    osascript -e 'display dialog "Error: Binary file not found." buttons {"OK"} default button "OK" with icon stop'
    exit 1
fi

# Запускаем AppleScript
osascript <<EOF
tell application "Terminal"
    -- Проверяем состояние Terminal
    set wasRunning to running
    set hadWindows to (exists window 1)
    
    -- Активируем Terminal
    activate
    
    -- Создаём новое окно
    if wasRunning and hadWindows then
        tell application "System Events" to keystroke "n" using command down
        delay 0.1
        set targetWindow to window 1
        do script "cd " & quoted form of "$CURRENT_DIR" & " && clear && " & quoted form of "$BINARY_PATH" in targetWindow
    else
        set targetWindow to window 1
        do script "cd " & quoted form of "$CURRENT_DIR" & " && clear && " & quoted form of "$BINARY_PATH" in targetWindow
    end if
    
    -- Настраиваем окно
    set custom title of targetWindow to "$BinaryName"
    set bounds of targetWindow to {100, 100, 800, 500}
    set frontmost of targetWindow to true
    
    -- Ждём завершения программы
    repeat while busy of targetWindow
        delay 0.5
    end repeat
    
    -- Закрываем окно
    close targetWindow
    
    -- Если Terminal не был открыт до запуска и нет других окон, завершаем его
    if not wasRunning and (count windows) is 0 then
        quit
    end if
end tell

-- Возвращаемся в Finder
tell application "Finder" to activate
EOF