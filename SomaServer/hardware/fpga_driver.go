package hardware

import (
        "math/rand"
        "fmt"
)

// FPGADriver simulates the 64-qubit "Entanglement Station" Grid.
type FPGADriver struct {
        // We use a 64-bit bitmask to represent 8 clusters of 8 qubits
        Register    uint64 
        Phase       float64
        ActiveCells int 
        RoutingMode string // "idle", "grover", "shor", "bell", "station"
        IsStation   bool   // Toggle between Single-Cube (8) and Station-Grid (64)
}

func NewFPGADriver() *FPGADriver {
        return &FPGADriver{
                Phase:       0.0,
                ActiveCells: 8,
                RoutingMode: "idle",
                IsStation:   false,
        }
}

func (f *FPGADriver) Poll() {
        f.Phase += 0.1
        if f.Phase > 6.28 { f.Phase = 0.0 }

        noise := rand.Float64() * 0.2
        master_state := (f.Phase + noise) > 3.14

        if f.IsStation {
                // STATION MODE: 64-bit Entangled Grid
                // Every bit in the 64-bit register mirrors the Master Station Anchor
                if master_state {
                        f.Register = 0xFFFFFFFFFFFFFFFF
                } else {
                        f.Register = 0x0000000000000000
                }
        } else {
                // SINGLE CUBE MODE: 8-bit GHZ State
                if master_state {
                        f.Register = 0x00000000000000FF
                } else {
                        f.Register = 0x0000000000000000
                }
        }
}

func (f *FPGADriver) SetRoutingMode(mode string) {
        f.RoutingMode = mode
        if mode == "station" {
                f.IsStation = true
                f.ActiveCells = 64
        } else {
                f.IsStation = false
                f.ActiveCells = 8
        }
        fmt.Printf("[HARDWARE RECONFIG] Silicon routing set to: %s (Station: %v)\n", mode, f.IsStation)
}

func (f *FPGADriver) TriggerDPR(action string) {
        if action == "spawn" && f.ActiveCells < 64 { f.ActiveCells += 8 } 
        if action == "collapse" && f.ActiveCells > 8 { f.ActiveCells -= 8 }
}

func (f *FPGADriver) GetHardwareData() HardwareState {
        return HardwareState{
                Register:    f.Register,
                ThermalLoad: rand.Float64()*5.0 + 35.0,
                PhaseField: f.Phase,
                ActiveCells: f.ActiveCells,
                RoutingMode: f.RoutingMode,
        }
}

type HardwareState struct {
        Register    uint64  `json:"register"`    // 64-bit quantum state
        ThermalLoad float64 `json:"thermal_load"`
        PhaseField  float64 `json:"phase_field"`
        ActiveCells int     `json:"active_cells"`
        RoutingMode string  `json:"routing_mode"`
}
