package parser

type TypeId int32
type ExpOpTypeId int

type Statement interface{ isStatement() }
type Expression interface{ isExpression() }

type Program struct {
	Functions []FunctionStatement
}

type ExpressionStatement struct {
	Expression Expression
	line       uint
}

type FunctionStatement struct {
	Name       string
	Arguments  []Identifier
	ReturnType string
	Body       BlockStatement
	line       uint
}

type BlockStatement struct {
	Body []Statement
	line uint
}

type ReturnStatement struct {
	Expression Expression
	ReturnType TypeId
	line       uint
}

type WhileStatement struct {
	Test Expression
	Body Statement
	line uint
}

type IfStatement struct {
	Test       Expression
	Consequent Statement
	Alternate  Statement
	line       uint
}

type PrintStatement struct {
	// temp
	Expression   Expression
	ResolvedArgs string
	line         uint
}

func (IfStatement) isStatement()         {}
func (BlockStatement) isStatement()      {}
func (WhileStatement) isStatement()      {}
func (PrintStatement) isStatement()      {}
func (ReturnStatement) isStatement()     {}
func (FunctionStatement) isStatement()   {}
func (ExpressionStatement) isStatement() {}

type NullLiteral struct{}

type Identifier struct {
	Name     string
	line     uint
	Variable *Variable
}

type Variable struct {
	VarType TypeId
}

type IntegerLiteral struct {
	Value int64
}

type FloatLiteral struct {
	Value float64
}

type StringLiteral struct {
	Value string
}

type BooleanLiteral struct {
	Value bool
}

type UnaryExpression struct {
	unaryType TokenType
	Value     Expression
	TypeInfo  ExpOpTypeId
	line      uint
}

type BinaryExpression struct {
	Left       Expression
	binaryType TokenType
	Right      Expression
	line       uint
	TypeInfo   BinaryTypeValue
}

type AssignmentExpression struct {
	IsNew bool
	BinaryExpression
}

func (Identifier) isExpression()           {}
func (NullLiteral) isExpression()          {}
func (FloatLiteral) isExpression()         {}
func (StringLiteral) isExpression()        {}
func (BooleanLiteral) isExpression()       {}
func (IntegerLiteral) isExpression()       {}
func (UnaryExpression) isExpression()      {}
func (BinaryExpression) isExpression()     {}
func (AssignmentExpression) isExpression() {}

type Pair[L any, R any] struct {
	left  L
	right R
}
