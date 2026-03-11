#!/bin/bash

# =====================================================================
# SomaOS: Unified System Launcher (HPQC Master Logging Edition)
# =====================================================================
# This script launches the Go Hardware Proxy, the Flutter Web UI, 
# and the Central WebSocket Logger concurrently.

echo "========================================================"
echo "    SomaOS v3.5: Initializing Master Execution Grid     "
echo "========================================================"

# Trap SIGINT (Ctrl+C) to clean up all background processes
trap "echo 'Shutting down SomaOS...'; pkill -P $$; exit" SIGINT SIGTERM

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$PROJECT_ROOT/master_execution.log"

# Clear old logs
echo ">> [INIT] Resetting Master Execution Log: $LOG_FILE"
echo "SomaOS Session Start: $(date)" > "$LOG_FILE"

# 1. Start the WebSocket Master Logger
echo ">> [1/3] Booting WebSocket Logger on port 8082..."
cd "$PROJECT_ROOT/logger"
go run main.go >> "$LOG_FILE" 2>&1 &
LOGGER_PID=$!
sleep 2

# 2. Start the Go Hardware Server (Tailing to Master Log)
echo ">> [2/3] Booting Go FPGA Hardware Proxy on port 8081..."
cd "$PROJECT_ROOT/SomaServer"
# We use a subshell to prefix hardware logs before appending to master log
( go run main.go | while read line; do echo "[HARDWARE] $line"; done ) >> "$LOG_FILE" 2>&1 &
GO_PID=$!
sleep 2

# 3. Start the Flutter Web Frontend
echo ">> [3/3] Booting Flutter Web Visualizer on port 5173..."
cd "$PROJECT_ROOT/SomaUI/soma_flutter"
( flutter run -d web-server --web-port 5173 | while read line; do echo "[WEB-BUILD] $line"; done ) >> "$LOG_FILE" 2>&1 &
FLUTTER_PID=$!

echo "========================================================"
echo " SUCCESS: SomaOS HPQC Environment is Live!"
echo " 👉 Visualizer: http://localhost:5173"
echo " 👉 Master Log: tail -f master_execution.log"
echo " "
echo " Press Ctrl+C to terminate the session."
echo "========================================================"

# Tailing the log to the terminal so the user can watch the master execution
tail -f "$LOG_FILE"

# Keep the script running
wait $LOGGER_PID $GO_PID $FLUTTER_PID
