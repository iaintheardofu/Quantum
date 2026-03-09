import React from 'react';
import { HardwareState } from '../App';
import { LogicBlock } from './LogicBlock';
import { Activity, Thermometer, Cpu, Radio } from 'lucide-react';

export const Dashboard = ({ hwState }: { hwState: HardwareState }) => {
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
            <h3>Topology State (d=16)</h3>
          </div>
          <div className="stat-value" style={{ fontSize: '1.2rem', fontWeight: 'bold' }}>
            |{hwState.q0 !== hwState.q1 ? '1111' : '0000'}⟩ GHZ
          </div>
        </div>
      </div>

      <div className="visualization-panel">
        <h2>Topological Entanglement Bus (d=16 GHZ Macro-State)</h2>
        <LogicBlock state={hwState} />
      </div>
    </div>
  );
};
