package utils

import (
	"fmt"
)

type TypeChecker struct {
	fileName     string
	program      Program
	functions    map[string][]string
	typeRegister map[string]TypeId
}

var binaryRules = map[Pair[TypeId, TypeId]]BinaryTypeResult{
	{T_INT, T_INT}:     {T_INT, nil, nil},
	{T_FLOAT, T_FLOAT}: {T_FLOAT, nil, nil},
}

func NewTypeChecker(fileName string, program Program) *TypeChecker {
	defaultType := map[string]TypeId{
		"int":    T_INT,
		"float":  T_FLOAT,
		"string": T_STRING,
	}
	return &TypeChecker{
		fileName:     fileName,
		program:      program,
		functions:    make(map[string][]string),
		typeRegister: defaultType,
	}
}

func (t *TypeChecker) resolve() error {
	for _, fnStmt := range t.program.functions {
		err := t.handleFnBodyStatement(fnStmt.body)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TypeChecker) handleFnBodyStatement(body BlockStatement) error {
	for _, stmt := range body.body {
		switch stmt := stmt.(type) {
		case *ExpressionStatement:
			_, err := t.handleExpression(stmt.expression)
			if err != nil {
				return err
			}
		case *PrintStatement:
			// test
			_, err := t.handleExpression(stmt.expression)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("DEV-INFO: Skipped %T Type for typechecking\n", stmt)
		}
	}
	return nil
}

func (t *TypeChecker) handleExpression(expr Expression) (TypeId, error) {
	switch expr := expr.(type) {
	case *IntegerLiteral:
		return T_INT, nil
	case *FloatLiteral:
		return T_FLOAT, nil
	case *StringLiteral:
		return T_STRING, nil
	case *BooleanLiteral:
		return T_BOOLEAN, nil
	case *BinaryExpression:
		return t.handleBinaryExpression(expr)
	}
	return T_UNKNOWN, nil
}

func (t *TypeChecker) handleBinaryExpression(expr *BinaryExpression) (TypeId, error) {
	left, err := t.handleExpression(expr.left)
	if err != nil {
		return T_UNKNOWN, err
	}
	right, err := t.handleExpression(expr.right)
	if err != nil {
		return T_UNKNOWN, err
	}
	pair := Pair[TypeId, TypeId]{left, right}
	resolvedInfo, exists := binaryRules[pair]
	if exists {
		expr.resolvedInfo = resolvedInfo
		return resolvedInfo.result, nil
	}
	return T_UNKNOWN, t.fail("Invalided Binary Expression", expr.line)
}

func (t *TypeChecker) fail(message string, line uint) error {
	return error(CompileError{
		Message:  message,
		FileName: t.fileName,
		Line:     line,
		Col:      0,
	})
}
