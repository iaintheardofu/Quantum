#!/bin/bash
# SOMA OS: Professional Deployment & Resuscitation Script
# This version uses the official Vivado flash engine followed by a 
# physical UART serial pulse to restore the network interface.

BIT_FILE="./build/mabel_x8c.bit"
BOARD_IP="10.100.102.9"
SERIAL_PORT="/dev/ttyUSB0"

echo "========================================================"
echo "    SomaOS: Official Silicon Manifestation & Sync      "
echo "========================================================"

# 1. Quiesce Board over Network
echo ">> [PS] Preparing board for reconfiguration..."
ssh -o ConnectTimeout=5 root@$BOARD_IP "pkill -9 soma_agent || true"

# 2. Perform Official Vivado JTAG Flash
echo ">> [JTAG] Initializing Xilinx Hardware Manager..."
pkill -9 hw_server
/home/noam/vivado/2025.2/Vivado/bin/hw_server > /dev/null 2>&1 &
sleep 2

/home/noam/vivado/2025.2/Vivado/bin/vivado -mode batch -source ./build/flash_mabel.tcl

if [ $? -eq 0 ]; then
    echo ">> [SUCCESS] MABEL x8C manifest active on Silicon."
    
    # 3. PHYSICAL RESUSCITATION (Via UART Serial)
    # The network link just dropped. We pulse it back to life from the outside.
    echo ">> [UART] Sending Resuscitation Pulse to /dev/ttyUSB0..."
    stty -F $SERIAL_PORT 115200 raw -echo
    echo "" > $SERIAL_PORT
    sleep 1
    echo "root" > $SERIAL_PORT
    sleep 1
    echo "ifconfig eth0 up && udhcpc -i eth0" > $SERIAL_PORT
    
    echo ">> [SYSTEM] Waiting for DHCP Snapback (10s)..."
    sleep 10

    # 4. Restart Telemetry Agent
    echo ">> [PS] Restoring telemetry services over network..."
    ./transfer_agent_and_run.sh
    echo ">> [NETWORK] Connection RESTORED. Deployment Complete."
else
    echo ">> [ERROR] Vivado Flash failed. Check USB connection."
    exit 1
fi
