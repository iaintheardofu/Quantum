package compiler

import (
	"clojurev/parser"
	"fmt"
	"strings"
)

func EmitVerilog(ast *parser.Program) (string, error) {
	var sb strings.Builder

	for _, node := range ast.Body {
		if defn, ok := node.(*parser.Defn); ok {
			code, err := emitVerilogDefn(defn)
			if err != nil {
				return "", err
			}
			sb.WriteString(code)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

func emitVerilogDefn(d *parser.Defn) (string, error) {
	var v strings.Builder

	v.WriteString(fmt.Sprintf("// SomaOS v3.0.0: Verilog Manifestation of %s\n", d.Name))
	if d.IsAI && d.Intent != "" {
		v.WriteString(fmt.Sprintf("// AI INTENT: %s\n", d.Intent))
	}

	modName := strings.ReplaceAll(d.Name, "-", "_")
	if d.IsFractal {
		modName = "sphy_core"
	}
	v.WriteString(fmt.Sprintf("module %s (\n", modName))
	
	var thirdParam string
	for i, p := range d.Params {
		pName := strings.ReplaceAll(p, "-", "_")
		if i == 2 && d.IsFractal {
			v.WriteString(fmt.Sprintf("    input wire [23:0] in_flux,\n"))
			thirdParam = pName
		} else if pName == "clk" || pName == "rst_n" {
			v.WriteString(fmt.Sprintf("    input wire %s,\n", pName))
		} else {
			v.WriteString(fmt.Sprintf("    input wire [23:0] %s,\n", pName))
		}
	}
	v.WriteString("    output reg [23:0] out\n);\n")

	if thirdParam != "" {
		v.WriteString(fmt.Sprintf("    wire [23:0] %s = in_flux;\n", thirdParam))
	}

	// Determine if we need local variables (mid, etc.)
	localVars := make(map[string]bool)
	needsMid := false
	hasSVA := false
	for _, node := range d.Body {
		walkAST(node, func(n parser.Node) {
			if call, ok := n.(*parser.Call); ok {
				if call.Callee == "quat-map" || call.Callee == "qurq/quat-map" {
					needsMid = true
				}
				if call.Callee == "assert-invariant" || call.Callee == "qurq/assert-invariant" {
					hasSVA = true
				}
				if call.Callee == "let" {
					if len(call.Args) > 0 {
						if bindings, ok := call.Args[0].(*parser.List); ok {
							for i := 0; i < len(bindings.Elements); i += 2 {
								if ident, isIdent := bindings.Elements[i].(*parser.Identifier); isIdent {
									localVars[strings.ReplaceAll(ident.Name, "-", "_")] = true
								}
							}
						}
					}
				}
			}
		})
	}

	if needsMid {
		v.WriteString("    reg [23:0] mid;\n")
	}
	for vName := range localVars {
		v.WriteString(fmt.Sprintf("    reg [23:0] %s;\n", vName))
	}

	v.WriteString("    always @(posedge clk) begin\n")
	v.WriteString("        if (!rst_n) begin\n            out = 24'h0;\n")
	if needsMid {
		v.WriteString("            mid = 24'h0;\n")
	}
	v.WriteString("        end else begin\n")

	hasAssignments := false
	for _, stmt := range d.Body {
		assignCode, err := emitVerilogStatement(stmt)
		if err != nil {
			return "", err
		}
		if assignCode != "" {
			v.WriteString(assignCode)
			hasAssignments = true
		}
	}

	if !hasAssignments {
		if len(d.Params) > 2 {
			v.WriteString(fmt.Sprintf("            out = %s;\n", strings.ReplaceAll(d.Params[len(d.Params)-1], "-", "_")))
		} else {
			v.WriteString("            out = 24'h0;\n")
		}
	}

	v.WriteString("        end\n    end\n")

	if hasSVA {
		v.WriteString("    // SVA: Asserting Topological Invariance\n")
		v.WriteString("    assert property (@(posedge clk) out != 24'h0);\n")
	}

	v.WriteString("endmodule\n")
	return v.String(), nil
}

func walkAST(node parser.Node, fn func(parser.Node)) {
	if node == nil {
		return
	}
	fn(node)
	switch n := node.(type) {
	case *parser.Defn:
		for _, b := range n.Body {
			walkAST(b, fn)
		}
	case *parser.Call:
		for _, a := range n.Args {
			walkAST(a, fn)
		}
	case *parser.List:
		for _, e := range n.Elements {
			walkAST(e, fn)
		}
	}
}

func emitVerilogStatement(node parser.Node) (string, error) {
	call, ok := node.(*parser.Call)
	if !ok {
		return "", nil // ignore non-calls at top level for now
	}

	op := strings.TrimPrefix(call.Callee, "qurq/")
	var v strings.Builder

	formatArg := func(n parser.Node) string {
		switch arg := n.(type) {
		case *parser.Identifier:
			return strings.ReplaceAll(arg.Name, "-", "_")
		case *parser.Number:
			if strings.HasPrefix(arg.Value, "0x") {
				return "24'h" + strings.TrimPrefix(arg.Value, "0x")
			}
			return arg.Value
		default:
			return "24'hAAAAAA" // Fallback placeholder
		}
	}

	if op == "let" {
	        // Handle let bindings topologically
	        if len(call.Args) > 0 {
	                if bindings, ok := call.Args[0].(*parser.List); ok {
	                        for i := 0; i < len(bindings.Elements); i += 2 {
	                                if i+1 < len(bindings.Elements) {
	                                        if ident, isIdent := bindings.Elements[i].(*parser.Identifier); isIdent {
	                                                if valCall, isCall := bindings.Elements[i+1].(*parser.Call); isCall {
	                                                        valOp := strings.TrimPrefix(valCall.Callee, "qurq/")
	                                                        if valOp == "measure-intent-pressure" {
	                                                                v.WriteString(fmt.Sprintf("            // Sensing intent pressure into %s\n", strings.ReplaceAll(ident.Name, "-", "_")))
	                                                        } else if valOp == "read-topological-dimension" {
	                                                                v.WriteString(fmt.Sprintf("            // Reading topological dimension into %s\n", strings.ReplaceAll(ident.Name, "-", "_")))
	                                                        } else if valOp == "spawn-station-bus" {
	                                                                v.WriteString(fmt.Sprintf("            wire [63:0] %s;\n", strings.ReplaceAll(ident.Name, "-", "_")))
	                                                                if len(valCall.Args) > 0 {
	                                                                        v.WriteString(fmt.Sprintf("            assign %s = {64{%s[0]}};\n", strings.ReplaceAll(ident.Name, "-", "_"), formatArg(valCall.Args[0])))
	                                                                }
	                                                        } else {
	                                                                v.WriteString(fmt.Sprintf("            // Let binding: %s\n", strings.ReplaceAll(ident.Name, "-", "_")))
	                                                        }	                                                } else if valNum, isNum := bindings.Elements[i+1].(*parser.Number); isNum {
	                                                        v.WriteString(fmt.Sprintf("            %s = %s;\n", strings.ReplaceAll(ident.Name, "-", "_"), formatArg(valNum)))
	                                                }
	                                        }
	                                }
	                        }
	                }
	        }
	        if len(call.Args) > 1 {
	                for _, n := range call.Args[1:] {
	                        code, err := emitVerilogStatement(n)
	                        if err == nil {
	                                v.WriteString(code)
	                        }
	                }
	        }
	        return v.String(), nil
	}
	switch op {
	case "if":
	        v.WriteString("            // IF Intent Evaluated\n")
		condStr := "1"
		if len(call.Args) > 0 {
			if condCall, isCall := call.Args[0].(*parser.Call); isCall {
				opCond := strings.TrimPrefix(condCall.Callee, "qurq/")
				if opCond == "=" || opCond == "equal" {
					condStr = fmt.Sprintf("%s == %s", formatArg(condCall.Args[0]), formatArg(condCall.Args[1]))
				}
			}
		}
		v.WriteString(fmt.Sprintf("            if (%s) begin\n", condStr))
	        if len(call.Args) > 1 {
	                code, _ := emitVerilogStatement(call.Args[1])
	                v.WriteString(code)
	        }
		v.WriteString("            end else begin\n")
	        if len(call.Args) > 2 {
	                code, _ := emitVerilogStatement(call.Args[2])
	                v.WriteString(code)
	        }
		v.WriteString("            end\n")
	case "when":
	        v.WriteString("            // WHEN Intent Evaluated\n")
	        if len(call.Args) > 0 {
	                condCode, _ := emitVerilogStatement(call.Args[0])
	                v.WriteString(condCode)
	        }
	        if len(call.Args) > 1 {
	                for _, n := range call.Args[1:] {
	                        code, _ := emitVerilogStatement(n)
	                        v.WriteString(code)
	                }
	        }
	case "greater-than":
	        v.WriteString(fmt.Sprintf("            // Evaluate greater-than: %s > %s\n", formatArg(call.Args[0]), formatArg(call.Args[1])))
	case "less-than":
	        v.WriteString(fmt.Sprintf("            // Evaluate less-than: %s < %s\n", formatArg(call.Args[0]), formatArg(call.Args[1])))
	case "equal":
	        v.WriteString(fmt.Sprintf("            // Evaluate equal: %s == %s\n", formatArg(call.Args[0]), formatArg(call.Args[1])))
	case "spawn-macro-cell":
	        if len(call.Args) >= 2 {
	                v.WriteString(fmt.Sprintf("            // TRIGGER DPR: Spawn %s and connect to %s\n", formatArg(call.Args[0]), formatArg(call.Args[1])))
	        }
	case "spawn-macro-cube":
	        if len(call.Args) >= 3 {
	                v.WriteString(fmt.Sprintf("            // TRIGGER DPR: Spawn Macro Cube %s and connect to %s\n", formatArg(call.Args[0]), formatArg(call.Args[2])))
	        }
	case "collapse-macro-cell":
	        if len(call.Args) >= 1 {
	                v.WriteString(fmt.Sprintf("            // TRIGGER DPR: Collapse %s to release dimensionality\n", formatArg(call.Args[0])))
	        }
	case "assign":		if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = %s;\n", dest, src))
		}
	case "bit-xor", "bit-and", "bit-or", "bit-shift-left", "bit-shift-right", "bit-clear", "bit-set":
		symMap := map[string]string{"bit-xor": "^", "bit-and": "&", "bit-or": "|", "bit-shift-left": "<<", "bit-shift-right": ">>", "bit-clear": "& ~", "bit-set": "|"}
		if len(call.Args) >= 3 {
			dest := formatArg(call.Args[0])
			left := formatArg(call.Args[1])
			right := formatArg(call.Args[2])
			if op == "bit-clear" {
				v.WriteString(fmt.Sprintf("            %s = %s & ~%s;\n", dest, left, right))
			} else {
				v.WriteString(fmt.Sprintf("            %s = %s %s %s;\n", dest, left, symMap[op], right))
			}
		}
	case "bit-not":
		if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = ~%s;\n", dest, src))
		}
	case "sum-split":
		if len(call.Args) >= 2 {
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            out = { (%s[23:12] + %s[11:0]), (%s[23:12] + %s[11:0]) };\n", src, src, src, src))
		}
	case "sum-pair":
		if len(call.Args) >= 2 {
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            out = %s[23:12] + %s[11:0];\n", src, src))
		}
	case "matrix-dot":
		if len(call.Args) >= 1 {
			mask := "24'h0"
			str, ok := call.Args[0].(*parser.StringLiteral)
			if ok {
				if strings.Contains(str.Value, ". . . . . .") {
					mask = "24'hAAA000"
					v.WriteString(fmt.Sprintf("            out = in & %s;\n", mask))
					v.WriteString("            // Matrix Dot Mask\n")
					return v.String(), nil
				} else if strings.Contains(str.Value, ". . . .") {
					mask = "24'hAA0000"
				}
			}
			v.WriteString(fmt.Sprintf("            out = in & %s;\n", mask))
		}
	case "quat-map":
		if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = %s;\n", dest, src))
		}
	case "torsional-pair":
		if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = ~%s;\n", dest, src))
		}
	case "fractal-zip":
		v.WriteString("            // FractalZip: High-density state compression initiated\n")
	case "photonic-stream":
		if len(call.Args) >= 2 {
			reg := call.Args[1].(*parser.Number).Value
			v.WriteString(fmt.Sprintf("            %s = photon_flux_reg_%s;\n", formatArg(call.Args[0]), reg))
			v.WriteString(fmt.Sprintf("            // Photonic Stream: Mapping physical reflection to register %s\n", reg))
		}
	case "swave-interaction":
		if len(call.Args) >= 2 {
			v.WriteString(fmt.Sprintf("            out = %s ^ swave_interaction_mask;\n", formatArg(call.Args[1])))
			v.WriteString("            // Topological Interaction: Applying field perturbations\n")
		}
	case "read-qudot":
		if len(call.Args) >= 2 {
			v.WriteString("            // ReadQuDot: Topological collapse of 784-qudit field\n")
			v.WriteString(fmt.Sprintf("            out = %s; // Collapsed state\n", formatArg(call.Args[1])))
		}
	case "read-sound-pixel":
		if len(call.Args) >= 2 {
			v.WriteString("            // ReadSoundPixel: Aural collapse of SPHY waveform (?)\n")
			v.WriteString(fmt.Sprintf("            out = %s; // Collapsed sound pixel\n", formatArg(call.Args[1])))
		}
	case "read-thought-pixel":
		if len(call.Args) >= 2 {
			v.WriteString("            // ReadThoughtPixel: Internal observation of the intent manifold (?)\n")
			v.WriteString(fmt.Sprintf("            out = %s; // Collapsed thought-sound pixel\n", formatArg(call.Args[1])))
		}
	case "transmit-qu-dot":
		if len(call.Args) >= 2 {
			v.WriteString("            // TransmitQuDot: Manifestation of digital intent into 784-qudit field (go)\n")
			v.WriteString(fmt.Sprintf("            out = %s; // Transmitted state\n", formatArg(call.Args[1])))
		}
	case "transmit-sound-pixel":
		if len(call.Args) >= 2 {
			v.WriteString("            // TransmitSoundPixel: Manifestation of inner voice into an aural fragment (go)\n")
			v.WriteString(fmt.Sprintf("            out = %s; // Transmitted sound pixel\n", formatArg(call.Args[1])))
		}
	case "mod-exp":
		v.WriteString("            // Auto-generated Modular Exponentiation Operator\n")
	case "stochastic-compensate":
		if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = ~%s;\n", dest, src))
		}
	case "temporal-void":
		v.WriteString("            // Auto-generated Temporal Void Operator\n")
	case "phi-scale":
		if len(call.Args) >= 3 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			scale := formatArg(call.Args[2])
			if scale == "-1.0" {
				v.WriteString(fmt.Sprintf("            %s = (%s * -1024) >> 10;\n", dest, src))
			} else {
				v.WriteString(fmt.Sprintf("            %s = (%s * 1657) >> 10;\n", dest, src))
			}
		} else if len(call.Args) >= 2 {
			dest := formatArg(call.Args[0])
			src := formatArg(call.Args[1])
			v.WriteString(fmt.Sprintf("            %s = (%s * 1657) >> 10;\n", dest, src))
		} else if len(call.Args) >= 1 {
			dest := formatArg(call.Args[0])
			v.WriteString(fmt.Sprintf("            %s = (%s * 1657) >> 10;\n", dest, dest))
		}
	}

	return v.String(), nil
}
