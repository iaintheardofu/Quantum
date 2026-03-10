import React, { useEffect, useState } from 'react'
import { Dashboard } from './components/Dashboard'
import { ClojureVIDE } from './components/ClojureVIDE'
import { TopologicalFlowWindow } from './components/TopologicalFlowWindow'

export type HardwareState = {
  // 64-bit register representing 8 Entanglement Stations
  register: number;
  thermal_load: number;
  phase_field: number;
  active_cells: number;
  routing_mode: string;
};

function App() {
  const [hwState, setHwState] = useState<HardwareState>({
    register: 0,
    thermal_load: 35.0,
    phase_field: 0.0,
    active_cells: 8,
    routing_mode: 'idle'
  });

  const [isIdeOpen, setIsIdeOpen] = useState(false);

  useEffect(() => {
    const fetchState = async () => {
      try {
        const res = await fetch('http://localhost:8081/api/state');
        if (res.ok) {
          const data = await res.json();
          setHwState(data);
        }
      } catch (err) {
        // Silently fail if server is down, keep showing old state
      }
    };

    // Poll server at 10Hz to match hardware simulator
    const interval = setInterval(fetchState, 100);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="app-container">
      <header className="header">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h1>SomaOS</h1>
            <p>Geometric Virtualization Visualizer | {hwState.routing_mode === 'station' ? '64-Qubit Station Hub' : '8-Qubit Macro-Cube'}</p>
          </div>
          <button 
            onClick={() => setIsIdeOpen(true)}
            style={{ 
              backgroundColor: '#2563eb', 
              color: 'white', 
              padding: '0.5rem 1rem', 
              borderRadius: '0.25rem', 
              fontSize: '0.875rem', 
              fontWeight: 'bold',
              border: 'none',
              cursor: 'pointer'
            }}
          >
            OPEN CLOJUREV IDE
          </button>
        </div>
      </header>
      <main className="main-content">
        <Dashboard hwState={hwState} />
      </main>

      {/* Floating Topological Flow Window (Persists over everything) */}
      <TopologicalFlowWindow state={hwState} />

      {isIdeOpen && <ClojureVIDE onClose={() => setIsIdeOpen(false)} />}
    </div>
  )
}

export default App
