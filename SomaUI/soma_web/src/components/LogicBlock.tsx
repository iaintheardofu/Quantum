import React, { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Line, Sphere, Text, Float } from '@react-three/drei';
import * as THREE from 'three';
import { HardwareState } from '../App';

// A single Unbraided NAND Loop
const RingOscillator = ({ position, color, isActive, speed = 1 }: { position: [number, number, number], color: string, isActive: boolean, speed?: number }) => {
  const groupRef = useRef<THREE.Group>(null);
  
  useFrame((state, delta) => {
    if (groupRef.current && isActive) {
      groupRef.current.rotation.x += delta * speed;
      groupRef.current.rotation.y += delta * speed * 0.5;
    }
  });

  return (
    <group ref={groupRef} position={position}>
      <mesh>
        <torusGeometry args={[0.2, 0.02, 16, 50]} />
        <meshStandardMaterial color={isActive ? color : '#444c5e'} wireframe={true} emissive={isActive ? color : '#000'} emissiveIntensity={0.5} />
      </mesh>
    </group>
  );
};

// A Single 2x2 Macro-Cell (Geometric Qubit) - Smaller for Station View
const MacroCell = ({ position, isActive, label, scale = 1 }: { position: [number, number, number], isActive: boolean, label: string, scale?: number }) => {
  const qColor = isActive ? '#00FFcc' : '#FF0000';

  return (
    <group position={position} scale={[scale, scale, scale]}>
      <RingOscillator position={[-0.3, 0, 0]} color={qColor} isActive={isActive} speed={2} />
      <RingOscillator position={[0.3, 0, 0]} color={qColor} isActive={isActive} speed={1.5} />
      <Sphere position={[0, 0, 0]} args={[0.08, 16, 16]}>
        <meshStandardMaterial color={isActive ? "#FFFF00" : "#5a667d"} emissive={isActive ? "#FFFF00" : "#000"} emissiveIntensity={0.8} />
      </Sphere>
      <Text position={[0, 0.4, 0]} fontSize={0.15} color="#AAA">{label}</Text>
    </group>
  );
};

// A Cluster of 8 Qubits (A single 8-Qubit Cube)
const QubitCube = ({ position, register, startIndex, activeCells, label }: { position: [number, number, number], register: number, startIndex: number, activeCells: number, label: string }) => {
  const cubeRef = useRef<THREE.Group>(null);

  // Layout for the 8 qubits in the cube
  const cellPositions = [
    [0, 1.2, 0], [-1, 0, 0.5], [1, 0, 0.5], [-1, 0, -0.5], [1, 0, -0.5], [-1.5, -1, 0], [1.5, -1, 0], [0, -1.5, 0]
  ];

  return (
    <group position={position} ref={cubeRef}>
      {/* Central Bus for this Cube */}
      <Sphere position={[0, 0, 0]} args={[0.15, 16, 16]}>
        <meshStandardMaterial color="#00FFcc" emissive="#00FFcc" emissiveIntensity={0.5} />
      </Sphere>
      <Text position={[0, 1.8, 0]} fontSize={0.25} color="#00FFcc">{label}</Text>

      {cellPositions.map((pos, i) => {
        const qubitIndex = startIndex + i;
        const isActive = (register & (1 << qubitIndex)) !== 0;
        const isVisible = qubitIndex < activeCells;

        if (!isVisible) return null;

        return (
          <group key={i}>
            <Line points={[[0, 0, 0], [pos[0], pos[1], pos[2]]]} color={isActive ? "#FFFF00" : "#333"} lineWidth={1} dashed={true} />
            <MacroCell position={pos as [number, number, number]} isActive={isActive} label={`Q${qubitIndex}`} scale={0.6} />
          </group>
        );
      })}
    </group>
  );
};

const HyperStationKnot = ({ state }: { state: HardwareState }) => {
  const knotRef = useRef<THREE.Group>(null);

  useFrame(() => {
    if (knotRef.current && state.hardware_connected) {
      knotRef.current.rotation.y = state.phase_field * 0.1;
    }
  });

  if (state.routing_mode === 'station') {
    // RENDER 8 CUBES IN A CIRCULAR LAYOUT (The Station Grid)
    const cubeCount = Math.ceil(state.active_cells / 8);
    const radius = 6;

    return (
      <group ref={knotRef}>
        {/* Master Station Hub */}
        <Sphere position={[0, 0, 0]} args={[0.5, 32, 32]}>
          <meshStandardMaterial color="#00FFcc" emissive="#00FFcc" emissiveIntensity={1} wireframe={true} />
        </Sphere>
        <Text position={[0, 1, 0]} fontSize={0.4} color="#00FFcc">MASTER STATION HUB</Text>

        {[...Array(8)].map((_, i) => {
          const angle = (i / 8) * Math.PI * 2;
          const x = Math.cos(angle) * radius;
          const z = Math.sin(angle) * radius;
          const isVisible = i < cubeCount;

          if (!isVisible) return null;

          return (
            <group key={i}>
              <Line points={[[0, 0, 0], [x, 0, z]]} color="#00FFcc" lineWidth={2} />
              <QubitCube 
                position={[x, 0, z]} 
                register={state.register} 
                startIndex={i * 8} 
                activeCells={state.active_cells} 
                label={`CUBE ${i}`} 
              />
            </group>
          );
        })}
      </group>
    );
  }

  // FALLBACK: RENDER THE SINGLE 8-QUBIT CUBE
  return (
    <group ref={knotRef}>
      <QubitCube 
        position={[0, 0, 0]} 
        register={state.register} 
        startIndex={0} 
        activeCells={state.active_cells} 
        label="8-QUBIT MACRO-CUBE" 
      />
    </group>
  );
};

export const LogicBlock = ({ state }: { state: HardwareState }) => {
  return (
    <div className="canvas-container">
      <Canvas camera={{ position: [0, 5, 15], fov: 50 }}>
        <color attach="background" args={['#0a0f1d']} />
        <ambientLight intensity={0.6} />
        <pointLight position={[10, 10, 10]} intensity={2.5} />
        <HyperStationKnot state={state} />
        <OrbitControls enableZoom={true} />
      </Canvas>
    </div>
  );
};
