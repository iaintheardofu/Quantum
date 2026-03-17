import React, { useState, useEffect } from 'react';
import { FileCode, Play, Terminal, ChevronRight, Save, FolderOpen, Loader2 } from 'lucide-react';

const EXAMPLES = {
  'shors_algorithm.cljv': `(ns ClojureV.quantum)

;; Solving Integer Factorization via Period Finding
;; N = p * q. This algorithm provides exponential speedup over classical RSA.

(defn-ai quantum-modular-exponentiation [clk rst_n base power modulus]
  "Core of Shor's: calculates (base^power) mod modulus using topological interference."
  (let [phi (qurq/phi-scale base)]
    (qurq/mod-exp out phi power modulus)))

(defn-fractal shors-factorization [clk rst_n n]
  (let [a (qurq/select-random-coprime n)
        ;; Create superposition of all possible exponents
        x (qurq/superposition 64) 
        f_x (quantum-modular-exponentiation clk rst_n a x n)]
    ;; Quantum Fourier Transform collapses the periodic state
    (let [qft (qurq/qft f_x)]
      (qurq/collapsed-period out qft))))`,

  'grovers_search.cljv': `(ns ClojureV.quantum)

;; Searching an Unstructured Database of 2^64 elements.
;; Classical: O(N) | Quantum: O(sqrt(N))

(defn-ai oracle [clk rst_n signal target]
  "Flips the phase of the target element"
  (if (= signal target)
    (qurq/phi-scale out signal -1.0)
    (qurq/assign out signal)))

(defn-ai diffusion-operator [clk rst_n signal]
  "Inversion about the mean: amplifies the probability of the marked state"
  (let [avg (qurq/read-average signal)]
    (qurq/phi-scale out (- (* 2 avg) signal))))

(defn-fractal grover-search [clk rst_n target]
  (loop [iters 100 signal (qurq/superposition 64)]
    (if (zero? iters)
      (qurq/assign out signal)
      (let [marked (oracle clk rst_n signal target)]
        (recur (dec iters) (diffusion-operator clk rst_n marked))))))`,

  'vqe_molecular_sim.cljv': `(ns ClojureV.quantum)

;; Variational Quantum Eigensolver (VQE)
;; Finding the ground state energy of a Hydrogen molecule (H2)

(defn-ai ansatz-circuit [clk rst_n params]
  "Parameterized quantum circuit representing the molecular wavefunction"
  (let [q0 (qurq/hadamard 0)
        q1 (qurq/rotation-z q0 (first params))]
    (qurq/entangle q0 q1)))

(defn-fractal vqe-hydrogen-sim [clk rst_n initial-params]
  (loop [params initial-params energy 100.0]
    (let [wavefunction (ansatz-circuit clk rst_n params)
          measured-e (qurq/measure-hamiltonian wavefunction :H2-Hamiltonian)]
      (if (< measured-e 0.001) ;; Convergence check
        (qurq/assign out energy)
        (recur (qurq/optimize-params params measured-e) measured-e)))))`
};

