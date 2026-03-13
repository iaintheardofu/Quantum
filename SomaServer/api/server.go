package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"soma_server/hardware"
	"time"
)

type Server struct {
	driver *hardware.FPGADriver
}

func NewServer(driver *hardware.FPGADriver) *Server {
	return &Server{driver: driver}
}

func (s *Server) Start(addr string) error {
	go func() {
		for {
			s.driver.Poll()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/dpr", s.handleDPR)
	mux.HandleFunc("/api/reconfigure", s.handleReconfigure)
	mux.HandleFunc("/api/synthesize", s.handleSynthesize)

	return http.ListenAndServe(addr, mux)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	state := s.driver.GetHardwareData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

func (s *Server) handleSynthesize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var payload struct {
		Code string `json:"code"`
		Mode string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// 1. Save code to a temporary ClojureV file
	tmpFile := filepath.Join(os.TempDir(), "soma_live_synthesis.cljv")
	err := os.WriteFile(tmpFile, []byte(payload.Code), 0644)
	if err != nil {
		http.Error(w, "Failed to write temp file", http.StatusInternalServerError)
		return
	}

	// 2. Run the actual Go transpiler
	// We assume we are running from the project root or SomaServer dir
	// We'll use absolute paths for the tools
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	compilerPath := filepath.Join(projectRoot, "ClojureV", "toolchain", "go", "cmd", "clojurev", "main.go")
	outputVerilog := filepath.Join(projectRoot, "build", "rtl", "sphy_core.v")

	fmt.Printf("[SYNTHESIS] Compiling %s to %s...\n", tmpFile, outputVerilog)
	
	cmd := exec.Command("go", "run", compilerPath, "-target=verilog", "-in="+tmpFile, "-out="+outputVerilog)
	output, err := cmd.CombinedOutput()

	// 3. Trigger physical hardware reset and routing update
	s.driver.SetRoutingMode(payload.Mode)

	// 4. Return the real compiler output
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"output":  string(output),
		"error":   err != nil,
		"message": "Topological Synthesis Complete",
	})
}

func (s *Server) handleReconfigure(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var intent struct {
		Mode string `json:"mode"`
	}
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, "Invalid intent payload", http.StatusBadRequest)
		return
	}
	s.driver.SetRoutingMode(intent.Mode)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "Reconfiguration complete", "mode": intent.Mode})
}

func (s *Server) handleDPR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var intent struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&intent); err != nil {
		http.Error(w, "Invalid intent payload", http.StatusBadRequest)
		return
	}
	s.driver.TriggerDPR(intent.Action)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "DPR sequence initiated", "action": intent.Action})
}
