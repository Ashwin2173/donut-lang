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
}

type BlockStatement struct {
	body []Statement
	line uint
}

type ReturnStatement struct {
	Statement
	expression *ExpressionStatement
}

func (BlockStatement) isStatement()      {}
func (ReturnStatement) isStatement()     {}
func (FunctionStatement) isStatement()   {}
func (ExpressionStatement) isStatement() {}

type Identifier struct {
	name string
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
}

func (Identifier) isExpression()       {}
func (FloatLiteral) isExpression()     {}
func (StringLiteral) isExpression()    {}
func (BooleanLiteral) isExpression()   {}
func (IntegerLiteral) isExpression()   {}
func (UnaryExpression) isExpression()  {}
func (BinaryExpression) isExpression() {}
