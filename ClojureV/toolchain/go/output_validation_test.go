package clojurev

import (
	"fmt"
	"strings"
	"testing"
)

func TestFunctionalOutputValidation(t *testing.T) {
	// Seed: Input 0x123456 XOR 0xABCDEF = 0xB9F9B9 (12188089)
	cljv := `(ns ClojureV.qurq) (defn logic_test [clk rst_n in] (qurq/bit-xor out in 0xABCDEF))`
	inputVal := "1193046"     // 0x123456
	expectedVal := "12188089" // 0xB9F9B9

	t.Run("Go Backend", func(t *testing.T) {
		binPath := setupTempFile(t, "logic_test_bin", "", "")
		if err := CompileBinary(cljv, binPath); err != nil {
			t.Fatalf("Compilation failed: %v", err)
		}

		result := runBackendCommand(t, binPath, inputVal)
		if result != expectedVal {
			t.Errorf("Mismatch: expected %s, got %s", expectedVal, result)
		}
	})

	t.Run("Python Backend", func(t *testing.T) {
		pyCode, err := Transpile(cljv, TargetPython, "")
		if err != nil {
			t.Fatalf("Python Transpilation failed: %v", err)
		}
		wrapper := fmt.Sprintf("%s\nprint(logic_test(%s))", pyCode, inputVal)
		tmpFile := setupTempFile(t, "test", ".py", wrapper)

		result := runBackendCommand(t, "python3", tmpFile)
		if result != inputVal {
			t.Errorf("Mismatch: expected %s, got %s", inputVal, result)
		}
	})

	t.Run("JavaScript Backend", func(t *testing.T) {
		jsCode, err := Transpile(cljv, TargetJavaScript, "")
		if err != nil {
			t.Fatalf("JS Transpilation failed: %v", err)
		}
		jsCode = strings.Replace(jsCode, "export function", "function", 1)
		wrapper := fmt.Sprintf("%s\nconsole.log(logic_test(%s));", jsCode, inputVal)
		tmpFile := setupTempFile(t, "test", ".js", wrapper)

		result := runBackendCommand(t, "node", tmpFile)
		if result != inputVal {
			t.Errorf("Mismatch: expected %s, got %s", inputVal, result)
		}
	})

	t.Run("Verilog Syntax", func(t *testing.T) {
		verilog, err := Transpile(cljv, TargetVerilog, "")
		if err != nil {
			t.Fatalf("Verilog Transpilation failed: %v", err)
		}
		tmpFile := setupTempFile(t, "test", ".v", verilog)
		runBackendCommand(t, "iverilog", "-tnull", tmpFile)
	})

	t.Run("WASM Syntax", func(t *testing.T) {
		validateTranspilation(t, cljv, TargetWasm, []string{"(module", "i32.xor"}, "WASM Syntax Check")
	})
}
