package clojurev

import (
	"os"
	"path/filepath"
	"testing"
)

func TestVerilogTranspilation(t *testing.T) {
	tests := []struct {
		name     string
		cljv     string
		expected []string
	}{
		{
			"Basic Passthrough",
			"(ns ClojureV.qurq)\n(defn test_seed [clk rst_n in] (qurq/assign out in))",
			[]string{
				"module test_seed",
				"input wire [23:0] in",
				"output reg [23:0] out",
				"always @(posedge clk)",
				"out = in;",
			},
		},
		{
			"Bit-XOR Transformation",
			"(ns ClojureV.qurq)\n(defn xor_seed [clk rst_n in] (qurq/bit-xor out in 0xABCDEF))",
			[]string{"out = in ^ 24'hABCDEF;"},
		},
		{
			"Bell State Sum-Split",
			"(ns ClojureV.qurq)\n(defn bell_seed [clk rst_n in] (qurq/sum-split out in1 in2))",
			[]string{"{ (in1[23:12] + in1[11:0]), (in1[23:12] + in1[11:0]) };"},
		},
		{
			"Matrix Dot Operator",
			`(ns ClojureV.qurq) (defn matrix_seed [clk rst_n in] (qurq/matrix-dot ". . . ."))`,
			[]string{"out = in & 24'hAA0000;"},
		},
		{
			"AI Intent (defn-ai)",
			`(ns ClojureV.qurq) (defn-ai ai_seed [clk rst_n in] "Manifesting SPHY resonance" (qurq/assign out in))`,
			[]string{"INTENT: Manifesting SPHY resonance"},
		},
		{
			"Fractal-Mimetic Operators",
			`(ns ClojureV.qurq) (defn fractal_seed [clk rst_n in] (qurq/quat-map mid in) (qurq/torsional-pair out mid) (qurq/phi-scale out))`,
			[]string{
				"mid = in;",
				"out = ~mid;",
				"out = (out * 1657) >> 10;",
			},
		},
		{
			"Fractal-Zip Operator",
			`(ns ClojureV.qurq) (defn fzip_seed [clk rst_n in] (qurq/fractal-zip out in))`,
			[]string{"FractalZip: High-density state compression initiated"},
		},
		{
			"Bitwise Operations Suite",
			"(ns ClojureV.qurq)\n(defn logic_suite [clk rst_n in] (qurq/bit-and out in 0x00FF00) (qurq/bit-or out in 0xFF00FF) (qurq/bit-not out in) (qurq/bit-shift-left out in 4) (qurq/bit-shift-right out in 8))",
			[]string{
				"out = in & 24'h00FF00;",
				"out = in | 24'hFF00FF;",
				"out = ~in;",
				"out = in << 4;",
				"out = in >> 8;",
			},
		},
		{
			"Bit Manipulation Suite",
			"(ns ClojureV.qurq)\n(defn manip_suite [clk rst_n in] (qurq/sum-pair out in) (qurq/bit-clear out in 0x0000FF) (qurq/bit-set out in 0xFF0000))",
			[]string{
				"out = in[23:12] + in[11:0];",
				"out = in & ~24'h0000FF;",
				"out = in | 24'hFF0000;",
			},
		},
		{
			"Spacing Regression",
			`(ns ClojureV.qurq)
(defn-ai seed_439 [clk rst_n in]
"Manifesting Linear & Fractal Coherence"
(qurq/quat-map mid in)
(qurq/torsional-pair out mid)
(qurq/phi-scale out)
(qurq/matrix-dot ". . . . . ."))`,
			[]string{
				"mid = in;",
				"out = ~mid;",
				"out = (out * 1657) >> 10;",
				"in & 24'hAAA000",
				"Matrix Dot Mask",
			},
		},
		{
			"SystemVerilog Assertions (SVA)",
			`(ns ClojureV.qurq) (defn-ai sva_seed [clk rst_n in] (qurq/assert-invariant out (qurq/not-equal out 0x0)))`,
			[]string{
				"assert property (@(posedge clk) out != 24'h0);",
			},
		},
		{
			"Photonic Reactor Stream",
			`(ns ClojureV.qurq) (defn photon_seed [clk rst_n in] (qurq/photonic-stream out 0x1))`,
			[]string{
				"out = photon_flux_reg_0x1;",
				"// Photonic Stream: Mapping physical reflection to register 0x1",
			},
		},
		{
			"Topological Interaction",
			`(ns ClojureV.qurq) (defn touch_seed [clk rst_n in] (qurq/swave-interaction out in))`,
			[]string{
				"out = in ^ swave_interaction_mask;",
				"// Topological Interaction: Applying field perturbations",
			},
		},
		{
			"Measurement: ReadQuDot",
			`(ns ClojureV.qurq) (defn measure_seed [clk rst_n in] (qurq/read-qudot out in))`,
			[]string{
				"// ReadQuDot: Topological collapse of 784-qudit field",
				"out = in; // Collapsed state",
			},
		},
		{
			"Measurement: ReadSoundPixel",
			`(ns ClojureV.qurq) (defn aural_seed [clk rst_n in] (qurq/read-sound-pixel out in))`,
			[]string{
				"// ReadSoundPixel: Aural collapse of SPHY waveform (?)",
				"out = in; // Collapsed sound pixel",
			},
		},
		{
			"Measurement: ReadThoughtPixel",
			`(ns ClojureV.qurq) (defn thought_seed [clk rst_n in] (qurq/read-thought-pixel out in))`,
			[]string{
				"// ReadThoughtPixel: Internal observation of the intent manifold (?)",
				"out = in; // Collapsed thought-sound pixel",
			},
		},
		{
			"Transmission: TransmitQuDot",
			`(ns ClojureV.qurq) (defn tx_qudot [clk rst_n in] (qurq/transmit-qu-dot out in))`,
			[]string{
				"// TransmitQuDot: Manifestation of digital intent into 784-qudit field (go)",
				"out = in; // Transmitted state",
			},
		},
		{
			"Transmission: TransmitSoundPixel",
			`(ns ClojureV.qurq) (defn tx_sound [clk rst_n in] (qurq/transmit-sound-pixel out in))`,
			[]string{
				"// TransmitSoundPixel: Manifestation of inner voice into an aural fragment (go)",
				"out = in; // Transmitted sound pixel",
			},
		},
	}

	transpileAndVerifyTable(t, TargetVerilog, tests)
}

