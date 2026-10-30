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
	// fmt.Println("Tokens:", tokens)
	return &Parser{
		fileName: fileName,
		tokens:   NewGenerator(tokens),
	}
}

func (p *Parser) ParseProgram() (*Program, error) {
	if len(p.tokens.data) == 0 {
		return nil, p.fail("provided file is empty; start by writing 'function' keyword")
	}
	functions := []FunctionStatement{}
	for {
		token, isEnd := p.tokens.peek()
		if isEnd {
			break
		}
		if token.Type == K_FUNCTION {
			p.tokens.next()
			function, err := p.scanFunctionDef()
			if err != nil {
				return nil, err
			}
			functions = append(functions, *function)
		} else {
			return nil, p.fail("Expected function definition in global scope")
		}
	}
	return &Program{
		functions: functions,
	}, nil
}

func (p *Parser) scanFunctionDef() (*FunctionStatement, error) {
	name, err := p.expectToken(IDENTIFIER)
	if err != nil {
		return nil, p.fail("Expected a function name")
	}
	_, err = p.expectToken(OPEN_PARAM)
	if err != nil {
		return nil, p.fail("Expected '(' after the function name")
	}
	_, err = p.expectToken(CLOSE_PARAM)
	if err != nil {
		return nil, p.fail("Expected ')'")
	}
	body, err := p.scanBlockStatement()
	if err != nil {
		return nil, err
	}
	return &FunctionStatement{
		name: name.Raw,
		body: *body,
	}, nil
}

func (p *Parser) scanBlockStatement() (*BlockStatement, error) {
	openBrace, err := p.expectToken(OPEN_BRACE)
	if err != nil {
		return nil, p.fail("Expected '{' at the start of block")
	}
	body := []Statement{}
	for {
		token, isEnd := p.tokens.peek()
		if isEnd {
			return nil, p.fail("Reached End of File while parsing block; Expected '}'")
		}
		if token.Type == CLOSE_BRACE {
			p.tokens.next()
			return &BlockStatement{
				body: body,
				line: openBrace.Line,
			}, nil
		}
		statement, err := p.scanLocalStatement(token)
		if err != nil {
			return nil, err
		}
		body = append(body, statement)
	}
}

func (p *Parser) scanLocalStatement(token Word) (Statement, error) {
	switch token.Type {
	case K_RETURN:
		return p.scanReturnStatement()
	case K_PRINT: // temp code
		return p.scanPrintStatement()
	}
	return p.scanExpressionStatement()
}

func (p *Parser) scanPrintStatement() (*PrintStatement, error) {
	// temp code
	ret, _ := p.tokens.next()
	token, isEnd := p.tokens.peek()
	if isEnd {
		return nil, p.fail("Expected an expression or semicolon after return")
	}
	var expression Expression
	if token.Type == SEMICOLON {
		p.tokens.next()
		expression = NullLiteral{}
	} else {
		expstmt, err := p.scanExpressionStatement()
		if err != nil {
			return nil, err
		}
		expression = expstmt.expression
	}
	return &PrintStatement{
		expression: expression,
		line:       ret.Line,
	}, nil
}

func (p *Parser) scanReturnStatement() (*ReturnStatement, error) {
	ret, _ := p.tokens.next()
	token, isEnd := p.tokens.peek()
	if isEnd {
		return nil, p.fail("Expected an expression or semicolon after return")
	}
	var expression Expression
	if token.Type == SEMICOLON {
		p.tokens.next()
		expression = NullLiteral{}
	} else {
		expstmt, err := p.scanExpressionStatement()
		if err != nil {
			return nil, err
		}
		expression = expstmt.expression
	}
	return &ReturnStatement{
		expression: expression,
		line:       ret.Line,
	}, nil
}

func (p *Parser) scanExpressionStatement() (*ExpressionStatement, error) {
	token, _ := p.tokens.peek()
	exp, err := p.scanExpression()
	if err != nil {
		return nil, err
	}
	token, isEnd := p.tokens.next()
	if isEnd || token.Type != SEMICOLON {
		return nil, p.fail("Expected Semicolon after expression")
	}
	return &ExpressionStatement{
		expression: exp,
		line:       token.Line,
	}, nil
}

func (p *Parser) scanExpression() (Expression, error) {
	return p.scanAssignment()
}

func (p *Parser) scanAssignment() (Expression, error) {
	expr, err := p.makeBinary(p.scanLogicalOr, EQUAL)
	if err != nil {
		return nil, err
	}
	binExpr, ok := expr.(*BinaryExpression)
	if ok && binExpr.binaryType == EQUAL {
		switch binExpr.left.(type) {
		case *Identifier:
			assignExpr := AssignmentExpression(*binExpr)
			return &assignExpr, nil
		default:
			return nil, p.fail("Expected Identifier on left side of Assignment")
		}
	}
	return expr, nil
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
	return p.makeBinary(p.scanUnary, STAR, SLASH, MOD)
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
		return &StringLiteral{value: token.Raw[1 : len(token.Raw)-1]}, nil
	case IDENTIFIER:
		return &Identifier{name: token.Raw, line: token.Line}, nil
	case OPEN_PARAM:
		expr, err := p.scanExpression()
		if err != nil {
			return nil, err
		}
		token, isEnd := p.tokens.next()
		if isEnd || token.Type != CLOSE_PARAM {
			return nil, p.fail("Unclosed parentheses; Expected ) to complete the expression")
		}
		return expr, nil
	case K_TRUE, K_FALSE:
		value := false
		if token.Type == K_TRUE {
			value = true
		}
		return &BooleanLiteral{value: value}, nil
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
			line:       token.Line,
		}
	}
}

func (p *Parser) expectToken(tokenType TokenType) (Word, error) {
	token, isEnd := p.tokens.next()
	if isEnd {
		message := fmt.Sprintf("Expected %s but reached end of file", GROUP_NAMES[tokenType])
		return token, p.fail(message)
	}
	if token.Type != tokenType {
		message := fmt.Sprintf("Expected %s but got %s", GROUP_NAMES[tokenType], GROUP_NAMES[token.Type])
		return token, p.fail(message)
	}
	return token, nil
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
