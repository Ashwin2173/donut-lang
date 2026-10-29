package utils

type Word struct {
	Raw   string
	Type  TokenType
	Line  uint
	Start uint
	End   uint
}

type Statement interface{ isStatement() }
type Expression interface{ isExpression() }

type Program struct {
	functions []FunctionStatement
}

type ExpressionStatement struct {
	expression Expression
	line       uint
}

type FunctionStatement struct {
	name      string
	arguments []Identifier
	body      BlockStatement
	maxLocals int
}

type BlockStatement struct {
	body []Statement
	line uint
}

type ReturnStatement struct {
	expression Expression
	line       uint
}

type PrintStatement struct {
	// temp
	expression   Expression
	resolvedArgs string
	line         uint
}

func (BlockStatement) isStatement()      {}
func (PrintStatement) isStatement()      {}
func (ReturnStatement) isStatement()     {}
func (FunctionStatement) isStatement()   {}
func (ExpressionStatement) isStatement() {}

type NullLiteral struct{}

type Identifier struct {
	name     string
	line     uint
	variable *Variable
}

type Variable struct {
	varType TypeId
	index   int
}

type IntegerLiteral struct {
	value int64
}

type FloatLiteral struct {
	value float64
}

type StringLiteral struct {
	value string
}

type BooleanLiteral struct {
	value bool
}

type UnaryExpression struct {
	unaryType TokenType
	value     Expression
}

type BinaryExpression struct {
	left       Expression
	binaryType TokenType
	right      Expression
	line       uint
	typeInfo   BinaryTypeValue
}

type AssignmentExpression BinaryExpression

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

type BinaryTypeKey struct {
	left      TypeId
	operation TokenType
	right     TypeId
}

type BinaryTypeValue struct {
	result                         TypeId
	operationType, leftTo, rightTo BinOpTypeId
}
