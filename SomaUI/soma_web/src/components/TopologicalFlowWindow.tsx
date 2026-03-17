import React, { useRef, useMemo, useState, useEffect } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Sphere, Line, Text, Float } from '@react-three/drei';
import * as THREE from 'three';
import { HardwareState } from '../App';

const ElectronicFlow = ({ mode }: { mode: string }) => {
  const pointsRef = useRef<THREE.Points>(null);
  const particleCount = 1500;

  // Create different flow paths based on the routing mode
  const particles = useMemo(() => {
    const positions = new Float32Array(particleCount * 3);
    const velocities = new Float32Array(particleCount);
    for (let i = 0; i < particleCount; i++) {
      positions[i * 3] = (Math.random() - 0.5) * 4;
      positions[i * 3 + 1] = (Math.random() - 0.5) * 4;
      positions[i * 3 + 2] = (Math.random() - 0.5) * 4;
      velocities[i] = 0.05 + Math.random() * 0.1; // Electricity flow
    }
    return { positions, velocities };
  }, [mode]);

  useFrame((state) => {
    const time = state.clock.getElapsedTime();
    const positions = pointsRef.current!.geometry.attributes.position.array as Float32Array;

    for (let i = 0; i < particleCount; i++) {
      const idx = i * 3;
      
      if (mode === 'grover') {
        // Spiral Braiding Flow
        const angle = time * 3 + i * 0.1;
        const radius = 1.5 + Math.sin(time + i) * 0.5;
        positions[idx] = Math.cos(angle) * radius;
        positions[idx + 1] = Math.sin(angle) * radius;
        positions[idx + 2] = Math.sin(angle * 0.5) * 2;
      } else if (mode === 'shor') {
        // Harmonic Frequency Flow (QFT)
        const freq = 0.5 + (i % 8) * 0.2;
        positions[idx] = Math.sin(time * freq + i) * 2;
        positions[idx + 1] = Math.cos(time * freq * 0.5 + i) * 2;
        positions[idx + 2] = (i / particleCount - 0.5) * 4;
      } else if (mode === 'station') {
        // Master Station Hub: Fractal Star Burst
        const t = time * 2;
        const radius = (i % 3) + 1;
        const angle1 = (i / particleCount) * Math.PI * 2 + t;
        const angle2 = ((i * 3) / particleCount) * Math.PI * 2 - t;
        positions[idx] = Math.cos(angle1) * Math.sin(angle2) * radius;
        positions[idx + 1] = Math.sin(angle1) * Math.sin(angle2) * radius;
        positions[idx + 2] = Math.cos(angle2) * radius;
      } else {
        // Idle Random Flux (Thermal)
        positions[idx] += (Math.random() - 0.5) * 0.04;
        positions[idx + 1] += (Math.random() - 0.5) * 0.04;
        positions[idx + 2] += (Math.random() - 0.5) * 0.04;
        
        if (Math.abs(positions[idx]) > 3) positions[idx] = 0;
      }
    }
    pointsRef.current!.geometry.attributes.position.needsUpdate = true;
  });

  return (
    <points ref={pointsRef}>
      <bufferGeometry>
        <bufferAttribute
          attach="attributes-position"
          count={particleCount}
          array={particles.positions}
          itemSize={3}
        />
      </bufferGeometry>
      <pointsMaterial
        size={0.06}
        color={mode === 'idle' ? '#4a5568' : '#3b82f6'} // Electric Blue / Cobalt High Voltage
        transparent
        opacity={0.8}
        blending={THREE.AdditiveBlending}
      />
    </points>
  );
};

export const TopologicalFlowWindow = ({ state }: { state: HardwareState }) => {
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const dragStart = useRef({ x: 0, y: 0 });

  const handleMouseDown = (e: React.MouseEvent) => {
    setIsDragging(true);
    dragStart.current = { x: e.clientX - position.x, y: e.clientY - position.y };
  };

  useEffect(() => {
    if (!isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      setPosition({
        x: e.clientX - dragStart.current.x,
        y: e.clientY - dragStart.current.y
      });
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    window.addEventListener('mousemove', handleMouseMove);
    window.addEventListener('mouseup', handleMouseUp);

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      window.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isDragging]);

  return (
    <div 
      className="flow-window"
      style={{ transform: `translate(${position.x}px, ${position.y}px)` }}
    >
      <div 
        className="flow-window-header"
        onMouseDown={handleMouseDown}
        style={{ cursor: isDragging ? 'grabbing' : 'grab' }}
      >
        <span className="text-[10px] font-bold tracking-widest text-blue-400 uppercase pointer-events-none select-none">
          Silicon Routing: {state.routing_mode.toUpperCase()}
        </span>
      </div>
      <div className="flow-canvas-container">
        <Canvas camera={{ position: [0, 0, 8], fov: 40 }}>
          <color attach="background" args={['#0a0f1d']} />
          <ambientLight intensity={0.5} />
          
          <Float speed={2} rotationIntensity={0.5} floatIntensity={0.5}>
            {/* The schematic logic block loop */}
            <mesh rotation={[Math.PI / 2, 0, 0]}>
              <torusGeometry args={[2, 0.02, 16, 100]} />
              <meshBasicMaterial color="#1e293b" transparent opacity={0.3} />
            </mesh>
            
            <ElectronicFlow mode={state.routing_mode} />
            
            {/* Legend Text */}
            <Text position={[0, -2.5, 0]} fontSize={0.2} color="#4b5563">
              REAL-TIME ELECTRONIC TOPOLOGICAL INTERFERENCE
            </Text>
          </Float>
        </Canvas>
      </div>
    </div>
  );
};
