package utils

import (
	"fmt"
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
	exp, err := p.scanPrimary()
	if err != nil {
		return nil, err
	}
	return &ExpressionStatement{
		expression: exp,
	}, nil
}

func (p *Parser) scanPrimary() (Expression, error) {
	token, isEnd := p.tokens.next()
	if isEnd {
		return nil, p.throw("Expected expression but reached end of file")
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
	return nil, p.throw("Invalid Syntax")
}

func (p *Parser) expectToken(tokenType TokenType) error {
	token, isEnd := p.tokens.next()
	if isEnd {
		message := fmt.Sprintf("Expected %s but reached end of file", GROUP_NAMES[tokenType])
		return p.throw(message)
	}
	if token.Type != tokenType {
		message := fmt.Sprintf("Expected %s but got %s", GROUP_NAMES[tokenType], GROUP_NAMES[token.Type])
		return p.throw(message)
	}
	return nil
}

func (p *Parser) throw(message string) error {
	token, _ := p.tokens.peek()
	return error(CompileError{
		Message:  message,
		FileName: p.fileName,
		Line:     token.Line,
		Col:      token.Start,
	})
}
