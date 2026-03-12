#!/bin/bash

# =====================================================================
# SomaOS: Unified System Launcher (HPQC Master Logging Edition)
# =====================================================================
# This script launches the Go Hardware Proxy, the Flutter Web UI, 
# and the Central WebSocket Logger concurrently.

# Trap SIGINT (Ctrl+C) to clean up all background processes
cleanup() {
    echo "Shutting down SomaOS..."
    pkill -P $$
    exit
}
trap cleanup SIGINT SIGTERM

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
LOG_FILE="$PROJECT_ROOT/master_execution.log"

# Clear old logs
echo ">> [INIT] Resetting Master Execution Log: $LOG_FILE"
echo "SomaOS Session Start: $(date)" > "$LOG_FILE"

# 1. Start the WebSocket Master Logger
echo ">> [1/3] Booting WebSocket Logger on port 8082..."
cd "$PROJECT_ROOT/logger"
go run main.go >> "$LOG_FILE" 2>&1 &
sleep 2

# 2. Start the Go Hardware Server
echo ">> [2/3] Booting Go FPGA Hardware Proxy on port 8081..."
cd "$PROJECT_ROOT/SomaServer"
stdbuf -oL go run main.go 2>&1 | while read line; do echo "[HARDWARE] $line" >> "$LOG_FILE"; done &
sleep 2

# 3. Start the SomaAI Cortex Router (Vertex AI Multimodal)
echo ">> [3/4] Booting SomaAI Cortex Router on port 8083..."
cd "$PROJECT_ROOT"
source SomaAI/.venv/bin/activate
stdbuf -oL python3 SomaAI/src/router.py 2>&1 | while read line; do echo "[AI-CORTEX] $line" >> "$LOG_FILE"; done &
sleep 2

# 4. Start the Flutter Web Frontend (Serving Release Build)
echo ">> [4/4] Booting Flutter Web Visualizer on port 5173..."
cd "$PROJECT_ROOT/SomaUI/soma_flutter/build/web"
# Ensure the server runs from the correct build directory
python3 -u -m http.server 5173 2>&1 | while read line; do echo "[WEB-SERVER] $line" >> "$LOG_FILE"; done &

echo "========================================================"
echo " SUCCESS: SomaOS HPQC Environment is Live!"
echo " 👉 Visualizer: http://localhost:5173"
echo " 👉 Master Log: master_execution.log"
echo "========================================================"

# Wait for background processes to finish
wait
