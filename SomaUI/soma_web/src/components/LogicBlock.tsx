import React, { useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { OrbitControls, Line, Sphere, Text } from '@react-three/drei';
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
      {/* The physical ring */}
      <mesh>
        <torusGeometry args={[0.5, 0.05, 16, 50]} />
        <meshStandardMaterial color={isActive ? color : '#444c5e'} wireframe={true} emissive={isActive ? color : '#000'} emissiveIntensity={0.5} />
      </mesh>
      {/* The electron state traversing */}
      {isActive && (
        <mesh position={[0.5, 0, 0]}>
          <sphereGeometry args={[0.1, 16, 16]} />
          <meshBasicMaterial color="#FFF" />
        </mesh>
      )}
    </group>
  );
};

// A Single 2x2 Macro-Cell (Geometric Qubit)
const MacroCell = ({ position, state, label, isControl }: { position: [number, number, number], state: HardwareState, label: string, isControl?: boolean }) => {
  const q0Color = state.q0 ? (isControl ? '#00FF00' : '#00FFFF') : '#FF0000';
  const q1Color = state.q1 ? (isControl ? '#0088FF' : '#FF00FF') : '#AA0000';

  return (
    <group position={position}>
      {/* Loop 0 (Polar) */}
      <RingOscillator position={[-0.8, 0, 0]} color={q0Color} isActive={state.q0} speed={2} />
      <Text position={[-0.8, -0.9, 0]} fontSize={0.2} color="white">Q0</Text>
      
      {/* Loop 1 (Azimuthal) */}
      <RingOscillator position={[0.8, 0, 0]} color={q1Color} isActive={state.q1} speed={1.5} />
      <Text position={[0.8, -0.9, 0]} fontSize={0.2} color="white">Q1</Text>

      {/* The 4-NAND Braiding Operator (XOR Bridge) */}
      <Line
        points={[[-0.3, 0, 0], [0.3, 0, 0]]}
        color={(state.q0 || state.q1) ? "#FFFF00" : "#5a667d"}
        lineWidth={3}
      />
      <Sphere position={[0, 0, 0]} args={[0.15, 16, 16]}>
        <meshStandardMaterial color={(state.q0 !== state.q1) ? "#FFFF00" : "#5a667d"} emissive={(state.q0 !== state.q1) ? "#FFFF00" : "#000"} emissiveIntensity={0.8} />
      </Sphere>
      <Text position={[0, 0.5, 0]} fontSize={0.25} color={isControl ? "#00FFcc" : "#AAA"}>{label}</Text>
    </group>
  );
};

const InnerKnot = ({ state }: { state: HardwareState }) => {
  const knotRef = useRef<THREE.Group>(null);

  // The SPHY Phase Field causes the entire macro-cell to pulsate/rotate
  useFrame(() => {
    if (knotRef.current) {
      knotRef.current.rotation.y = state.phase_field * 0.2; // Rotate slower for the whole grid
      // Pulse scale based on thermal load mapping (simulated expansion)
      const scale = 1.0 + (state.thermal_load - 35) * 0.02;
      knotRef.current.scale.set(scale, scale, scale);
    }
  });

  return (
    <group ref={knotRef}>
      {/* C0: Control Qubit (Top) */}
      <MacroCell position={[0, 2.5, 0]} state={state} label="C0 (Control)" isControl={true} />
      
      {/* The Topological Entanglement Bus (Central Splitter) */}
      <Sphere position={[0, 0, 0]} args={[0.3, 32, 32]}>
        <meshStandardMaterial color={state.q0 !== state.q1 ? "#FFFF00" : "#5a667d"} emissive={state.q0 !== state.q1 ? "#FFFF00" : "#000"} emissiveIntensity={1} />
      </Sphere>
      <Text position={[0, -0.5, 0]} fontSize={0.2} color="#FFF">Entanglement Bus</Text>

      {/* Fan-Out Routing Lines */}
      {/* C0 to Bus */}
      <Line points={[[0, 2.0, 0], [0, 0.3, 0]]} color={state.q0 !== state.q1 ? "#FFFF00" : "#5a667d"} lineWidth={2} dashed={true} />
      {/* Bus to C1 */}
      <Line points={[[-0.3, 0, 0], [-2.5, -2.0, 0]]} color={state.q0 !== state.q1 ? "#FFFF00" : "#5a667d"} lineWidth={2} />
      {/* Bus to C2 */}
      <Line points={[[0, -0.3, 0], [0, -2.0, 0]]} color={state.q0 !== state.q1 ? "#FFFF00" : "#5a667d"} lineWidth={2} />
      {/* Bus to C3 */}
      <Line points={[[0.3, 0, 0], [2.5, -2.0, 0]]} color={state.q0 !== state.q1 ? "#FFFF00" : "#5a667d"} lineWidth={2} />

      {/* C1, C2, C3: Target Qubits (Bottom row) */}
      <MacroCell position={[-2.5, -2.5, 0]} state={state} label="C1 (Target)" />
      <MacroCell position={[0, -2.5, 0]} state={state} label="C2 (Target)" />
      <MacroCell position={[2.5, -2.5, 0]} state={state} label="C3 (Target)" />
    </group>
  );
};

export const LogicBlock = ({ state }: { state: HardwareState }) => {
  return (
    <div className="canvas-container">
      <Canvas camera={{ position: [0, 0, 10], fov: 50 }}>
        <color attach="background" args={['#1c2841']} />
        <ambientLight intensity={0.6} />
        <pointLight position={[10, 10, 10]} intensity={2.5} />
        <InnerKnot state={state} />
        <OrbitControls enableZoom={true} />
      </Canvas>
    </div>
  );
};
