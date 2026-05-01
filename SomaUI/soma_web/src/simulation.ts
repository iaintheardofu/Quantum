import type { HardwareState } from './App';

// Client-side FPGA simulation — mirrors the Go backend's simulateState()
// so the UI runs live even when no backend is available (e.g. Vercel deploy).

let phase = 0;
let windingNumber = 0;
let activeCells = 8;
let routingMode = 'idle';

function randomBytes(n: number): string {
  const hex = '0123456789abcdef';
  let out = '';
  for (let i = 0; i < n * 2; i++) {
    out += hex[Math.floor(Math.random() * 16)];
  }
  return out;
}

function shannonEntropy(manifold: string): number {
  const bytes: number[] = [];
  for (let i = 0; i < manifold.length; i += 2) {
    bytes.push(parseInt(manifold.slice(i, i + 2), 16));
  }
  const counts: Record<number, number> = {};
  for (const b of bytes) counts[b] = (counts[b] || 0) + 1;
  let H = 0;
  for (const c of Object.values(counts)) {
    const p = c / bytes.length;
    H -= p * Math.log2(p);
  }
  return H;
}

export function simulateHardwareState(): HardwareState {
  // Register: random value masked to active cells
  const register = activeCells >= 32
    ? Math.floor(Math.random() * 0xFFFFFFFF)
    : Math.floor(Math.random() * ((1 << activeCells) - 1));

  // Thermal baseline with realistic fluctuation
  const baseTemp = 36.5 + (Math.random() - 0.5) * 0.8;
  const thermalLoad = baseTemp + Math.random() * 0.5;

  // Advance phase (simulated SPHY oscillator)
  phase += 0.15;
  if (phase > 2 * Math.PI) {
    phase = 0;
    windingNumber++;
  }

  // 128-byte manifold vector
  const manifold = randomBytes(128);

  // Shannon entropy from manifold
  const entropy = shannonEntropy(manifold);

  // Decoherence rate
  let decoherence = 0.02 + (thermalLoad - 35.0) * 0.01;
  if (routingMode === 'idle') decoherence += 0.05;
  if (decoherence < 0.001) decoherence = 0.001;
  decoherence += Math.random() * 0.01;

  // Coherence time
  let coherenceTime = 1.0 / (decoherence + 0.0001);
  coherenceTime += (Math.random() - 0.5) * (coherenceTime * 0.02);

  // Histogram from manifold
  const histogram: Record<string, number> = {};
  for (let i = 0; i < manifold.length; i += 2) {
    const bin = manifold.slice(i, i + 2);
    histogram[bin] = (histogram[bin] || 0) + 1;
  }

  return {
    register,
    thermal_load: thermalLoad,
    phase_field: phase,
    active_cells: activeCells,
    routing_mode: routingMode,
    manifold,
    shannon_entropy: entropy,
    coherence_time: coherenceTime,
    state_histogram: histogram,
    hardware_connected: true,
    live_mode: false,
  };
}

export function simDPR(action: 'spawn' | 'collapse') {
  if (action === 'spawn' && activeCells < 64) activeCells += 8;
  if (action === 'collapse' && activeCells > 8) activeCells -= 8;
}

export function simSetRoutingMode(mode: string) {
  routingMode = mode;
}
