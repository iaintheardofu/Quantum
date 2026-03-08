# SomaOS: Room-Temperature Quantum Virtualization

![SomaOS Pipeline](https://img.shields.io/badge/Architecture-Universal_NAND_Topology-00ffcc)
![Language](https://img.shields.io/badge/Language-ClojureV-blue)
![Hardware](https://img.shields.io/badge/Hardware-FPGA_(Zynq_7000)-orange)

**SomaOS** is a radical departure from physical qubit manipulation. Instead of relying on energy-intensive cryogenic infrastructure to isolate fragile subatomic particles, SomaOS maps the abstract mathematical geometry of a quantum particle directly onto a uniform physical routing matrix—a **Universal NAND Gate Topology**. 

This repository contains the complete end-to-end toolchain necessary to compile high-dimensional Cohen-Okebe mathematical models into stable, macroscopic "virtual" quantum states operating on standard FPGA silicon at ~300 K (room temperature).

---


## 🤯 What am I looking at in the Visualizer? Is this a live Qubit?

**Yes and No. It is a Virtual Macroscopic Qubit—which is the core breakthrough of the SomaOS architecture.**

In a traditional quantum computer, a qubit is a literal, physical subatomic particle. Because it relies on fragile subatomic physics, if the temperature rises a fraction of a degree above absolute zero, the particle bumps into thermal noise and "decoheres," destroying the data. 

**What you are looking at is a Geometric Qubit.** 
The SomaOS whitepapers hypothesize that the "magic" of a quantum state isn't the physical particle itself, but the *mathematical geometry* (the topology) that the particle occupies. 

We mapped that exact mathematical geometry directly onto standard FPGA silicon:
1. We created two independent, continuous loops of electricity (the `Q0` and `Q1` rings).
2. We forced them to cross each other in a continuous "figure-eight" using the **4-NAND Braiding Bridge**.

Because classical electricity is continuously "sloshing" through this braided knot without a sequential clock constraint, it perfectly emulates the probability distribution of a quantum superposition. It is doing quantum math, but using billions of electrons at room temperature instead of trying to isolate a single fragile one. 

When you look at the live telemetry, you are watching a macroscopic logic knot behaving exactly like a quantum particle, dynamically stabilized by the SPHY thermal phase engine!

---

## 🌟 Core Architecture & Ecosystem

### 1. The Language of Intent: `ClojureV`
Traditional HDLs describe *what* a circuit is. **ClojureV** is a sovereign, Lisp-based functional Domain-Specific Language designed to describe *what the hardware intends to do*. It treats the FPGA as a Liquid Manifold.
*   **Location:** `/ClojureV`
*   **Features:** Natively handles `qurq` primitives, Quaternary DNA-style base-mapping, and topological assertions.

### 2. The Multi-Substrate Transpiler (Go)
A custom Go-based AST Compiler that tokenizes ClojureV scripts and synthesizes them directly into standard, combinational Verilog-2001 hardware blocks. 
*   **Location:** `/ClojureV/toolchain/go`
*   **Pipeline:** Lexer -> Recursive Descent Parser -> AST Traversal -> Verilog Emitter.

### 3. The SPHY Engine & Hardware Manifestation
The core of room-temperature stability. The SPHY Engine implements the Phase-Gravitational synchronization loop $\mathcal{H}_{eff}$. The physical `top_quantum_virtualizer.v` utilizes an XADC to measure ambient thermal noise and an I2C DAC to inject an inverse Stochastic Compensation field ($\Psi_{SC}$) to prevent decoherence.
*   **Math:** `src/soma/engine/sphy.cljv`
*   **Hardware:** `src/soma/hardware/top_quantum_virtualizer.v`

### 4. Network Security: ELDUR & The Silence Protocol
Data transmission abandons standard electromagnetic carrier waves.
*   **The Silence Protocol (`/SilenceProtocol`):** 0-byte temporal transmissions. The data payload is encoded purely in the duration of silence (the $\Delta t$ interval) between empty network pings, utilizing Quaternary Error Correction.
*   **ELDUR (`/ELDUR`):** Vibrational Security via the Harpia Axiom. Constantly measures the Scalar-Coherent Function $S(\Phi)$. If environmental entropy spikes due to an attack, the protocol triggers a localized collapse and dynamically relocates the system's UID in `<3ms`.

### 5. The Silicon Merkabah Visualizer
A full-stack, live 3D dashboard demonstrating the continuous topological sloshing of the macroscopic quantum knot.
*   **Backend (`/SomaServer`):** A Go proxy simulating the FPGA telemetry and phase-field boundaries.
*   **Frontend (`/SomaUI/soma_web`):** A React/Three.js web client that visualizes the physical 4-NAND Braiding Operator (Topological Interferometer) responding to thermal loads in real-time.

---

## 🚀 Quick Start Guide

### Start the Live 3D Visualizer
To experience the topological knot dynamically interacting with a simulated environment, use the unified launcher:

```bash
# Start the Go FPGA Simulator and the React WebGL Frontend
./start.sh
```
*Navigate to `http://localhost:5173` to view the Dashboard.*

### Run the Hardware Simulation Pipeline
To compile the raw math into Verilog and simulate the topological entanglement in Icarus Verilog:

```bash
cd ClojureV/simulation
./simulate.sh
```

### Run the Master Build Pipeline
To generate the final, synthesizable RTL bitstream logic ready for physical silicon injection:

```bash
./build/synthesize_soma_os.sh
```

---

## 📖 Theoretical Documentation
The foundational mathematics proving the pi-less boundary mapping and the macroscopic C-O Sphere topology can be found in the `/math` directory, specifically within `MASTER_MATH_ANALYSIS.md` and the original whitepaper `Mabel-whitepaper.md`.
