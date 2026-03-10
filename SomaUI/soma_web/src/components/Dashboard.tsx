import React from 'react';
import { HardwareState } from '../App';
import { LogicBlock } from './LogicBlock';
import { TopologicalFlowWindow } from './TopologicalFlowWindow';
import { Activity, Thermometer, Cpu, Radio, Plus, Minus } from 'lucide-react';

export const Dashboard = ({ hwState }: { hwState: HardwareState }) => {

  const triggerDPR = async (action: 'spawn' | 'collapse') => {
    try {
      await fetch('http://localhost:8081/api/dpr', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ action })
      });
    } catch (e) {
      console.error("DPR Injection Failed", e);
    }
  };

  const dimension = Math.pow(2, hwState.active_cells);
  
  // Format the 64-bit register as a binary string, padded based on active cells
  const ghzStateString = hwState.register.toString(2).padStart(hwState.active_cells, '0');

  return (
    <div className="dashboard">
      <div className="telemetry-panel">
        <div className="stat-card">
          <div className="stat-header">
            <Thermometer size={20} color="#ff5555" />
            <h3>XADC Thermal Load</h3>
          </div>
          <div className="stat-value">{hwState.thermal_load.toFixed(2)} °C</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Radio size={20} color="#55ff55" />
            <h3>SPHY Phase Tuner</h3>
          </div>
          <div className="stat-value">Φ {hwState.phase_field.toFixed(3)} rad</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Cpu size={20} color="#5555ff" />
            <h3>Topology State (d=2^{hwState.active_cells})</h3>
          </div>
          <div className="stat-value" style={{ fontSize: '0.9rem', fontWeight: 'bold', wordBreak: 'break-all', fontFamily: 'monospace' }}>
            |{ghzStateString}⟩ GHZ
          </div>
        </div>
        
        <div className="stat-card dpr-panel">
          <div className="stat-header">
            <h3>DPR Intent Engine</h3>
          </div>
          <div className="dpr-controls" style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
             <button onClick={() => triggerDPR('spawn')} className="dpr-btn" style={{flex: 1, padding: '10px', background: '#2a3b5c', color: '#00ffcc', border: '1px solid #00ffcc', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px'}}>
               <Plus size={16} /> Spawn Cell
             </button>
             <button onClick={() => triggerDPR('collapse')} className="dpr-btn" disabled={hwState.active_cells <= 8} style={{flex: 1, padding: '10px', background: '#2a3b5c', color: '#ff5555', border: '1px solid #ff5555', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px', opacity: hwState.active_cells <= 8 ? 0.5 : 1}}>
               <Minus size={16} /> Collapse
             </button>
          </div>
        </div>

      </div>

      <div className="visualization-panel" style={{ position: 'relative' }}>
        <h2>{hwState.routing_mode === 'station' ? 'Fractal Hypercube: 64-Qubit Station Hub' : 'Topological Entanglement Bus: 8-Qubit Macro-Cube'}</h2>
        <LogicBlock state={hwState} />
        <TopologicalFlowWindow state={hwState} />
      </div>
    </div>
  );
};