func TestWebAndScriptTranspilation(t *testing.T) {
	t.Run("JavaScript Emission", func(t *testing.T) {
		cljv := "(ns ClojureV.qurq)\n(defn js_seed [clk rst_n in] (qurq/assign out in))"
		js, err := Transpile(cljv, TargetJavaScript, "")
		if err != nil {
			t.Fatalf("JS Transpilation failed: %v", err)
		}
		assertContains(t, js, "export function js_seed", "JS Emission")
	})

	t.Run("Python Correctness", func(t *testing.T) {
		cljv := `(ns ClojureV.qurq) (defn py_seed [clk rst_n in] (qurq.math/sin out in))`
		py, err := Transpile(cljv, TargetPython, "")
		if err != nil {
			t.Fatalf("Python Transpilation failed: %v", err)
		}
		assertContains(t, py, "def py_seed(input_flux):", "Python signature")
	})

	t.Run("WASM Correctness", func(t *testing.T) {
		cljv := `(ns ClojureV.qurq) (defn wasm_seed [clk rst_n in] (qurq/bit-xor out in 0xFFFFFF))`
		wat, err := Transpile(cljv, TargetWasm, "")
		if err != nil {
			t.Fatalf("WASM Transpilation failed: %v", err)
		}
		assertContains(t, wat, "(func $wasm_seed", "WASM definition")
	})

	t.Run("Fractal-Zip Multi-Target", func(t *testing.T) {
		cljv := `(ns ClojureV.qurq) (defn fzip_seed [clk rst_n in] (qurq/fractal-zip out in))`
		
		js, err := Transpile(cljv, TargetJavaScript, "")
		if err != nil {
			t.Errorf("JS FZip failed: %v", err)
		} else {
			assertContains(t, js, "export function fzip_seed(input_flux)", "JS FZip")
		}

		py, err := Transpile(cljv, TargetPython, "")
		if err != nil {
			t.Errorf("Python FZip failed: %v", err)
		} else {
			assertContains(t, py, "def fzip_seed(input_flux):", "Python FZip")
		}
	})
}

func TestQuantumTranspilation(t *testing.T) {
	t.Run("Cirq Export Correctness", func(t *testing.T) {
		cljv := `(ns ClojureV.quantum) (defn quantum_seed [clk rst_n in] (q/h q0) (q/cx q0 q1) (q/measure q0))`
		py, err := Transpile(cljv, TargetPython, "")
		if err != nil {
			t.Fatalf("Cirq Transpilation failed: %v", err)
		}

		expectedParts := []string{
			"def quantum_seed(input_flux):",
			"# Translated from ClojureV AST",
		}

		for _, part := range expectedParts {
			assertContains(t, py, part, "Cirq Export")
		}
	})
}

