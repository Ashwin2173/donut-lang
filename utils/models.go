package utils

import (
	"encoding/binary"
)

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
	expression Expression
	line       uint
}

func (BlockStatement) isStatement()      {}
func (ReturnStatement) isStatement()     {}
func (FunctionStatement) isStatement()   {}
func (ExpressionStatement) isStatement() {}

type NullLiteral struct{}

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
func (NullLiteral) isExpression()      {}
func (FloatLiteral) isExpression()     {}
func (StringLiteral) isExpression()    {}
func (BooleanLiteral) isExpression()   {}
func (IntegerLiteral) isExpression()   {}
func (UnaryExpression) isExpression()  {}
func (BinaryExpression) isExpression() {}

type ByteCode struct {
	byteCode []byte
}

func NewByteCode() ByteCode {
	return ByteCode{
		byteCode: []byte{},
	}
}

func (b *ByteCode) length() int {
	return len(b.byteCode)
}

func (b *ByteCode) writeU1(value byte) {
	b.byteCode = append(b.byteCode, value)
}

func (b *ByteCode) writeU2(value uint16) {
	b.byteCode = binary.BigEndian.AppendUint16(b.byteCode, value)
}

func (b *ByteCode) writeU4(value uint32) {
	b.byteCode = binary.BigEndian.AppendUint32(b.byteCode, value)
}

func (b *ByteCode) dump(byteCode ByteCode) {
	b.byteCode = append(b.byteCode, byteCode.byteCode...)
}

func (b *ByteCode) dumpRaw(bytes []byte) {
	b.byteCode = append(b.byteCode, bytes...)
}
