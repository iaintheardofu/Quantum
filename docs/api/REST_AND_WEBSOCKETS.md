# MABEL x8C: SomaOS API & Interface Guide

This document outlines the software and hardware interfaces for interacting with the SomaOS v4.3 ecosystem. It is intended for software engineers, systems integrators, and AI agents looking to query or reconfigure the High-Performance Quantum Computing (HPQC) grid.

## 1. The SomaServer (Go Hardware Proxy)
The Go backend acts as the secure interface between the user/UI and the physical `soma_agent` running on the ALINX 7020 Zynq processor.

**Default Address:** `http://localhost:8081`

### `GET /api/state`
Returns the real-time telemetry of the physical logic matrix.
**Response (JSON):**
```json
{
  "register": 12188089,
  "thermal_load": 36.4,
  "phase_field": 3.14159,
  "active_cells": 64,
  "routing_mode": "station",
  "live_mode": true,
  "winding_number": 1405,
  "decoherence_rate": 0.052,
  "compensation_vector": -0.078,
  "fidelity": 0.998
}
```

### `POST /api/reconfigure`
Triggers a physical logic reconfiguration and forces a hardware-level PL reset.
**Payload (JSON):**
```json
{ "mode": "grover" } // Options: "idle", "grover", "shor", "bell", "station"
```

### `POST /api/dpr`
Triggers Dynamic Partial Reconfiguration (DPR) to expand or collapse the topological dimensionality by 8-qubit increments.
**Payload (JSON):**
```json
{ "action": "spawn" } // Options: "spawn", "collapse"
```

### `POST /api/synthesize`
The core compiler bridge. Accepts raw ClojureV source code, transpiles it to Verilog, lints it, and triggers a hardware reconfiguration to match the intent.
**Payload (JSON):**
```json
{ 
  "code": "(ns ClojureV.qurq) (defn-ai oracle [in] (qurq/phi-scale out in -1.0))",
  "mode": "grover"
}
```

---

## 2. The SomaAI Cortex (Vertex AI Router)
The Python-based AI router provides multimodal analysis of the physical hardware environment.

**Default Address:** `http://localhost:8083`

### `POST /api/ai/vision`
Accepts a base64 encoded image frame (usually from the Sentinel Eye) and a natural language prompt. Uses Gemini 1.5 Pro to physically inspect the hardware LEDs and state.
**Payload (JSON):**
```json
{
  "image": "data:image/jpeg;base64,/9j/4AAQSkZJ...",
  "prompt": "Are the DONE and INIT_B LEDs illuminated?"
}
```

### `POST /api/ai/generate_code`
Passes a natural language engineering intent to the AI to generate syntactically correct ClojureV code for the IDE.
**Payload (JSON):**
```json
{
  "intent": "Write a topological operator that negates the input field."
}
```

### `WS /ws/telemetry_analysis`
A WebSocket connection that accepts streaming chunks of the `master_execution.log`. It buffers the logs and periodically responds with an AI-generated audit of thermal stability and protocol integrity.

---

## 3. The Master Execution Logger
A Go-based WebSocket hub that aggregates output from the hardware, the UI, and the transpiler.

**Default Address:** `ws://localhost:8082/log`
All output sent to this socket is simultaneously printed to the terminal and appended to `master_execution.log`.