func TestGoTranspilation(t *testing.T) {
	cljv := `(ns ClojureV.qurq) (defn go_seed [clk rst_n in] (qurq/bit-xor out in 0xFFFFFF))`

	t.Run("Go Source Generation", func(t *testing.T) {
		goCode, err := Transpile(cljv, TargetGo, "")
		if err != nil {
			t.Fatalf("Go Transpilation failed: %v", err)
		}
		assertContains(t, goCode, "func go_seed", "Go function")
		assertContains(t, goCode, "package main", "Go package")
	})

	t.Run("Go Bitwise Operators", func(t *testing.T) {
		cljv := `(ns ClojureV.qurq) (defn bit_logic [clk rst_n in] (qurq/bit-and out in 0x00FF00) (qurq/bit-or out in 0xFF00FF) (qurq/bit-not out in) (qurq/bit-shift-left out in 4) (qurq/bit-shift-right out in 8) (qurq/bit-clear out in 0x0000FF) (qurq/bit-set out in 0xFF0000))`
		goCode, err := Transpile(cljv, TargetGo, "")
		if err != nil {
			t.Fatalf("Go Bitwise Transpilation failed: %v", err)
		}
		assertContains(t, goCode, "float64(uint64(in) & uint64(0x00FF00))", "Go Bit-AND")
		assertContains(t, goCode, "float64(uint64(in) | uint64(0xFF00FF))", "Go Bit-OR")
		assertContains(t, goCode, "^uint64(in)", "Go Bit-NOT")
		assertContains(t, goCode, "uint64(in) << 4", "Go Bit-SHL")
		assertContains(t, goCode, "uint64(in) >> 8", "Go Bit-SHR")
		assertContains(t, goCode, "float64(uint64(in) &^ uint64(0x0000FF))", "Go Bit-Clear")
		assertContains(t, goCode, "float64(uint64(in) | uint64(0xFF0000))", "Go Bit-Set")
	})

	// t.Run("LaTeX to cljv Transpilation", func(t *testing.T) {
	// 	latex := `f(x) = \frac{x^2 + 1}{\sqrt{x}}`
	// 	cljv, err := TranspileLaTeX(latex)
	// 	if err != nil {
	// 		t.Fatalf("LaTeX transpilation failed: %v", err)
	// 	}
	// 	if !strings.Contains(cljv, "defn f [x]") || !strings.Contains(cljv, "math/sqrt x") {
	// 		t.Errorf("Mismatch in basic LaTeX: got %s", cljv)
	// 	}

	// 	latexTrig := `g(theta) = \sin(theta) + \cos(theta)`
	// 	cljvTrig, err := TranspileLaTeX(latexTrig)
	// 	if err != nil {
	// 		t.Fatalf("LaTeX trig transpilation failed: %v", err)
	// 	}
	// 	if !strings.Contains(cljvTrig, "qurq.math/sin theta") || !strings.Contains(cljvTrig, "qurq.math/cos theta") {
	// 		t.Errorf("Mismatch in trig LaTeX: got %s", cljvTrig)
	// 	}

	// 	latexRound := `h(x) = \text{round}(|x|)`
	// 	cljvRound, err := TranspileLaTeX(latexRound)
	// 	if err != nil {
	// 		t.Fatalf("LaTeX round transpilation failed: %v", err)
	// 	}
	// 	if !strings.Contains(cljvRound, "math/round") || !strings.Contains(cljvRound, "math/abs") {
	// 		t.Errorf("Mismatch in round LaTeX: got %s", cljvRound)
	// 	}
	// })

	t.Run("Go Advanced Operators", func(t *testing.T) {
		cljv := `(ns ClojureV.qurq) (defn advanced_logic [clk rst_n in] (qurq/sum-pair out in))`
		goCode, err := Transpile(cljv, TargetGo, "")
		if err != nil {
			t.Fatalf("Go Advanced Transpilation failed: %v", err)
		}
		assertContains(t, goCode, "float64(uint64(int64(in) >> 12) + uint64(int64(in) & 0xFFF))", "Go Sum-Pair")
	})

	t.Run("Go Binary Compilation", func(t *testing.T) {
		outputPath := filepath.Join(t.TempDir(), "go_seed_bin")
		err := CompileBinary(cljv, outputPath)
		if err != nil {
			t.Fatalf("Go Binary Compilation failed: %v", err)
		}

		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Errorf("Binary file was not created at %s", outputPath)
		}
	})
}

func TestLinter(t *testing.T) {
	tests := []struct {
		name     string
		cljv     string
		errCount int
	}{
		{"Valid Code", `(ns ClojureV.qurq) (defn valid [clk rst_n in] (qurq/assign out in))`, 0},
		{"Invalid Namespace", `(ns Wrong.Namespace) (defn invalid [clk rst_n in] (qurq/assign out in))`, 1},
		{"Unbalanced Parentheses", `(ns ClojureV.qurq) (defn unbal [clk rst_n in] (qurq/assign out in`, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := Lint(tt.cljv)
			if len(errs) != tt.errCount {
				t.Errorf("%s: expected %d lint errors, got %d: %v", tt.name, tt.errCount, len(errs), errs)
			}
		})
	}
}
