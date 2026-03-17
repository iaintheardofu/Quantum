package parser

import (
	"fmt"
)

type Parser struct {
	lexer *Lexer
	cur   Token
	peek  Token
}

func NewParser(input string) *Parser {
	p := &Parser{lexer: NewLexer(input)}
	p.advance()
	p.advance()
	return p
}

func (p *Parser) advance() {
	p.cur = p.peek
	p.peek = p.lexer.NextToken()
}

func (p *Parser) Parse() (*Program, error) {
	prog := &Program{}

	for p.cur.Type != TokenEOF {
		node, err := p.parseForm()
		if err != nil {
			return nil, err
		}
		if node != nil {
			prog.Body = append(prog.Body, node)
		}
	}

	return prog, nil
}

func (p *Parser) parseForm() (Node, error) {
	switch p.cur.Type {
	case TokenLParen:
		return p.parseListForm()
	case TokenLBrack:
		return p.parseBracketList()
	case TokenIdent:
		node := &Identifier{Name: p.cur.Value}
		p.advance()
		return node, nil
	case TokenNum:
		node := &Number{Value: p.cur.Value}
		p.advance()
		return node, nil
	case TokenStr:
		node := &StringLiteral{Value: p.cur.Value}
		p.advance()
		return node, nil
	case TokenEOF:
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s at %d:%d", p.cur.Value, p.cur.Line, p.cur.Col)
	}
}

func (p *Parser) parseListForm() (Node, error) {
	p.advance() // skip '('
	
	if p.cur.Type == TokenEOF {
		return nil, fmt.Errorf("unexpected EOF")
	}

	if p.cur.Type != TokenIdent {
		return nil, fmt.Errorf("expected identifier after '(', got %s", p.cur.Value)
	}

	callee := p.cur.Value
	p.advance()

	if callee == "ns" {
		return p.parseNamespace()
	} else if callee == "defn" || callee == "defn-ai" || callee == "defn-fractal" || callee == "defn-ui" {
		return p.parseDefn(callee)
	}

	call := &Call{Callee: callee}
	
	for p.cur.Type != TokenRParen && p.cur.Type != TokenEOF {
		arg, err := p.parseForm()
		if err != nil {
			return nil, err
		}
		call.Args = append(call.Args, arg)
	}

	if p.cur.Type != TokenRParen {
		return nil, fmt.Errorf("expected ')' closing call to %s", callee)
	}
	p.advance() // skip ')'

	return call, nil
}

func (p *Parser) parseNamespace() (*Namespace, error) {
	if p.cur.Type != TokenIdent {
		return nil, fmt.Errorf("expected namespace name")
	}
	ns := &Namespace{Name: p.cur.Value}
	p.advance()

	// skip rest of ns declaration
	for p.cur.Type != TokenRParen && p.cur.Type != TokenEOF {
		p.advance()
	}
	if p.cur.Type == TokenRParen {
		p.advance()
	}
	return ns, nil
}

func (p *Parser) parseDefn(defType string) (*Defn, error) {
	if p.cur.Type != TokenIdent {
		return nil, fmt.Errorf("expected function name")
	}
	d := &Defn{
		Name:      p.cur.Value,
		IsAI:      defType == "defn-ai",
		IsFractal: defType == "defn-fractal",
	}
	p.advance()

	if p.cur.Type == TokenLBrack {
		p.advance() // skip '['
		for p.cur.Type != TokenRBrack && p.cur.Type != TokenEOF {
			if p.cur.Type == TokenIdent {
				d.Params = append(d.Params, p.cur.Value)
			}
			p.advance()
		}
		if p.cur.Type == TokenRBrack {
			p.advance()
		}
	}

	if d.IsAI && p.cur.Type == TokenStr {
		d.Intent = p.cur.Value
		p.advance()
	}

	for p.cur.Type != TokenRParen && p.cur.Type != TokenEOF {
		bodyNode, err := p.parseForm()
		if err != nil {
			return nil, err
		}
		d.Body = append(d.Body, bodyNode)
	}

	if p.cur.Type == TokenRParen {
		p.advance()
	}

	return d, nil
}

func (p *Parser) parseBracketList() (*List, error) {
	p.advance() // skip '['
	l := &List{}
	for p.cur.Type != TokenRBrack && p.cur.Type != TokenEOF {
		elem, err := p.parseForm()
		if err != nil {
			return nil, err
		}
		l.Elements = append(l.Elements, elem)
	}
	if p.cur.Type != TokenRBrack {
		return nil, fmt.Errorf("expected ']'")
	}
	p.advance()
	return l, nil
}
