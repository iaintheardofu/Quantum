import React from 'react';
import { HardwareState } from '../App';
import { LogicBlock } from './LogicBlock';
import { TopologicalFlowWindow } from './TopologicalFlowWindow';
import { Activity, Thermometer, Cpu, Radio, Plus, Minus, BarChart3, Database, Binary } from 'lucide-react';

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

  // Convert histogram map to sorted array for display
  const histogramData = Object.entries(hwState.state_histogram || {})
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8); // Top 8 states

  // Format the 64-bit register as a binary string
  const ghzStateString = hwState.register.toString(2).padStart(hwState.active_cells, '0');

  return (
    <div className="dashboard">
      <div className="telemetry-panel">
        <div className="stat-card">
          <div className="stat-header">
            <Thermometer size={20} className="text-red-400" />
            <h3>XADC Thermal Load</h3>
          </div>
          <div className="stat-value">{hwState.thermal_load.toFixed(2)} °C</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Radio size={20} className="text-green-400" />
            <h3>SPHY Phase Tuner</h3>
          </div>
          <div className="stat-value">Φ {(hwState.phase_field / Math.PI).toFixed(3)} π rad</div>
        </div>

        <div className="stat-card">
          <div className="stat-header">
            <Cpu size={20} className="text-blue-400" />
            <h3>Topology State (d=2^{hwState.active_cells})</h3>
          </div>
          <div className="stat-value" style={{ fontSize: '0.8rem', fontWeight: 'bold', wordBreak: 'break-all', fontFamily: 'monospace' }}>
            |{ghzStateString}⟩ GHZ
          </div>
        </div>
        
        <div className="stat-card dpr-panel">
          <div className="stat-header">
            <Activity size={18} className="text-yellow-400" />
            <h3>DPR Intent Engine</h3>
          </div>
          <div className="dpr-controls" style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
             <button onClick={() => triggerDPR('spawn')} className="dpr-btn" style={{flex: 1, padding: '8px', background: '#1e293b', color: '#00ffcc', border: '1px solid #00ffcc', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px', borderRadius: '4px'}}>
               <Plus size={14} /> Spawn
             </button>
             <button onClick={() => triggerDPR('collapse')} className="dpr-btn" disabled={hwState.active_cells <= 8} style={{flex: 1, padding: '8px', background: '#1e293b', color: '#ff5555', border: '1px solid #ff5555', cursor: 'pointer', display: 'flex', justifyContent: 'center', alignItems: 'center', gap: '5px', opacity: hwState.active_cells <= 8 ? 0.5 : 1, borderRadius: '4px'}}>
               <Minus size={14} /> Collapse
             </button>
          </div>
        </div>
      </div>

      <div style={{ display: 'flex', flexDirection: 'column', gap: '2rem', flex: 1, minHeight: '500px' }}>
        <div className="visualization-panel" style={{ position: 'relative' }}>
          <h2>{hwState.routing_mode === 'station' ? 'Electronic Hyper-Braid: 64-Cell Station' : 'Electronic Entanglement Bus: 8-Cell Macro-Cube'}</h2>
          <LogicBlock state={hwState} />
        </div>
      </div>
    </div>
  );
};