package parser

import (
	"fmt"
	"slices"
	"strconv"

	"donut/lang"
)

type Parser struct {
	fileName string
	tokens   *Generator[Word]
}

func New(fileName, source string) *Parser {
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
		if token.Type == k_function {
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
	return p.analyzeProgram(Program{
		functions: functions,
	})
}

func (p *Parser) analyzeProgram(program Program) (*Program, error) {
	analyzer := NewAnalyzer(p.fileName, program)
	err := analyzer.Analyze()
	if err != nil {
		return nil, err
	}
	return &analyzer.program, nil
}

func (p *Parser) scanFunctionDef() (*FunctionStatement, error) {
	name, err := p.expectToken(t_identifier)
	if err != nil {
		return nil, p.fail("Expected a function name")
	}
	_, err = p.expectToken(t_openParen)
	if err != nil {
		return nil, p.fail("Expected '(' after the function name")
	}
	_, err = p.expectToken(t_closeParen)
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
	openBrace, err := p.expectToken(t_openBrace)
	if err != nil {
		return nil, p.fail("Expected '{' at the start of block")
	}
	body := []Statement{}
	for {
		token, isEnd := p.tokens.peek()
		if isEnd {
			return nil, p.fail("Reached End of File while parsing block; Expected '}'")
		}
		if token.Type == t_closeBrace {
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
	case k_return:
		return p.scanReturnStatement()
	case k_while:
		return p.scanWhileStatement()
	case k_if:
		return p.scanIfStatement()
	case t_openBrace:
		return p.scanBlockStatement()
	case k_print: // temp code
		return p.scanPrintStatement()
	}
	return p.scanExpressionStatement()
}

func (p *Parser) scanIfStatement() (*IfStatement, error) {
	token, _ := p.tokens.next()
	expr, err := p.scanExpression()
	if err != nil {
		return nil, err
	}
	token, isEnd := p.tokens.peek()
	if isEnd {
		return nil, p.fail("Required a statement for if stmt, but reached end of file")
	}
	consequent, err := p.scanLocalStatement(token)
	if err != nil {
		return nil, err
	}
	var alternate Statement
	token, isEnd = p.tokens.peek()
	if !isEnd && token.Type == k_else {
		p.tokens.next()
		token, isEnd = p.tokens.peek()
		if isEnd {
			return nil, p.fail("Required a statement for else stmt, but reached end of file")
		}
		alternate, err = p.scanLocalStatement(token)
		if err != nil {
			return nil, err
		}
	}
	return &IfStatement{
		test:       expr,
		consequent: consequent,
		alternate:  alternate,
		line:       token.Line,
	}, nil
}

func (p *Parser) scanWhileStatement() (*WhileStatement, error) {
	token, _ := p.tokens.next()
	expr, err := p.scanExpression()
	if err != nil {
		return nil, err
	}
	token, isEnd := p.tokens.peek()
	if isEnd {
		return nil, p.fail("Required a statement for while loop, but reached end of file")
	}
	bodystmt, err := p.scanLocalStatement(token)
	if err != nil {
		return nil, err
	}
	return &WhileStatement{
		test: expr,
		body: bodystmt,
		line: token.Line,
	}, nil
}

func (p *Parser) scanPrintStatement() (*PrintStatement, error) {
	// temp code
	ret, _ := p.tokens.next()
	token, isEnd := p.tokens.peek()
	if isEnd {
		return nil, p.fail("Expected an expression or semicolon after return")
	}
	var expression Expression
	if token.Type == t_semicolon {
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
	if token.Type == t_semicolon {
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
	if isEnd || token.Type != t_semicolon {
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
	expr, err := p.makeBinary(p.scanLogicalOr, t_equal)
	if err != nil {
		return nil, err
	}
	binExpr, ok := expr.(*BinaryExpression)
	if ok && binExpr.binaryType == t_equal {
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
	return p.makeBinary(p.scanLogicalAnd, k_or)
}

func (p *Parser) scanLogicalAnd() (Expression, error) {
	return p.makeBinary(p.scanEquality, k_and)
}

func (p *Parser) scanEquality() (Expression, error) {
	return p.makeBinary(p.scanComparison, t_equalEqual, t_bangEqual)
}

func (p *Parser) scanComparison() (Expression, error) {
	return p.makeBinary(p.scanTerm, t_lesser, t_lesserEquals, t_greater, t_greaterEquals)
}

func (p *Parser) scanTerm() (Expression, error) {
	return p.makeBinary(p.scanFactor, t_plus, t_minus)
}

func (p *Parser) scanFactor() (Expression, error) {
	return p.makeBinary(p.scanUnary, t_star, t_slash, t_mod)
}

func (p *Parser) scanUnary() (Expression, error) {
	token, isEnd := p.tokens.peek()
	if !isEnd && (token.Type == t_minus || token.Type == k_not) {
		p.tokens.next()
		right, err := p.scanUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryExpression{
			unaryType: token.Type,
			value:     right,
			line:      token.Line,
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
	case t_integerLit:
		value, _ := strconv.Atoi(token.Raw)
		return &IntegerLiteral{value: int64(value)}, nil
	case t_floatLit:
		value, _ := strconv.ParseFloat(token.Raw, 64)
		return &FloatLiteral{value: value}, nil
	case t_stringLit:
		return &StringLiteral{value: token.Raw[1 : len(token.Raw)-1]}, nil
	case t_identifier:
		return &Identifier{name: token.Raw, line: token.Line}, nil
	case t_openParen:
		expr, err := p.scanExpression()
		if err != nil {
			return nil, err
		}
		token, isEnd := p.tokens.next()
		if isEnd || token.Type != t_closeParen {
			return nil, p.fail("Unclosed parentheses; Expected ) to complete the expression")
		}
		return expr, nil
	case k_true, k_false:
		value := false
		if token.Type == k_true {
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
		message := fmt.Sprintf("Expected %s but reached end of file", groupNames[tokenType])
		return token, p.fail(message)
	}
	if token.Type != tokenType {
		message := fmt.Sprintf("Expected %s but got %s", groupNames[tokenType], groupNames[token.Type])
		return token, p.fail(message)
	}
	return token, nil
}

func (p *Parser) fail(message string) error {
	token, _ := p.tokens.peek()
	return error(lang.CompileError{
		Message:  message,
		FileName: p.fileName,
		Line:     token.Line,
		Col:      token.Column,
	})
}
