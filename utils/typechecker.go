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

var binaryRules = map[BinaryTypeKey]BinaryTypeValue{
	{T_INT, PLUS, T_INT}:  {T_INT, (*ByteCode).opIAdd, nil, nil},
	{T_INT, MINUS, T_INT}: {T_INT, (*ByteCode).opISub, nil, nil},
	{T_INT, STAR, T_INT}:  {T_INT, (*ByteCode).opIMul, nil, nil},
	{T_INT, SLASH, T_INT}: {T_INT, (*ByteCode).opIDiv, nil, nil},

	{T_INT, PLUS, T_FLOAT}:  {T_FLOAT, (*ByteCode).opFAdd, (*ByteCode).opI2F, nil},
	{T_INT, MINUS, T_FLOAT}: {T_FLOAT, (*ByteCode).opFSub, (*ByteCode).opI2F, nil},
	{T_INT, STAR, T_FLOAT}:  {T_FLOAT, (*ByteCode).opFMul, (*ByteCode).opI2F, nil},
	{T_INT, SLASH, T_FLOAT}: {T_FLOAT, (*ByteCode).opFDiv, (*ByteCode).opI2F, nil},

	{T_FLOAT, PLUS, T_FLOAT}:  {T_FLOAT, (*ByteCode).opFAdd, nil, nil},
	{T_FLOAT, MINUS, T_FLOAT}: {T_FLOAT, (*ByteCode).opFSub, nil, nil},
	{T_FLOAT, STAR, T_FLOAT}:  {T_FLOAT, (*ByteCode).opFMul, nil, nil},
	{T_FLOAT, SLASH, T_FLOAT}: {T_FLOAT, (*ByteCode).opFDiv, nil, nil},

	{T_FLOAT, PLUS, T_INT}:  {T_FLOAT, (*ByteCode).opFAdd, nil, (*ByteCode).opI2F},
	{T_FLOAT, MINUS, T_INT}: {T_FLOAT, (*ByteCode).opFSub, nil, (*ByteCode).opI2F},
	{T_FLOAT, STAR, T_INT}:  {T_FLOAT, (*ByteCode).opFMul, nil, (*ByteCode).opI2F},
	{T_FLOAT, SLASH, T_INT}: {T_FLOAT, (*ByteCode).opFDiv, nil, (*ByteCode).opI2F},

	{T_STRING, PLUS, T_STRING}: {T_STRING, nil, nil, nil},
}

func NewTypeChecker(fileName string, program Program) *TypeChecker {
	defaultType := map[string]TypeId{
		"int":     T_INT,
		"float":   T_FLOAT,
		"string":  T_STRING,
		"boolean": T_BOOLEAN,
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
	var exists bool
	key := BinaryTypeKey{left, expr.binaryType, right}
	expr.typeInfo, exists = binaryRules[key]
	if !exists {
		return T_UNKNOWN, t.fail("Invalid Binary Expression", expr.line)
	}
	return expr.typeInfo.result, nil
}

func (t *TypeChecker) fail(message string, line uint) error {
	return error(CompileError{
		Message:  message,
		FileName: t.fileName,
		Line:     line,
		Col:      0,
	})
}
