import React, { useEffect, useState } from 'react'
import { Dashboard } from './components/Dashboard'
import { ClojureVIDE } from './components/ClojureVIDE'
import { Wifi, WifiOff, AlertTriangle, Monitor } from 'lucide-react'
import { TopologicalFlowWindow } from './components/TopologicalFlowWindow'
import { TopologicalToolbar } from './components/TopologicalToolbar'

export type HardwareState = {
  // 64-bit register representing 8 Entanglement Stations
  register: number;
  thermal_load: number;
  phase_field: number;
  active_cells: number;
  routing_mode: string;
  manifold?: string;
  shannon_entropy: number;
  coherence_time: number;
  state_histogram: Record<string, number>;
  hardware_connected: bool;
};

function App() {
  const [hwState, setHwState] = useState<HardwareState>({
    register: 0,
    thermal_load: 35.0,
    phase_field: 0.0,
    active_cells: 8,
    routing_mode: 'idle',
    shannon_entropy: 0,
    coherence_time: 0,
    state_histogram: {},
    hardware_connected: false
  });

  const [isIdeOpen, setIsIdeOpen] = useState(false);

  useEffect(() => {
    const fetchState = async () => {
      try {
        const res = await fetch('http://localhost:8081/api/state');
        if (res.ok) {
          const data = await res.json();
          setHwState(data);
        } else {
          setHwState(prev => ({ ...prev, hardware_connected: false, register: 0, phase_field: 0, thermal_load: 0, shannon_entropy: 0, coherence_time: 0, state_histogram: {} }));
        }
      } catch (err) {
        // Server is down, reflect disconnected state
        setHwState(prev => ({ ...prev, hardware_connected: false, register: 0, phase_field: 0, thermal_load: 0, shannon_entropy: 0, coherence_time: 0, state_histogram: {} }));
      }
    };

    // Poll server at 10Hz to match hardware simulator
    const interval = setInterval(fetchState, 100);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="app-container">
      <div className="top-toolbar">
        <div className="top-toolbar-left">
          {hwState.hardware_connected ? (
            <div className="connection-success">
              <div className="flex items-center gap-3">
                <Wifi className="text-green-500" size={24} />
                <div className="flex flex-col">
                  <span className="text-green-400 font-bold text-[10px] uppercase tracking-widest">LIVE ON FPGA</span>
                  <span className="text-gray-400 text-[8px] uppercase font-bold">Connected to soma_agent</span>
                </div>
              </div>
            </div>
          ) : (
            <div className="connection-warning">
              <div className="flex items-center gap-3">
                <div className="relative">
                  <WifiOff className="text-red-500 animate-pulse" size={24} />
                  <AlertTriangle className="absolute -top-2 -right-2 text-yellow-500" size={14} />
                </div>
                <div className="flex flex-col">
                  <span className="text-red-400 font-bold text-[10px] uppercase tracking-widest">ALINX LINK FAILURE</span>
                  <span className="text-gray-500 text-[8px] uppercase font-bold">Simulating SPHY Manifold...</span>
                </div>
              </div>
            </div>
          )}
        </div>
        
        <div className="top-toolbar-right">
          <button 
            onClick={() => setIsIdeOpen(true)}
            style={{ 
              backgroundColor: '#2563eb', 
              color: 'white', 
              padding: '0.4rem 0.8rem', 
              borderRadius: '0.25rem', 
              fontSize: '0.75rem', 
              fontWeight: 'bold',
              border: 'none',
              cursor: 'pointer'
            }}
          >
            OPEN CLOJUREV IDE
          </button>
        </div>
      </div>

      <header className="app-header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>SomaOS</h1>
            <p>Geometric Virtualization Visualizer | {hwState.routing_mode === 'station' ? '64-Qubit Station Hub' : '8-Qubit Macro-Cube'}</p>
          </div>
        </div>
      </header>
      <main className="main-content">
        <Dashboard hwState={hwState} />
      </main>

      <TopologicalFlowWindow state={hwState} />
      <TopologicalToolbar hwState={hwState} />

      {isIdeOpen && <ClojureVIDE onClose={() => setIsIdeOpen(false)} />}
    </div>
  )
}

export default App
