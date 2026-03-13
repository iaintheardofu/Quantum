package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"soma_server/hardware"
	"testing"
)

func TestExhaustiveAPIEndpoints(t *testing.T) {
	driver := hardware.NewFPGADriver()
	server := NewServer(driver)

	tests := []struct {
		name       string
		method     string
		endpoint   string
		payload    interface{}
		wantStatus int
	}{
		{"State GET", "GET", "/api/state", nil, http.StatusOK},
		{"State OPTIONS", "OPTIONS", "/api/state", nil, http.StatusOK},
		{"Reconfigure POST Valid", "POST", "/api/reconfigure", map[string]string{"mode": "grover"}, http.StatusOK},
		{"Reconfigure OPTIONS", "OPTIONS", "/api/reconfigure", nil, http.StatusOK},
		{"Reconfigure GET Invalid Method", "GET", "/api/reconfigure", nil, http.StatusMethodNotAllowed},
		{"Reconfigure POST Invalid JSON", "POST", "/api/reconfigure", "bad json", http.StatusBadRequest},
		{"DPR POST Valid", "POST", "/api/dpr", map[string]string{"action": "spawn"}, http.StatusOK},
		{"DPR OPTIONS", "OPTIONS", "/api/dpr", nil, http.StatusOK},
		{"DPR GET Invalid Method", "GET", "/api/dpr", nil, http.StatusMethodNotAllowed},
		{"DPR POST Invalid JSON", "POST", "/api/dpr", "bad json", http.StatusBadRequest},
		{"Synthesize OPTIONS", "OPTIONS", "/api/synthesize", nil, http.StatusOK},
		{"Synthesize POST Invalid JSON", "POST", "/api/synthesize", "bad json", http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.payload != nil {
				if str, ok := tc.payload.(string); ok {
					req = httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBufferString(str))
				} else {
					jsonData, _ := json.Marshal(tc.payload)
					req = httptest.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(jsonData))
				}
			} else {
				req = httptest.NewRequest(tc.method, tc.endpoint, nil)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			switch tc.endpoint {
			case "/api/state":
				server.handleState(rr, req)
			case "/api/reconfigure":
				server.handleReconfigure(rr, req)
			case "/api/dpr":
				server.handleDPR(rr, req)
			case "/api/synthesize":
				server.handleSynthesize(rr, req)
			}

			if status := rr.Code; status != tc.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tc.wantStatus)
			}
		})
	}
}
