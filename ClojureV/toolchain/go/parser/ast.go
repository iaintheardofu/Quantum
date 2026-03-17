package parser

// NodeType represents the type of an AST node
type NodeType string

const (
	NodeProgram    NodeType = "Program"
	NodeNamespace  NodeType = "Namespace"
	NodeDefn       NodeType = "Defn"
	NodeCall       NodeType = "Call"
	NodeIdentifier NodeType = "Identifier"
	NodeNumber     NodeType = "Number"
	NodeString     NodeType = "String"
	NodeList       NodeType = "List" // e.g. [clk rst_n in]
)

// Node is the base interface for all AST nodes
type Node interface {
	Type() NodeType
	String() string
}

// Program is the root of the AST
type Program struct {
	Body []Node
}

func (p *Program) Type() NodeType { return NodeProgram }
func (p *Program) String() string { return "Program" }

// Namespace declaration (ns ...)
type Namespace struct {
	Name string
}

func (n *Namespace) Type() NodeType { return NodeNamespace }
func (n *Namespace) String() string { return "(ns " + n.Name + ")" }

// Defn represents a function definition (defn or defn-ai)
type Defn struct {
	Name      string
	IsAI      bool
	IsFractal bool
	Params    []string
	Intent    string // Only used if IsAI is true
	Body      []Node
}

func (d *Defn) Type() NodeType { return NodeDefn }
func (d *Defn) String() string {
	typ := "defn"
	if d.IsAI {
		typ = "defn-ai"
	}
	return "(" + typ + " " + d.Name + " ...)"
}

// Call represents a function call or list expression like (qurq/bit-xor out in 0x0)
type Call struct {
	Callee string
	Args   []Node
}

func (c *Call) Type() NodeType { return NodeCall }
func (c *Call) String() string { return "(" + c.Callee + " ...)" }

// Identifier represents a variable or keyword
type Identifier struct {
	Name string
}

func (i *Identifier) Type() NodeType { return NodeIdentifier }
func (i *Identifier) String() string { return i.Name }

// Number represents a numeric literal
type Number struct {
	Value string
}

func (n *Number) Type() NodeType { return NodeNumber }
func (n *Number) String() string { return n.Value }

// StringLiteral represents a string literal
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() NodeType { return NodeString }
func (s *StringLiteral) String() string { return "\"" + s.Value + "\"" }

// List represents a bracketed list [...]
type List struct {
	Elements []Node
}

func (l *List) Type() NodeType { return NodeList }
func (l *List) String() string { return "[...]" }
