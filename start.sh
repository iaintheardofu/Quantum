#!/bin/bash

# =====================================================================
# SomaOS: Unified System Launcher (Flutter Web Edition)
# =====================================================================
# This script launches both the Go FPGA Hardware Simulator 
# and the Flutter Web Visualizer concurrently.

echo "========================================================"
echo "    SomaOS v3.5: Initializing Visualization Environment "
echo "========================================================"

# Trap SIGINT (Ctrl+C) to clean up both background processes
trap "echo 'Shutting down SomaOS...'; pkill -P $$; exit" SIGINT SIGTERM

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 1. Start the Go Hardware Server
echo ">> [1/2] Booting Go FPGA Hardware Proxy on port 8081..."
cd "$PROJECT_ROOT/SomaServer"
go run main.go &
GO_PID=$!

# Wait a moment for the server to initialize
sleep 2

# 2. Start the Flutter Web Frontend
echo ">> [2/2] Booting Flutter Web Visualizer..."
cd "$PROJECT_ROOT/SomaUI/soma_flutter"
flutter run -d web-server --web-port 5173 &
FLUTTER_PID=$!

echo "========================================================"
echo " SUCCESS: SomaOS Environment is Live!"
echo " 👉 Visualizer: http://localhost:5173"
echo " 👉 API:        http://localhost:8081/api/state"
echo " "
echo " Press Ctrl+C to terminate both servers."
echo "========================================================"

# Keep the script running to hold the trap active
wait $GO_PID $FLUTTER_PID
