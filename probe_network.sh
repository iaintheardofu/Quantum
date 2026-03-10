#!/bin/bash
# SomaOS: Network Probe for ALINX 7020
# Searches for the board by scanning the ARP table or network for the ALINX OUI.

echo ">> Probing network for ALINX hardware..."

# ALINX boards often use specific MAC ranges. If we don't have the exact MAC, 
# we scan the local subnet.
LOCAL_SUBNET=$(ip route | grep eth0 | awk '{print $1}' | head -n 1)

if [ -z "$LOCAL_SUBNET" ]; then
    LOCAL_SUBNET=$(ip route | grep wlan0 | awk '{print $1}' | head -n 1)
fi

echo ">> Scanning subnet: $LOCAL_SUBNET"

# Use nmap to find active hosts and grep for known ALINX/Zynq fingerprints
# or simply list everything for the user to identify.
nmap -sn "$LOCAL_SUBNET" | grep -B 2 "MAC Address"

echo "--------------------------------------------------------"
echo ">> Current ARP Cache (Filtered for likely candidates):"
arp -a | grep -i "00:0a:35" # Xilinx/Alinx OUI prefix example
