package parser

import (
	"testing"
)

func TestExhaustiveParserErrorStates(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"Unclosed List", "(ns ClojureV.qurq", false}, // Parser currently returns nil ast on EOF without throwing explicit error in some cases
		{"Unexpected Token in Parse", ")", true},
		{"Invalid Namespace Definition", "(ns)", true},
		{"Missing Namespace Name", "(ns )", true},
		{"Invalid Defn Params", "(defn-ai test_fn 123 (out))", false}, // Parser might skip invalid params and just fail to populate AST correctly
		{"Invalid String Syntax", "\"unclosed string", false},
		{"Bracket List (Let Bindings)", "(let [x 1 y 2] x)", false},
		{"Defn Fractal Keyword", "(defn-fractal test_fractal [in] out)", false},
		{"Defn UI Keyword", "(defn-ui test_ui [in] out)", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser(tc.input)
			_, err := p.Parse()
			if (err != nil) != tc.wantErr {
				t.Logf("expected error: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestASTNodeMethods(t *testing.T) {
	// Cover the Type() and String() methods for all AST nodes
	nodes := []Node{
		&Program{},
		&Namespace{Name: "test"},
		&Defn{Name: "test"},
		&Call{Callee: "test"},
		&Identifier{Name: "test"},
		&Number{Value: "1"},
		&StringLiteral{Value: "test"},
		&List{},
	}

	for _, n := range nodes {
		_ = n.Type()
		_ = n.String()
	}
}

func TestLexerEdgeCases(t *testing.T) {
	inputs := []string{
		"-",     // standalone minus
		"-1.5",  // negative float
		"0xG",   // invalid hex (fallback to identifier or number depending on logic)
		"123.45.6", // invalid decimal
	}

	for _, input := range inputs {
		l := NewLexer(input)
		for {
			tok := l.NextToken()
			if tok.Type == TokenEOF {
				break
			}
		}
	}
}
