package hardware

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// FPGADriver handles both Simulated and Live hardware links.
type FPGADriver struct {
	Register    uint64
	Phase       float64
	ActiveCells int
	RoutingMode string
	IsStation   bool
	LiveMode    bool   // Toggle for physical ALINX connection
	BoardIP     string // ALINX board network address
}

func NewFPGADriver() *FPGADriver {
	return &FPGADriver{
		Phase:       0.0,
		ActiveCells: 8,
		RoutingMode: "idle",
		IsStation:   false,
		LiveMode:    false, // Default to simulation
		BoardIP:     "192.168.1.10", // Example ALINX IP
	}
}

func (f *FPGADriver) Poll() {
	if f.LiveMode {
		f.pollLiveBoard()
	} else {
		f.pollSimulation()
	}
}

func (f *FPGADriver) pollSimulation() {
	f.Phase += 0.1
	if f.Phase > 6.28 { f.Phase = 0.0 }
	noise := rand.Float64() * 0.2
	master_state := (f.Phase + noise) > 3.14

	if f.IsStation {
		if master_state { f.Register = 0xFFFFFFFFFFFFFFFF } else { f.Register = 0 }
	} else {
		if master_state { f.Register = 0xFF } else { f.Register = 0 }
	}
}

func (f *FPGADriver) pollLiveBoard() {
	// Fetching telemetry from the C-Agent running on the Zynq ARM core
	client := http.Client{Timeout: 50 * time.Millisecond}
	resp, err := client.Get(fmt.Sprintf("http://%s:8080/telemetry", f.BoardIP))
	if err != nil {
		fmt.Printf("[ALINX ERROR] Failed to reach board at %s: %v\n", f.BoardIP, err)
		return
	}
	defer resp.Body.Close()

	var boardData struct {
		Register uint64  `json:"reg"`
		Temp     float64 `json:"temp"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&boardData); err == nil {
		f.Register = boardData.Register
		// SPHY Engine synchronization logic
		f.Phase += 0.1 
	}
}

func (f *FPGADriver) SetRoutingMode(mode string) {
	f.RoutingMode = mode
	f.IsStation = (mode == "station")
	f.ActiveCells = 8
	if f.IsStation { f.ActiveCells = 64 }
	fmt.Printf("[HARDWARE RECONFIG] Silicon routing set to: %s\n", mode)
}

func (f *FPGADriver) TriggerDPR(action string) {
	if action == "spawn" && f.ActiveCells < 64 { f.ActiveCells += 8 } 
	if action == "collapse" && f.ActiveCells > 8 { f.ActiveCells -= 8 }
}

func (f *FPGADriver) GetHardwareData() HardwareState {
	return HardwareState{
		Register:    f.Register,
		ThermalLoad: rand.Float64()*5.0 + 35.0,
		PhaseField:  f.Phase,
		ActiveCells: f.ActiveCells,
		RoutingMode: f.RoutingMode,
		LiveMode:    f.LiveMode,
	}
}

type HardwareState struct {
	Register    uint64  `json:"register"`
	ThermalLoad float64 `json:"thermal_load"`
	PhaseField  float64 `json:"phase_field"`
	ActiveCells int     `json:"active_cells"`
	RoutingMode string  `json:"routing_mode"`
	LiveMode    bool    `json:"live_mode"`
}