export const ClojureVIDE = ({ onClose }: { onClose: () => void }) => {
  const [activeFile, setActiveFile] = useState('shors_algorithm.cljv');
  const [code, setCode] = useState(EXAMPLES['shors_algorithm.cljv' as keyof typeof EXAMPLES]);
  
  const [terminalOutput, setTerminalOutput] = useState<string[]>(['SomaOS ClojureV IDE v4.0.0 initialized.', 'Connected to ALINX 7020 Toolchain. Ready for synthesis...']);
  const [isCompiling, setIsCompiling] = useState(false);
  const [compilationStatus, setCompilationStatus] = useState('');

  const scrollRef = React.useRef<HTMLDivElement>(null);
  const textareaRef = React.useRef<HTMLTextAreaElement>(null);
  const terminalEndRef = React.useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (terminalEndRef.current) {
      terminalEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [terminalOutput, isCompiling]);

  const syncScroll = (e: React.UIEvent<HTMLTextAreaElement>) => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = e.currentTarget.scrollTop;
      scrollRef.current.scrollLeft = e.currentTarget.scrollLeft;
    }
  };

  const pollStatus = async (jobId: string, mode: string) => {
    let lastLogLength = 0;
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`http://localhost:8081/api/synthesize/status?id=${jobId}`);
        if (res.ok) {
          const job = await res.json();
          
          // Append new logs if any
          const fullLogs = (job.output || '') + (job.vivado || '');
          const lines = fullLogs.split('\n').filter((l: string) => l.trim() !== '');
          if (lines.length > lastLogLength) {
            setTerminalOutput(prev => [...prev, ...lines.slice(lastLogLength)]);
            lastLogLength = lines.length;
          }

          if (job.status === 'success') {
            clearInterval(interval);
            setTerminalOutput(prev => [...prev, `[SUCCESS] Physical manifestation complete. Mode: ${mode.toUpperCase()}`]);
            setIsCompiling(false);
            setCompilationStatus('');
          } else if (job.status === 'error') {
            clearInterval(interval);
            setTerminalOutput(prev => [...prev, '[ERROR] Hardware Synthesis Failed. Check Vivado logs above.']);
            setIsCompiling(false);
            setCompilationStatus('');
          }
        }
      } catch (e) {
        clearInterval(interval);
        setTerminalOutput(prev => [...prev, `[CRITICAL ERROR] Lost connection to status endpoint: ${e}`]);
        setIsCompiling(false);
        setCompilationStatus('');
      }
    }, 1000);
  };

  const handleRun = async () => {
    setIsCompiling(true);
    setCompilationStatus('Initializing HPQC Toolchain...');
    setTerminalOutput(prev => [...prev, `> Initiating Live Synthesis for ${activeFile}...`]);
    
    let mode = 'idle';
    if (activeFile.includes('grover')) mode = 'grover';
    else if (activeFile.includes('shor')) mode = 'shor';
    else if (activeFile.includes('vqe')) mode = 'station';

    try {
      const res = await fetch('http://localhost:8081/api/synthesize', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code, mode })
      });

      if (res.ok) {
        const data = await res.json();
        setTerminalOutput(prev => [...prev, `> Synthesis Job Started: ${data.job_id}`]);
        setCompilationStatus('Synthesis & JTAG Flash in Progress...');
        pollStatus(data.job_id, mode);
      } else {
        setTerminalOutput(prev => [...prev, '[ERROR] Synthesis Engine failed to start job.']);
        setIsCompiling(false);
        setCompilationStatus('');
      }
    } catch (e) {
      setTerminalOutput(prev => [...prev, `[CRITICAL ERROR] Connection to toolchain lost: ${e}`]);
      setIsCompiling(false);
      setCompilationStatus('');
    }
  };

  const selectFile = (filename: string) => {
    setActiveFile(filename);
    setCode(EXAMPLES[filename as keyof typeof EXAMPLES]);
    setTerminalOutput(prev => [...prev, `> Opened ${filename}`]);
  };

  const renderHighlightedCode = (content: string, lang: string) => {
    const lines = content.split('\n');
    return lines.map((line, i) => {
      let highlighted = line;
      if (lang === 'cljv') {
        highlighted = highlighted
          .replace(/(\(|\)|\[|\])/g, '<span style="color: #94a3b8">$1</span>')
          .replace(/(defn|defn-ai|defn-fractal|let|loop|if|recur|do|ns)/g, '<span style="color: #60a5fa">$1</span>')
          .replace(/(qurq\/[a-z-]+)/g, '<span style="color: #2dd4bf">$1</span>')
          .replace(/(:[a-z-]+)/g, '<span style="color: #fbbf24">$1</span>')
          .replace(/(;;.*)/g, '<span style="color: #64748b">$1</span>');
      }

      return (
        <div key={i} className="flex" style={{ whiteSpace: 'pre' }}>
          <span className="w-8 text-gray-600 text-right pr-4 select-none" style={{ whiteSpace: 'pre' }}>{i + 1}</span>
          <span dangerouslySetInnerHTML={{ __html: highlighted || '&nbsp;' }} />
        </div>
      );
    });
  };

  return (
    <div className="ide-overlay">
      <div className="ide-container">
        <header className="ide-header">
          <div className="flex items-center gap-2">
            <FileCode className="text-blue-400" size={20} />
            <span className="font-mono text-sm font-bold">ClojureV IDE (HPQC Toolchain)</span>
          </div>
          <div className="flex items-center gap-4">
            <button 
              onClick={handleRun}
              disabled={isCompiling}
              className={`flex items-center gap-1 px-3 py-1 rounded text-xs font-bold transition-all ${isCompiling ? 'bg-gray-700 text-gray-400' : 'bg-green-600 hover:bg-green-500 text-white'}`}
            >
              {isCompiling ? <Loader2 className="animate-spin" size={14} /> : <Play size={14} fill="currentColor" />}
              {isCompiling ? 'SYNTHESIZING & FLASHING...' : 'COMPILE, SYNTHESIZE & FLASH'}
            </button>
            <button onClick={onClose} className="text-gray-400 hover:text-white transition-colors">
              ✕
            </button>
          </div>
        </header>

        <div className="ide-main">
          <div className="ide-sidebar">
            <div className="p-2 border-b border-gray-800 flex items-center gap-2 text-gray-400">
              <FolderOpen size={14} />
              <span className="text-[10px] uppercase font-bold tracking-widest">Quantum Algorithms</span>
            </div>
            <div className="py-2">
              {Object.keys(EXAMPLES).map(filename => (
                <button
                  key={filename}
                  onClick={() => selectFile(filename)}
                  className={`w-full text-left px-4 py-2 text-xs font-mono transition-colors flex items-center gap-2 ${activeFile === filename ? 'bg-blue-900/30 text-blue-400 border-l-2 border-blue-400' : 'text-gray-400 hover:bg-gray-800'}`}
                >
                  <ChevronRight size={10} />
                  {filename.replace('.cljv', '').replace(/_/g, ' ')}
                </button>
              ))}
            </div>

            <div className="p-2 mt-auto border-t border-gray-800 flex items-center gap-2 text-gray-400">
              <span className="text-[10px] uppercase font-bold tracking-widest text-blue-300">ClojureV Quick Ref</span>
            </div>
            <div className="p-3 text-[10px] text-gray-400 font-mono overflow-y-auto" style={{ maxHeight: '250px' }}>
              <div className="mb-2"><span className="text-blue-400 font-bold">defn-fractal</span><br/>Defines the top-level SPHY core manifestation. Maps to the 8 physical macro-cells in continuous asynchronous mode.</div>
              <div className="mb-2"><span className="text-blue-400 font-bold">defn-ai</span><br/>Reusable modular components. Supports Intent String.</div>
              <div className="mb-2"><span className="text-teal-400 font-bold">qurq/phi-scale</span><br/>Scales a signal by Φ (1.618) to induce topological resonance.</div>
              <div className="mb-2"><span className="text-teal-400 font-bold">qurq/mod-exp</span><br/>Modular exponentiation hardware macro for Shor's.</div>
              <div className="mb-2"><span className="text-teal-400 font-bold">qurq/quat-map</span><br/>Maps 24-bit flux to quaternary intermediate states.</div>
            </div>

            <div className="mt-auto p-4 border-t border-gray-800">
              <a 
                href="/tutorial.html" 
                target="_blank" 
                rel="noopener noreferrer"
                className="flex items-center gap-2 text-blue-400 hover:text-blue-300 transition-colors text-[11px] font-bold uppercase tracking-wider group"
              >
                <FileCode size={14} className="group-hover:rotate-12 transition-transform" />
                Full ClojureV Tutorial
              </a>
            </div>
          </div>

          <div className="ide-editor-container relative">
            {isCompiling && (
              <div className="absolute inset-0 z-50 bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center gap-4 transition-all">
                <div className="relative">
                  <Loader2 className="text-blue-500 animate-spin" size={48} />
                  <div className="absolute inset-0 flex items-center justify-center">
                    <div className="w-2 h-2 bg-blue-400 rounded-full animate-ping"></div>
                  </div>
                </div>
                <div className="flex flex-col items-center">
                  <span className="text-blue-400 font-bold text-sm tracking-widest uppercase animate-pulse">{compilationStatus}</span>
                  <span className="text-gray-500 text-[10px] mt-1">Compilation, Synthesis & JTAG Flash in progress...</span>
                </div>
              </div>
            )}
            <div className="ide-editor-tabs">
              <div className="active-tab">
                {activeFile}
              </div>
            </div>
            
            <div className="ide-editor-viewport">
              <div className="ide-editor-content">
                <div ref={scrollRef} className="ide-highlighter">
                  {renderHighlightedCode(code, 'cljv')}
                </div>
                <textarea
                  ref={textareaRef}
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  onScroll={syncScroll}
                  spellCheck={false}
                  className="ide-textarea z-10"
                  style={{ background: 'transparent', color: 'transparent', caretColor: 'white' }}
                />
              </div>
            </div>
          </div>

          <div className="ide-terminal">
            <div className="p-2 bg-black border-b border-gray-800 flex items-center justify-between">
              <div className="flex items-center gap-2 text-gray-400">
                <Terminal size={14} />
                <span className="text-[10px] uppercase font-bold tracking-widest">Vivado & Toolchain Logs</span>
              </div>
              <div className="text-[9px] text-gray-600 font-mono">
                Xilinx Vivado v2025.2
              </div>
            </div>
            <div className="p-4 font-mono text-[11px] flex-1 overflow-y-auto bg-black text-blue-400/90">
              {terminalOutput.map((line, i) => (
                <div key={i} className={`mb-1 ${line.includes('[ERROR]') || line.includes('[CRITICAL]') ? 'text-red-400' : line.includes('[SUCCESS]') || line.includes('[HW]') || line.includes('[JTAG]') ? 'text-green-400' : line.startsWith('  ') ? 'text-gray-400' : ''}`}>
                  {line}
                </div>
              ))}
              {isCompiling && <div className="animate-pulse">_</div>}
              <div ref={terminalEndRef} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};