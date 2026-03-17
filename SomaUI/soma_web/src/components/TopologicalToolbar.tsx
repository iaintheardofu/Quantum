import React from 'react';
import { BarChart3, Database, Binary } from 'lucide-react';
import { HardwareState } from '../App';

export const TopologicalToolbar = ({ hwState }: { hwState: HardwareState }) => {
  // Convert histogram map to sorted array for display
  const histogramData = Object.entries(hwState.state_histogram || {})
    .sort((a, b) => b[1] - a[1])
    .slice(0, 5); // Show top 5 states in toolbar to save space

  return (
    <div className="topological-toolbar">
      <div className="toolbar-section">
        <div className="text-[9px] text-gray-500 font-bold uppercase mb-1 flex items-center gap-1">
          <BarChart3 size={10} /> Shannon Entropy (H)
        </div>
        <div className="text-sm font-mono text-blue-300">{hwState.shannon_entropy?.toFixed(4)}</div>
      </div>
      
      <div className="toolbar-section">
        <div className="text-[9px] text-gray-500 font-bold uppercase mb-1">Coherence (T2)</div>
        <div className="text-sm font-mono text-green-300">{hwState.coherence_time?.toFixed(2)} ms</div>
      </div>

      <div className="toolbar-section flex-1 border-x border-gray-800 px-4 mx-4">
        <div className="text-[9px] text-gray-500 font-bold uppercase mb-1 flex items-center gap-1">
          <Database size={10}/> Topological State Distribution (Top 5)
        </div>
        <div className="flex gap-4">
          {histogramData.length > 0 ? histogramData.map(([state, count]) => (
            <div key={state} className="flex items-center gap-1 flex-1">
              <span className="font-mono text-[9px] text-gray-400">{state}</span>
              <div className="flex-1 h-1 bg-gray-900 rounded-full overflow-hidden min-w-[20px]">
                <div 
                  className="h-full bg-blue-500/60 shadow-[0_0_4px_rgba(59,130,246,0.5)]" 
                  style={{ width: `${Math.min(100, (count / 128) * 100)}%`, transition: 'width 0.3s ease' }}
                />
              </div>
            </div>
          )) : (
            <div className="text-xs italic text-gray-600">Awaiting telemetry...</div>
          )}
        </div>
      </div>

      <div className="toolbar-section max-w-[250px]">
        <div className="text-[9px] text-gray-500 font-bold uppercase mb-1 flex items-center gap-1">
          <Binary size={10}/> Raw Manifold Vector
        </div>
        <div className="font-mono text-[8px] text-blue-400/70 truncate border border-blue-900/20 bg-black/50 p-1 rounded">
          {hwState.manifold ? `${hwState.manifold.substring(0, 64)}...` : '0x00000000000000...'}
        </div>
      </div>
    </div>
  );
};