package utils

import (
	"fmt"
	"slices"
	"strconv"
)

type Parser struct {
	fileName string
	tokens   *Generator[Word]
}

func NewParser(fileName, source string) *Parser {
	tokens := Tokenize(source)
	fmt.Println("Tokens:", tokens)
	return &Parser{
		fileName: fileName,
		tokens:   NewGenerator(tokens),
	}
}

func (p *Parser) ParseProgram() (*Program, error) {
	// mocking program for testing
	expression, err := p.scanExpression()
	if err != nil {
		return nil, err
	}
	return &Program{
		functions: []FunctionStatement{
			FunctionStatement{
				name:      "main",
				arguments: []Identifier{},
				body: BlockStatement{
					body: []Statement{
						expression,
					},
				},
			},
		},
	}, nil
}

func (p *Parser) scanExpression() (*ExpressionStatement, error) {
	exp, err := p.scanLogicalOr()
	if err != nil {
		return nil, err
	}
	return &ExpressionStatement{
		expression: exp,
	}, nil
}

func (p *Parser) scanLogicalOr() (Expression, error) {
	return p.makeBinary(p.scanLogicalAnd, K_OR)
}

func (p *Parser) scanLogicalAnd() (Expression, error) {
	return p.makeBinary(p.scanEquality, K_AND)
}

func (p *Parser) scanEquality() (Expression, error) {
	return p.makeBinary(p.scanComparison, EQUAL_EQUAL, BANG_EQUAL)
}

func (p *Parser) scanComparison() (Expression, error) {
	return p.makeBinary(p.scanTerm, LESSER, LESSER_EQUALS, GREATER, GREATER_EQUALS)
}

func (p *Parser) scanTerm() (Expression, error) {
	return p.makeBinary(p.scanFactor, PLUS, MINUS)
}

func (p *Parser) scanFactor() (Expression, error) {
	return p.makeBinary(p.scanUnary, STAR, SLASH)
}

func (p *Parser) scanUnary() (Expression, error) {
	token, isEnd := p.tokens.peek()
	if !isEnd && (token.Type == MINUS || token.Type == K_NOT) {
		p.tokens.next()
		right, err := p.scanPostfix()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			unaryType: token.Type,
			value:     right,
		}, nil
	}
	return p.scanPostfix()
}

func (p *Parser) scanPostfix() (Expression, error) {
	return p.scanPrimary()
}

func (p *Parser) scanPrimary() (Expression, error) {
	token, isEnd := p.tokens.next()
	if isEnd {
		return nil, p.fail("Expected expression but reached end of file")
	}
	switch token.Type {
	case INTEGER:
		value, _ := strconv.Atoi(token.Raw)
		return &IntegerLiteral{value: int64(value)}, nil
	case FLOAT:
		value, _ := strconv.ParseFloat(token.Raw, 64)
		return &FloatLiteral{value: value}, nil
	case STRING:
		return &StringLiteral{value: token.Raw}, nil
	case K_TRUE, K_FALSE:
		value := false
		if token.Type == K_TRUE {
			value = true
		}
		return &BooleanLiteral{value: value}, nil
	case IDENTIFIER:
		return &Identifier{name: token.Raw}, nil
	}
	return nil, p.fail("Invalid Syntax")
}

func (p *Parser) makeBinary(callNext func() (Expression, error), operators ...TokenType) (Expression, error) {
	left, err := callNext()
	if err != nil {
		return nil, err
	}
	for {
		token, isEnd := p.tokens.peek()
		if isEnd || !slices.Contains(operators, token.Type) {
			return left, nil
		}
		p.tokens.next()
		right, err := callNext()
		if err != nil {
			return nil, err
		}
		left = &BinaryExpression{
			left:       left,
			binaryType: token.Type,
			right:      right,
		}
	}
}

func (p *Parser) expectToken(tokenType TokenType) error {
	token, isEnd := p.tokens.next()
	if isEnd {
		message := fmt.Sprintf("Expected %s but reached end of file", GROUP_NAMES[tokenType])
		return p.fail(message)
	}
	if token.Type != tokenType {
		message := fmt.Sprintf("Expected %s but got %s", GROUP_NAMES[tokenType], GROUP_NAMES[token.Type])
		return p.fail(message)
	}
	return nil
}

func (p *Parser) fail(message string) error {
	token, _ := p.tokens.peek()
	return error(CompileError{
		Message:  message,
		FileName: p.fileName,
		Line:     token.Line,
		Col:      token.Start,
	})
}
